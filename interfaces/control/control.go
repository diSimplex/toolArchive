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
  "html/template"
  "net/http"
  "strings"
)

//////////////////////////////////////////////////////////////////////
// Control interface types
//

// Records the current control state of a given Nursery.
//
type NurseryState struct {
  State     string
  Processes uint
}

// Records the current control state of the federation of ConTeXt
// Nurseries.
//
// The map is indexed by the Nursery Name.
//
type FederationStateMap map[string]NurseryState

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

  // Return the http.Template used to format an HTML response listing the
  // control status information about the federation of ConTeXt Nurseries.
  //
  // This template expects to be bound to an FederationStateMap
  //
  ResponseListFederationStatusTemplate() *template.Template
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
  stateChange string,
  cc *clientConnection.CC,
) *FederationStateMap {
  respBody := cc.SendJsonMessage("/control/"+stateChange, http.MethodPut, []byte{})

  fmt.Printf("\ncontrol response [%s]\n\n", string(respBody))

  var fedMap FederationStateMap

  err := json.Unmarshal(respBody, &fedMap)
  if err != nil {
    cc.Log.MayBeError("Could not unmarshal respBody", err)
    fedMap = FederationStateMap{}
  }
  return &fedMap
}

// Add the Discovery RESTful HTTP interface to the current webserver.
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
  ws.DescribeRoute("/control", "???control description???")

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
      if ws.RepliedInJson(w, r, fedMap) { return }
      fedMapTemp := interfaceImpl.ResponseListFederationStatusTemplate()
      err := fedMapTemp.Execute(w, fedMap)
      ws.Log.MayBeError("Could not execute fedMapTemplate", err)
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

// TODO the action might fail since we should not honour a "kill" until
// there are not processes left... how do we deal with this case?

      fedMap := interfaceImpl.ResponseListFederationStatusJSON()
      jsonBytes, err := json.Marshal(fedMap)
      if err != nil {
        ws.Log.MayBeError("Could not json.marshal value in repliedInJson", err)
        jsonBytes = []byte{}
      }
      ws.Log.Log("Control Reply: "+string(jsonBytes))
      w.Write(jsonBytes)
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

// TODO the action might fail since we should not honour a "kill" until
// there are not processes left... how do we deal with this case?

      fedMap := interfaceImpl.ResponseListFederationStatusJSON()
      jsonBytes, err := json.Marshal(fedMap)
      if err != nil {
        ws.Log.MayBeError("Could not json.marshal value in repliedInJson", err)
        jsonBytes = []byte{}
      }
      ws.Log.Log("Control/all Reply: "+string(jsonBytes))
      w.Write(jsonBytes)
    },
  )
  ws.Log.MayBeError("Could not add PUT handler for [/control/all]", err)

}
