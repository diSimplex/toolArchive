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

// A RESTful HTTP interface responsible for communicating regular load
// average, discovery, and heartbeat messages.
//
package discovery

import (
  "encoding/json"
  //"fmt"
  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "html/template"
  "io/ioutil"
  "net/http"
)

//////////////////////////////////////////////////////////////////////
// Discovery interface types
//

// Records the amount of a particular type of memory used.
//
type MemoryTU struct {
  Total        uint64
  Used         uint64
  Percent_Used float64
}

// Records current information about a given ConTeXt Nursery.
//
type NurseryInfo struct {
  Name      string
  Port      string
  Base_Url  string
  State     string
  Processes uint
  Cores     uint
  Speed_Mhz float64
  Memory    MemoryTU
  Swap      MemoryTU
  Load      struct {
    Load1   float64
    Load5   float64
    Load15  float64
  }
}

// Records the current information about a federation of ConTeXt
// Nurseries.
//
// The map is indexed by the Nursery Name.
//
type NurseryInfoMap map[string]NurseryInfo

//////////////////////////////////////////////////////////////////////
// Discovery interface functions
//

// The Callbacks required to implement the Discovery RESTful HTTP
// interface responsible for communicating regular load average,
// discovery, and heartbeat messages.
//
type DiscoveryImpl interface {

  // Return the heartbeat status information about the federation of ConTeXt
  // Nurseries.
  //
  ResponseListNurseryInformationJSON() NurseryInfoMap

  // Return the http.Template used to format an HTML response listing
  // the heartbeat status information about the federation of ConTeXt
  // Nurseries.
  //
  // This template expects to be bound to an NurseryInfoMap 
  //
  ResponseListNurseryInformationTemplate() *template.Template

  // Update the heartbeat status information about the given Nursery in the
  // federation of Nurseries.
  //
  ActionUpdateNurseryInfo(ni NurseryInfo)
}

// Send a discovery message using the client connection
//
//  interface:
//    - url: /heartbeat
//      method: POST
//      jsonPost: NurseryInfo
//      credentials: CommonName of the Client X509 certificate
//      action: |
//        Adds or updates the NurseryInfo for the Named Nursery into the
//        Federation wide NurseryInfo map
//      response: |
//        Lists the currently known NurseryInfo of Nurseries in the Federation
//      jsonResp: NurseryInfoMap
//
func SendDiscoveryMessage(
  primaryUrl string,
  ni NurseryInfo,
  cc *clientConnection.CC,
) *NurseryInfoMap {
  jsonBytes, err := json.Marshal(ni)
  cc.Log.MayBeError("Could not marshal NurseryInfo", err)

  //fmt.Printf("\nbeat request [%s]\n\n", string(jsonBytes))

  respBody := cc.SendJsonMessage(primaryUrl, "/heartbeat", http.MethodPost, jsonBytes)

  //fmt.Printf("\nbeat response [%s]\n\n", string(respBody))

  var niMap NurseryInfoMap

  err = json.Unmarshal(respBody, &niMap)
  if err != nil {
    cc.Log.MayBeError("Could not unmarshal respBody", err)
    niMap = NurseryInfoMap{}
  }
  return &niMap
}

// Add the Discovery RESTful HTTP interface to the current webserver.
//
//  interface:
//    - url: /heartbeat
//      method: GET
//      action: None
//      credentials: CommonName of the Client X509 certificate
//      response: |
//        Lists the currently known NurseryInfo of Nurseries in the Federation
//      jsonResp: NurseryInfoMap
//
//    - url: /heartbeat
//      method: POST
//      jsonPost: NurseryInfo
//      credentials: CommonName of the Client X509 certificate
//      action: |
//        Adds or updates the NurseryInfo for the Named Nursery into the
//        Federation wide NurseryInfo map
//      response: |
//        Lists the currently known NurseryInfo of Nurseries in the Federation
//      jsonResp: NurseryInfoMap
//
func AddDiscoveryInterface(
  ws *webserver.WS,
  interfaceImpl DiscoveryImpl,
) {
  ws.DescribeRoute("/heartbeat", "???heartbeat description???", true)

//  interface:
//    - url: /heartbeat
//      method: GET
//      action: None
//      credentials: CommonName of the Client X509 certificate
//      response: |
//        Lists the currently known NurseryInfo of Nurseries in the Federation
//      jsonResp: NurseryInfoMap
//
  err := ws.AddGetHandler(
    "/heartbeat",
    func(w http.ResponseWriter, r *http.Request) {
      niMap := interfaceImpl.ResponseListNurseryInformationJSON()
      if ws.RepliedInJson(w, r, niMap) { return }
      niMapTemp := interfaceImpl.ResponseListNurseryInformationTemplate()
      err := niMapTemp.Execute(w, niMap)
      ws.Log.MayBeError("Could not execute niMapTemplate", err)
    },
  )
  ws.Log.MayBeError("Could not add GET handler for [/heartbeat]", err)

//  interface:
//    - url: /heartbeat
//      method: POST
//      jsonPost: NurseryInfo
//      credentials: CommonName of the Client X509 certificate
//      action: |
//        Adds or updates the NurseryInfo for the Named Nursery into the
//        Federation wide NurseryInfo map
//      response: |
//        Lists the currently known NurseryInfo of Nurseries in the Federation
//      jsonResp: NurseryInfoMap
//
  err = ws.AddPostHandler(
    "/heartbeat",
    func(w http.ResponseWriter, r *http.Request) {
      body, err := ioutil.ReadAll(r.Body)
      if err != nil {
        ws.Log.MayBeError("Could not read body of /heartbeat post request", err)
        http.Error(w, "can't read body", http.StatusBadRequest)
        return
      }
//      ws.Log.Log("heartBeat body: "+string(body))
      var ni NurseryInfo
      err = json.Unmarshal(body, &ni)
      if err != nil {
        ws.Log.MayBeError("Could not unmarshal heartbeat body", err)
        ni = NurseryInfo{}
      }
      interfaceImpl.ActionUpdateNurseryInfo(ni)

      niMap := interfaceImpl.ResponseListNurseryInformationJSON()
      if ws.RepliedInJson(w, r, niMap) { return }
      http.Redirect(w,r, "/hearBeat", http.StatusSeeOther)
    },
  )
  ws.Log.MayBeError("Could not add POST handler for [/heartbeat]", err)

}
