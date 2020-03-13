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
//  "bytes"
//  "encoding/json"
  "fmt"
  "html/template"
//  "io/ioutil"
//  "math/rand"
  "net/http"
//  "time"
)

////////////////////
// Control interface

type CNState struct {
  Mutex sync.RWMutex
  State control.NurseryState
}

func CreateCNState() *CNState {
  return &CNState{}
}

func (cnState *CNState) ActionChangeNurseryState(stateChange string) {

}

func (cnState *CNState) ActionChangeFederationState(stateChange string) {

}

func handleControl() {

  addBasePagePath("control", "control description")

  http.HandleFunc("/control", func(w http.ResponseWriter, r *http.Request) {
    cnNurseryLogf("url: [%s] method: [%s]", r.URL.Path, r.Method)

    if r.Method == http.MethodGet {
      ctlPath := r.URL.Path
      fmt.Printf("control path: [%s]\n", ctlPath)
//      hbInfoMap := getHeartBeatInfoMap()
//      cnNurseryJson("hbInfoMap: ", "NurseryInfoMap", hbInfoMap)

//      if repliedInJson(w, r, hbInfoMap) { return }

      // we are replying to a (human) browser

//      hbTemplate := heartBeatTemplate()
//      err := hbTemplate.Execute(w, hbInfoMap)
//      if err != nil {
//        cnNurseryMayBeError("Could not execute heart beat template", err)
//        w.Write([]byte("Could not provide any federation information\nPlease try a$
//      }
      w.Write([]byte("Hello from control\n"))
      return
    }
    http.Error(w, "Incorrect Request Method for /control", http.StatusMethodNotAllowed)
  })

}

func controlTemplate() *template.Template {
  controlTemplateStr := `
  <body>
    <h1>Federation Control Information</h1>
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

  theTemplate, err := theTemplate.Parse(controlTemplateStr)
  cnNurseryMayBeFatal("Could not parse the internal control template", err)

  return theTemplate
}

