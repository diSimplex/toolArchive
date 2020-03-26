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

package CNNurseries

import (
  "github.com/diSimplex/ConTeXtNursery/interfaces/discovery"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "html/template"
  "strings"
  "sync"
)

// The CNInfoMap collects a cnNursery's map of all other cnNurseries in 
// the federation. 
//
// CONSTRAINTS: Once created, the values in this structure SHOULD only be 
// altered by structure methods.
//
type CNInfoMap struct {
  IsPrimary bool
  Name      string
  Mutex     sync.RWMutex
  NI        discovery.NurseryInfoMap
  CNLog     *logger.LoggerType
}

func CreateCNInfoMap(
  config *ConfigType,
) *CNInfoMap {
  infoMap := CNInfoMap{}
  infoMap.IsPrimary = strings.Contains(config.Primary_Url, config.Name)
  infoMap.Name      = config.Name
  infoMap.NI        = make(discovery.NurseryInfoMap)
  infoMap.CNLog     = config.CNLog
  return &infoMap
}

type ANurseryAction func(string, discovery.NurseryInfo)

func (cniMap *CNInfoMap) DoToAllOthers(anAction ANurseryAction) {
  cniMap.Mutex.Lock()
  defer cniMap.Mutex.Unlock()

  for aKey, aValue := range cniMap.NI {
    if aKey == cniMap.Name { continue } // do not do this to myself!
    anAction(aKey, aValue)
  }
}

func (cniMap *CNInfoMap) DoToAll(anAction ANurseryAction) {
  cniMap.Mutex.Lock()
  defer cniMap.Mutex.Unlock()

  for aKey, aValue := range cniMap.NI {
    anAction(aKey, aValue)
  }
}

func (cniMap *CNInfoMap) DeleteAll(deadNurseries []string) {
  cniMap.Mutex.Lock()
  defer cniMap.Mutex.Unlock()

  for _, aNursery := range deadNurseries {
    delete(cniMap.NI, aNursery)
  }
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
  <head>
    <title>Federation Heart Beat Information</title>
    <meta http-equiv="refresh" content="5" />
  </head>
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
        <td><a href="{{$value.Base_Url}}">{{$value.Name}}</a></td>
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
  cniMap.CNLog.MayBeFatal("Could not parse the internal heartBeat template", err)

  return theTemplate
}

