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

// This file collects all of the globals required for the cnNursery process.
//
// Since cnNursery makes essential use of multi-threading, we need to 
// ensure all globals are thread safe. To do this we make liberal use
// of the sync.RWMutexes, one for each global.
//

package main

import (
  "github.com/diSimplex/ConTeXtNursery/interfaces/discovery"
  "html/template"
  "strings"
  "sync"
)

type CNInfoMap struct {
  IsPrimary bool
  Mutex     sync.RWMutex
  NI        discovery.NurseryInfoMap
}

func CreateCNInfoMap() *CNInfoMap {
  lConfig := getConfig()

  infoMap := CNInfoMap{}
  infoMap.IsPrimary = strings.Contains(lConfig.Primary_Url, lConfig.Name)
  infoMap.NI        = make(discovery.NurseryInfoMap)
  return &infoMap
}

func (cniMap *CNInfoMap) ActionUpdateNurseryInfo(ni discovery.NurseryInfo) {
  cniMap.Mutex.Lock()
  defer cniMap.Mutex.Unlock()

  if ni.Name != "" {
    cniMap.NI[ni.Name] = ni
  }
}

func (cniMap *CNInfoMap) ActionUpdateNurseryInfoMap(niMap *discovery.NurseryInfoMap) {
  cniMap.Mutex.Lock()
  defer cniMap.Mutex.Unlock()

  if !cniMap.IsPrimary && (0 < len(*niMap)) {
    cniMap.NI = *niMap
  }
}

func (cniMap *CNInfoMap) ResponseListNurseryInformationJSON() discovery.NurseryInfoMap {
  cniMap.Mutex.RLock()
  defer cniMap.Mutex.RUnlock()

  return cniMap.NI
}

func (cniMap *CNInfoMap) ResponseListNurseryInformationTemplate() *template.Template {
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
        <td><a href="https://{{$value.Name}}:{{$value.Port}}">{{$value.Name}}</a><$
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
  cnLog.MayBeFatal("Could not parse the internal heartBeat template", err)

  return theTemplate
}

