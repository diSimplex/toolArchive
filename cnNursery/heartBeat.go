// Copyright 2020 PerceptiSys Ltd, (Stephen Gaito)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
  "bytes"
  "crypto/tls"
  "encoding/json"
  "fmt"
//  "github.com/bvinc/go-sqlite-lite/sqlite3"
//  "github.com/cornelk/hashmap"
  "github.com/shirou/gopsutil/cpu"
  "github.com/shirou/gopsutil/load"
  "github.com/shirou/gopsutil/mem"
  "html/template"
  "io/ioutil"
  "math/rand"
  "net/http"
//  "strings"
//  "sync"
  "time"
)

func sendPeriodicHeartBeats() {
  lConfig := getConfig()

  // Setup HTTPS client
  tlsConfig := &tls.Config{
    ClientAuth:     tls.RequireAndVerifyClientCert,
    Certificates: []tls.Certificate{serverCert},
    RootCAs:        caCertPool,
    ClientCAs:      caCertPool,
  }

  transport := &http.Transport{
    TLSClientConfig:    tlsConfig,
    ForceAttemptHTTP2:  true,
    MaxIdleConns:       10,
    IdleConnTimeout:    30 * time.Second,
    DisableCompression: true,
  }

  client := &http.Client{
    Transport: transport,
  }

  for {
    time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
    ni := NurseryInfo{
      Name: lConfig.Name,
      Port: lConfig.Port,
      State: "up",
      Processes: 1,
    }

    loads, err := load.Avg()
    if err != nil {
      cnNurseryMayBeError("Could not read the load average", err)
      loads = &load.AvgStat{ Load1: 1.0, Load5: 1.0, Load15: 1.0, }
    }
    ni.Load.Load1  = loads.Load1
    ni.Load.Load5  = loads.Load5
    ni.Load.Load15 = loads.Load15

    cpuInfo, err := cpu.Info()
    if err != nil {
      cnNurseryMayBeError("Could not read the cpu information", err)
      cpuInfo = []cpu.InfoStat{ cpu.InfoStat{ Cores: 1, Mhz: 1000 } }
    }
    ni.Cores     = uint(len(cpuInfo))
    ni.Speed_Mhz = cpuInfo[0].Mhz

    virtMem, err := mem.VirtualMemory()
    if err != nil {
      cnNurseryMayBeError("Could not read the virtual memory information", err)
      virtMem = &mem.VirtualMemoryStat{ Total: 1000, Used: 1000 }
    }
    ni.Memory.Total = virtMem.Total
    ni.Memory.Used  = virtMem.Used

    swapMem, err := mem.SwapMemory()
    if err != nil {
      cnNurseryMayBeError("Could not read the swap memory information", err)
      swapMem = &mem.SwapMemoryStat{ Total: 1000, Used: 1000 }
    }
    ni.Swap.Total = swapMem.Total
    ni.Swap.Used  = swapMem.Used

//    jsonBytes, err := json.MarshalIndent(ni, "", "  ")
    jsonBytes, err := json.Marshal(ni)
    fmt.Printf("\nbeat request [%s]\n", string(jsonBytes))

    hbReq, err := http.NewRequest(http.MethodPost, 
      lConfig.Primary_Url + "/heartbeat",
      bytes.NewReader(jsonBytes),
    )
    if err != nil {
      cnNurseryMayBeError("Could not create heart beat request", err)
      continue
    }

    resp, err := client.Do(hbReq)
    if err != nil {
      cnNurseryMayBeError("Could not send heart beat request to the primary Nursery", err)
      continue
    }
    defer resp.Body.Close()

    respBody, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    if err != nil {
      cnNurseryMayBeFatal("Could not read the body of the heart beat response", err)
      continue
    }

    fmt.Printf("beat response [%s]\n\n", string(respBody))

    jsonUnmarshalFederationInfo(respBody)
  }

}

func handleHeartBeats() {

  initHeartBeatInfo()

  http.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
    cnNurseryLogf("url: [%s] method: [%s]", r.URL.Path, r.Method)

    if r.Method == http.MethodPost {

      body, err := ioutil.ReadAll(r.Body)
      if err != nil {
        cnNurseryMayBeError("Could not ready body of post request", err)
        http.Error(w, "can't read body", http.StatusBadRequest)
        return
      }
      cnNurseryLog("heartBeat body: "+string(body))

      err = setHeartBeatInfoNurseryFromJsonBytes(body)
      cnNurseryMayBeError("Could not set the heart beat info from the POST body", err)

      jsonBytes, err := jsonMarshalHeartBeatInfo()
      cnNurseryMayBeError("Could not marshal heart beat info", err)
      cnNurseryLog("heartBeat fi: "+string(jsonBytes))
      w.Write(jsonBytes)
      return
    }

    if r.Method == http.MethodGet {
      hbInfoMap := getHeartBeatInfoMap()
      cnNurseryJson("hbInfoMap: ", "NurseryInfoMap", hbInfoMap)

      if repliedInJson(w, r, hbInfoMap) { return }

      // we are replying to a (human) browser

      hbTemplate := heartBeatTemplate()
      err := hbTemplate.Execute(w, hbInfoMap)
      if err != nil {
        cnNurseryMayBeError("Could not execute heart beat template", err)
        w.Write([]byte("Could not provide any federation information\nPlease try again!"))
      }
      return
    }
    http.Error(w, "Incorrect Request Method", http.StatusMethodNotAllowed)
  })

}

func heartBeatTemplate() *template.Template {
  heartBeatTemplateStr := `
  <body>
    <h1>Federation Heart Beat Information</h1>
    <table>
      <tr>
        <th>Name</th>
        <th>Port</th>
        <th>State</th>
        <th>Processes</th>
        <th>Cores</th>
        <th>Speed Mhz</th>
        <th>Mem Total</th>
        <th>Mem Used</th>
        <th>Swap Total</th>
        <th>Swap Used</th>
        <th>Load 1 min</th>
        <th>Load 5 min</th>
        <th>Load 15 min</th>
      </tr>
{{ range $key, $value := . }}
      <tr>
        <td>{{$value.Name}}</td>
        <td>{{$value.Port}}</td>
        <td>{{$value.State}}</td>
        <td>{{$value.Processes}}</td>
        <td>{{$value.Cores}}</td>
        <td>{{$value.Speed_Mhz}}</td>
        <td>{{$value.Memory.Total}}</td>
        <td>{{$value.Memory.Used}}</td>
        <td>{{$value.Swap.Total}}</td>
        <td>{{$value.Swap.Used}}</td>
        <td>{{$value.Load.Load1}}</td>
        <td>{{$value.Load.Load5}}</td>
        <td>{{$value.Load.Load15}}</td>
      </tr>
{{ end }}
    </table>
  </body>
`
  theTemplate := template.New("body")

  theTemplate, err := theTemplate.Parse(heartBeatTemplateStr)
  cnNurseryMayBeFatal("Could not parse the internal heartBeat template", err)

  return theTemplate
}
