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
package control

import (
  "encoding/json"
  "fmt"
  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "net/http"
  "strings"
)

//////////////////////////////////////////////////////////////////////
// Control interface types
//

// Records the current control state of a given Nursery.
//
type NurseryState struct {
  Base_Url     string
  Url_Modifier string
  State        string
  Processes    uint
}

// Records the current control state of the federation of ConTeXt
// Nurseries.
//
// The map is indexed by the Nursery Name.
//
type FederationStateMap map[string]NurseryState

const (
  StateUp     = "up"
  StatePaused = "paused"
  StateDown   = "down"
  StateKill   = "kill"
)

//////////////////////////////////////////////////////////////////////
// Control interface functions
//

// The Callbacks required to implement the Control RESTful HTTP interface
// responsible for managing the up, down, and pause state of either a given
// Nursery or the whole federation.
//
type ControlImpl interface {

  // Change the control state of this Nursery.
  //
  ActionChangeNurseryState(string)

  // Change the control state of the federation of Nurseries.
  //
  ActionChangeFederationState(string)

  // Return the control status information about the federation of ConTeXt
  // Nurseries.
  //
  // NOTE: requests to "kill" a Nursery are kept pending the completion of
  // all outstanding processes. SO in this pending state, the status
  // information SHOULD also return the number of running processes left to
  // complete.
  //
  ResponseListFederationStatusJSON() *FederationStateMap
}

// Send a control message using the client connection
//
// NOTE: We DO NOT implement the federation level (/control/all) version
//
//  interface:
//    - url: /control/up
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: brings *this* Nursery back to the "up" state
//      response: The current state of the federation
//
//    - url: /control/pause
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: brings *this* Nursery to the "pasued" state
//      response: The current state of the federation
//
//    - url: /control/kill
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: *this* cnNursery is shutdown and no longer responds
//      response: None
//
func SendNurseryControlMessage(
  baseUrl, stateChange string,
  cc *clientConnection.CC,
) *FederationStateMap {
  respBody := cc.SendJsonMessage(baseUrl, "/control/"+stateChange, http.MethodPut, []byte{})

  fmt.Printf("\ncontrol response [%s]\n\n", string(respBody))

  var fedMap FederationStateMap

  err := json.Unmarshal(respBody, &fedMap)
  if err != nil {
    cc.Log.MayBeError("Could not unmarshal respBody", err)
    fedMap = FederationStateMap{}
  }
  return &fedMap
}

// Add the Control RESTful HTTP interface to the current webserver.
//
//  interface:
//    - url: /control
//      method: GET
//      credentials: CommonName of the Client X509 certificate
//      action: None
//      response: The current state of the federation
//
//    - url: /control/up
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: brings *this* Nursery back to the "up" state
//      response: The current state of the federation
//
//    - url: /control/pause
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: brings *this* Nursery to the "pasued" state
//      response: The current state of the federation
//
//    - url: /control/kill
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: *this* cnNursery is shutdown and no longer responds
//      response: None
//
//    - url: /control/all/up
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: |
//        Walks through the federation and sends the /control/up message
//      response: The current state of the federation
//
//    - url: /control/all/pause
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: |
//        Walks through the federation and sends the /control/pause message
//      response: The current state of the federation
//
//    - url: /control/all/kill
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: |
//        Walks through the federation and sends the /control/kill message
//      response: None
//
func AddControlInterface(
  ws *webserver.WS,
  interfaceImpl ControlImpl,
) {
  ws.DescribeRoute("/control",     "???control description???", true)
  ws.DescribeRoute("/control/all", "???control/all description???", true)

//  interface:
//    - url: /control
//      method: GET
//      credentials: CommonName of the Client X509 certificate
//      action: None
//      response: The current state of the federation
//
  err := ws.AddGetHandler(
    "/control",
    func(w http.ResponseWriter, r *http.Request) {
      fedMap := interfaceImpl.ResponseListFederationStatusJSON()
      ws.ReplyInJson(w, r, fedMap)
    },
  )
  ws.Log.MayBeError("Could not add GET handler for [/control]", err)

//  interface:
//    - url: /control/up
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: brings *this* Nursery back to the "up" state
//      response: The current state of the federation
//
//    - url: /control/pause
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: brings *this* Nursery to the "pasued" state
//      response: The current state of the federation
//
//    - url: /control/kill
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: *this* cnNursery is shutdown and no longer responds
//      response: None
//
  err = ws.AddPutHandler(
    "/control",
    func(w http.ResponseWriter, r *http.Request) {

      stateChange := strings.TrimPrefix(r.URL.Path, "/control/")
      ws.Log.Logf("control stateChange: [%s]", stateChange)
      interfaceImpl.ActionChangeNurseryState(stateChange)

      fedMap := interfaceImpl.ResponseListFederationStatusJSON()
      ws.Log.Json("Control Reply: ", "fedMap", fedMap)
      ws.ReplyInJson(w, r, fedMap)
    },
  )
  ws.Log.MayBeError("Could not add PUT handler for [/control]", err)

//  interface:
//    - url: /control/all/up
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: |
//        Walks through the federation and sends the /control/up message
//      response: The current state of the federation
//
//    - url: /control/all/pause
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: |
//        Walks through the federation and sends the /control/pause message
//      response: The current state of the federation
//
//    - url: /control/all/kill
//      method: PUT
//      credentials: CommonName of the Client X509 certificate
//      action: |
//        Walks through the federation and sends the /control/kill message
//      response: None
//
  err = ws.AddPutHandler(
    "/control/all",
    func(w http.ResponseWriter, r *http.Request) {

      stateChange := strings.TrimPrefix(r.URL.Path, "/control/all/")
      ws.Log.Logf("control/all stateChange: [%s]", stateChange)
      interfaceImpl.ActionChangeFederationState(stateChange)

      fedMap := interfaceImpl.ResponseListFederationStatusJSON()
      ws.Log.Json("Control Reply: ", "fedMap", fedMap)
      ws.ReplyInJson(w, r, fedMap)
    },
  )
  ws.Log.MayBeError("Could not add PUT handler for [/control/all]", err)

}
