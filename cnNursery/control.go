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
  "github.com/diSimplex/ConTeXtNursery/interfaces/discovery"
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
  lConfig := getConfig()
  return &CNState{
    State: control.NurseryState{
      Base_Url:     lConfig.Base_Url,
      Url_Modifier: "",
      State:        "up",
      Processes:    0,
    },
    Ws: ws,
    Cc: cc,
  }
}

func (cnState *CNState) SetState(newState string) {
  cnState.Mutex.Lock()
  defer cnState.Mutex.Unlock()

  cnState.State.State = newState // this is too permissive! but works for now.
}

func (cnState *CNState) GetState() string {
  cnState.Mutex.RLock()
  defer cnState.Mutex.RUnlock()

  return cnState.State.State
}

func (cnState *CNState) ActionChangeNurseryState(stateChange string) {
  switch stateChange {
    case control.StateUp     : cnState.SetState(stateChange)
    case control.StatePaused : cnState.SetState(stateChange)
    case control.StateDown   : cnState.SetState(stateChange)
    case control.StateKill   : cnState.SetState(stateChange)
      cnState.Ws.Server.Shutdown(context.Background())
    default                  :
      cnLog.Logf("Ignoring incorrect state change: [%s]", stateChange)
  }
}

func (cnState *CNState) ActionChangeFederationState(stateChange string) {
  cnInfoMap.DoToAllOthers(func (name string, ni discovery.NurseryInfo) {
    control.SendNurseryControlMessage(ni.Base_Url, stateChange, cnState.Cc)
  })
  lConfig := getConfig()
  control.SendNurseryControlMessage(lConfig.Base_Url, stateChange, cnState.Cc)
}

func (cnState *CNState) ResponseListFederationStatusJSON() *control.FederationStateMap {
  lConfig         := getConfig()
  fedStateMap     := control.FederationStateMap{}
  fedNumProcesses := uint(0)
  cnInfoMap.DoToAll(func(name string, ni discovery.NurseryInfo) {
    ns := control.NurseryState{
      Base_Url:     ni.Base_Url,
      Url_Modifier: "",
      State:        ni.State,
      Processes:    ni.Processes,
    }
    fedNumProcesses = fedNumProcesses + ni.Processes
    fedStateMap[name] = ns
  })
  fedStateMap["Federation"] = control.NurseryState{
    Base_Url:     lConfig.Primary_Url,
    Url_Modifier: "/all",
    State:        control.StateUp,
    Processes:    fedNumProcesses,
  }
  return &fedStateMap
}

func (cnState *CNState) ResponseListFederationStatusTemplate() *template.Template {
  controlTemplateStr := `
  <body>
    <h1>Federation Control Information</h1>
    <table>
      <tr>
        <th colspan=2></th>
        <th colspan=4>State</th>
      </tr>
      <tr>
        <th colspan=2></th>
        <th colspan=4><hr style="margin-top:0em;margin-bottom:0em;" /></th>
      </tr>
      <tr>
        <th>Name</th>
        <th>Processes</th>
        <th>Current</th>
        <th>Up</th>
        <th>Pause</th>
        <th>Kill</th>
      </tr>
{{ range $key, $value := . }}
      <tr>
        <td><a href="{{$value.Base_Url}}">{{$key}}</a></td>
        <td>{{$value.Processes}}</td>
        <td>{{$value.State}}</td>
        <td>
          <form method="post"
           action="{{$value.Base_Url}}/control{{$value.Url_Modifier}}/up?method=put">
            <input type="submit" value="Up" />
          </form>
        </td>
        <td>
          <form method="post"
           action="{{$value.Base_Url}}/control{{$value.Url_Modifier}}/paused?method=put">
            <input type="submit" value="Pause" />
          </form>
        </td>
        <td>
          <form method="post"
           action="{{$value.Base_Url}}/control{{$value.Url_Modifier}}/kill?method=put">
            <input type="submit" value="Kill" />
          </form>
        </td>
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
