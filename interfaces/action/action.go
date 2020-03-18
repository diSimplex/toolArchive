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

// A RESTful HTTP interface responsible for managing the up, down, and
// pause state of either a given Nursery or the whole federation.
//
package action

import (
//  "encoding/json"
//  "fmt"
//  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "html/template"
  "net/http"
//  "strings"
)

//////////////////////////////////////////////////////////////////////
// Action interface types
//

type Arguments []string

type EnvValue struct {
  key   string
  value string
}

type EnvVars   []EnvValue

type ActionConfig struct {
  Args Arguments
  Envs EnvVars
}

type ArgDesc struct {
  key  string
  desc string
}

type ArgumentDescs []string

type EnvDesc struct {
  key  string
  desc string
}

type EnvDescs []EnvDesc

type ActionDescription struct {
  Arg ArgumentDescs
  Env EnvDescs
}

//////////////////////////////////////////////////////////////////////
// Action interface functions
//

type ActionImpl interface {

  ResponseListActionsJSON() *ActionDescription

  ResponseListActionsTemplate() *template.Template

}

// Add the Action RESTful HTTP interface to the current webserver.
//
// interface:
//   - url: /action
//     method: GET
//     credentials: CommonName of the Client X509 certificate
//     action: None
//     response: The list of currently registered actions
//     jsonResp: []string
//
//   - url: /action/<anAction>
//     method: GET
//     action: None
//     response: List the available action arguments and environment variables.
//     jsonResp: ActionConfig
//
//   - url: /action/<anAction>
//     method: POST
//     jsonPost: ActionConfig
//     credentials: CommonName of the Client X509 certificate
//     action: Runs the <anAction>
//     response: |
//       Redirect to output file browser which longPolls the log file produced
//       by this action. (Note we could use mithril.js in an AJAX "pull" model
//       to ensure the user does not see the whole page refresh).
//
//   - url: /action/output/<anAction>
//     method: GET
//     action: None
//     response: List of available runs associated with this action
//     jsonResp: []string
//
//   - url: /action/output/<anAction>/<aRun>
//     method: GET
//     action: None
//     response: |
//       List the output files associated with <aRun> of the <anAction>.
//     jsonResp: []string
//
//   - url: /action/output/<anAction>/<aRun>/<outputFile>
//     method: GET
//     action: None
//     response: |
//       Browse the <outputFile> associated with <aRun> of the <anAction>.
//
//   - url: /action/clear/<anAction>/<aRun>
//     method: DELETE
//     action: |
//       Clears the associated <aRun> of the <anAction> (or all runs if no
//       <aRun> is provided)
//     response: List (remaining) runs associated with this action
//     jsonResp: []string
//
func AddActionInterface(
  ws *webserver.WS,
  interfaceImpl ActionImpl,
) {
  ws.DescribeRoute("/action", "???action description???")
  ws.DescribeRoute("/action/output", "???action/output description???")

// interface:
//   - url: /action
//     method: GET
//     credentials: CommonName of the Client X509 certificate
//     action: None
//     response: The list of currently registered actions
//     jsonResp: map[string]string
//
  err := ws.AddGetHandler(
    "/action",
    func(w http.ResponseWriter, r *http.Request) {
      actions := interfaceImpl.ResponseListActionsJSON()
      if ws.RepliedInJson(w, r, actions) { return }
      actionsTemp := interfaceImpl.ResponseListActionsTemplate()
      err := actionsTemp.Execute(w, actions)
      ws.Log.MayBeError("Could not execute actionsTemplate", err)
    },
  )
  ws.Log.MayBeError("Could not add GET handler for [/action]", err)



}

