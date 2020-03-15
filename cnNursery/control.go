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
  "context"
//  "encoding/json"
//  "fmt"
  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/interfaces/control"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "html/template"
//  "io/ioutil"
//  "math/rand"
//  "net/http"
  "sync"
//  "time"
)

////////////////////
// Control interface

type CNState struct {
  Mutex sync.RWMutex
  State control.NurseryState
  Ws    *webserver.WS
  Cc    *clientConnection.CC
}

func CreateCNState(ws *webserver.WS, cc *clientConnection.CC) *CNState {
  return &CNState{Ws: ws, Cc: cc}
}

func (cnState *CNState) ActionChangeNurseryState(stateChange string) {
  switch stateChange {
    case control.StateUp     : cnState.State.State = stateChange
    case control.StatePaused : cnState.State.State = stateChange
    case control.StateDown   : cnState.State.State = stateChange
    case control.StateKill   : cnState.State.State = stateChange
      cnState.Ws.Server.Shutdown(context.Background())
  }
}

func (cnState *CNState) ActionChangeFederationState(stateChange string) {
  cnInfoMap.DoToAllOthers(func (baseUrl string) {
    control.SendNurseryControlMessage(baseUrl, stateChange, cnState.Cc)
  })
  lConfig := getConfig()
  control.SendNurseryControlMessage(lConfig.Base_Url, stateChange, cnState.Cc)
}

func (cnState *CNState) ResponseListFederationStatusJSON() *control.FederationStateMap {

  return nil
}

func (cnState *CNState) ResponseListFederationStatusTemplate() *template.Template {
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
  cnLog.MayBeFatal("Could not parse the internal control template", err)

  return theTemplate
}
