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
  "encoding/json"
  "fmt"
  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "io"
  "io/ioutil"
  "net/http"
)

//////////////////////////////////////////////////////////////////////
// Action interface types
//

type ArgValue struct {
  Key   string
  Value string
}

type Arguments []ArgValue

type EnvValue struct {
  Key   string
  Value string
}

type EnvVars   []EnvValue

type ActionConfig struct {
  Args Arguments
  Envs EnvVars
}

type ArgumentDesc struct {
  Key  string
  Desc string
}

type ArgumentDescs []ArgumentDesc

type EnvironmentDesc struct {
  Key  string
  Desc string
}

type EnvironmentDescs []EnvironmentDesc

type ActionDescription struct {
  Name string
  Desc string
  Args ArgumentDescs
  Envs EnvironmentDescs
}

// A map of currently registered actions together with a brief description 
// of each action.
//
type ActionList map[string]ActionDescription

//////////////////////////////////////////////////////////////////////
// Action interface functions
//

type ActionImpl interface {

  // Returns the mapping of the currently registered actions
  // together with a brief description of each action.
  //
  ResponseListActionsJSON() ActionList

  ActionRunAction(string, *ActionConfig) string

  ResponseDescribeActionJSON(string) ActionDescription

  ResponseListActionsWithRunsJSON() map[string]string

  ResponseListRunsForActionJSON(string) map[string]string

  ResponseListOutputsForActionRunJSON(string, string) map[string]string

  ResponseOutputFileForActionRunReader(
    string, string, string,
  ) (
    io.Reader, string, error,
  )

  ActionDeleteAll()

  ActionDeleteRunsFor(string)

  ActionDeleteOutputFilesFor(string, string)

  ActionDeleteOutputFile(string, string, string)
}



// Send an action request using the client connection
//
// interface:
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
func SendActionRequestToNursery(
  baseUrl      string,
  action       string,
  actionConfig *ActionConfig,
  cc           *clientConnection.CC,
) {
  jsonBytes, err := json.Marshal(actionConfig)
  cc.Log.MayBeError("Could not marshal action configuration", err)

  fmt.Printf("\naction request [%s]\n\n", string(jsonBytes))

  respBody := cc.SendJsonMessage(
    baseUrl,
    "/action/"+action,
    http.MethodPost,
    jsonBytes,
  )

  fmt.Printf("\naction response [%s]\n\n", string(respBody))

  // TODO
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
//   - url: /action/output
//     method: GET
//     action: None
//     response: List of actions which have runs associated with them
//     jsonResp: []string
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
//   - url: /action/output/<anAction>
//     method: DELETE
//     action: |
//       Deletes all of the associated <aRun>s of the <anAction>
//     response: List (remaining) runs associated with this action
//     jsonResp: []string
//
//   - url: /action/output/<anAction>/<aRun>
//     method: DELETE
//     action: |
//       Clears the associated <aRun> of the <anAction>
//     response: List (remaining) runs associated with this action
//     jsonResp: []string
//
//   - url: /action/output/<anAction>/<aRun>/<outputFile>
//     method: DELETE
//     action: |
//       Deletes the <outputFile> associated with <aRun> of the <anAction>.
//     response: List (remaining) output files associated with this action
//     jsonResp: []string
//
func AddActionInterface(
  ws *webserver.WS,
  interfaceImpl ActionImpl,
) {
  ws.DescribeRoute("/action", "???action description???", true)
  ws.DescribeRoute("/action/output", "???action/output description???", true)

  // interface:
  //   - url: /action
  //     method: GET
  //     credentials: CommonName of the Client X509 certificate
  //     action: None
  //     response: The list of currently registered actions
  //     jsonResp: map[string]string
  //
  //   - url: /action/<anAction>
  //     method: GET
  //     action: None
  //     response: List the available action arguments and environment variables.
  //     jsonResp: ActionConfig
  //
  err := ws.AddGetHandler(
    "/action",
    func(w http.ResponseWriter, r *http.Request) {
      pathParts := ws.GetPathParts(r.URL.Path)
      if len(pathParts) < 2 {
        //
        // List currently registered actions
        //
        actions := interfaceImpl.ResponseListActionsJSON()
        ws.ReplyInJson(w, r, actions)
      } else {
        //
        // Describe anAction
        //
        actionDesc := interfaceImpl.ResponseDescribeActionJSON(pathParts[1])
        ws.ReplyInJson(w, r, actionDesc)
      }
    },
  )
  ws.Log.MayBeError("Could not add GET handler for [/action]", err)

  // interface:
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
  err = ws.AddPostHandler(
    "/action",
    func(w http.ResponseWriter, r *http.Request) {
      pathParts := ws.GetPathParts(r.URL.Path)
      if len(pathParts) < 2 {
        ws.Log.MayBeError("No action specified in /action post request", err)
        http.Error(w, "No action specified", http.StatusBadRequest)
        return
      }
      //
      // Run <anAction>
      //
      body, err := ioutil.ReadAll(r.Body)
      theAction := pathParts[1]
      if err != nil {
        ws.Log.MayBeError("Could not read body of /action post request", err)
        http.Error(w, "Could not read body", http.StatusBadRequest)
        return
      }
      ws.Log.Logf("[%s] action body: %s", theAction, string(body))
      var ac ActionConfig
      err = json.Unmarshal(body, &ac)
      if err != nil {
        ws.Log.MayBeError("Could not unmarshal action configuration body", err)
        http.Error(
          w,
          "Could not unmarshal action configuration",
          http.StatusBadRequest,
        )
        return
      }
      theRunId := interfaceImpl.ActionRunAction(theAction, &ac)

      http.Redirect(
        w, r,
        "/action/output/"+theAction+"/"+theRunId,
        http.StatusSeeOther,
      )
    },
  )
  ws.Log.MayBeError("Could not add POST handler for [/action]", err)

  // interface:
  //   - url: /action/output
  //     method: GET
  //     action: None
  //     response: List of actions which have runs associated with them
  //     jsonResp: []string
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
  err = ws.AddGetHandler(
    "/action/output",
    func(w http.ResponseWriter, r *http.Request) {
      pathParts  := ws.GetPathParts(r.URL.Path)
      theAction  := ""
      theRun     := ""
      outputFile := ""
      if 1 < len(pathParts) { theAction  = pathParts[1] }
      if 2 < len(pathParts) { theRun     = pathParts[2] }
      if 3 < len(pathParts) { outputFile = pathParts[3] }

      if len(pathParts) < 2 {
        //
        // List actions with assocaited runs
        //
        actions := interfaceImpl.ResponseListActionsWithRunsJSON()
        ws.ReplyInJson(w, r, actions)
      } else if len(pathParts) < 3 {
        //
        // List runs associated with <anAction>
        //
        runs := interfaceImpl.ResponseListRunsForActionJSON(theAction)
        ws.ReplyInJson(w, r, runs)
      } else if len(pathParts) < 4 {
        //
        // List output files associated with <aRun> of <anAction>
        //
        outputFiles :=
          interfaceImpl.ResponseListOutputsForActionRunJSON(theAction, theRun)
        ws.ReplyInJson(w, r, outputFiles)
      } else {
        //
        // Browse <outputFile> assocaited with <aRun> of <anAction>
        //
        ofReader, mimeType, err := 
          interfaceImpl.ResponseOutputFileForActionRunReader(
            theAction,
            theRun,
            outputFile,
          )
        ws.Log.MayBeError("Could not get output file for action", err)
        ws.ReplyAsRawFile(w, r, ofReader, mimeType)
      }
    },
  )
  ws.Log.MayBeError("Could not add GET handler for [/action]", err)

  // interface:
  //   - url: /action/output
  //     method: DELETE
  //     action: |
  //       Deletes all runs associated with any action
  //     response: Redirects to (GET) /action
  //
  //   - url: /action/output/<anAction>
  //     method: DELETE
  //     action: |
  //       Deletes all of the associated <aRun>s of the <anAction>
  //     response: Redirects to (GET) /action/output
  //
  //   - url: /action/output/<anAction>/<aRun>
  //     method: DELETE
  //     action: |
  //       Clears the associated <aRun> of the <anAction>
  //     response: Redirects to (GET) /action/output/<anAction>
  //
  //   - url: /action/output/<anAction>/<aRun>/<outputFile>
  //     method: DELETE
  //     action: |
  //       Deletes the <outputFile> associated with <aRun> of the <anAction>.
  //     response: Redirects to (GET) /action/output/<anAction/<aRun>
  //
  err = ws.AddDeleteHandler(
    "/action/output",
    func(w http.ResponseWriter, r *http.Request) {
      pathParts  := ws.GetPathParts(r.URL.Path)
      theAction  := ""
      theRun     := ""
      outputFile := ""
      if 1 < len(pathParts) { theAction  = pathParts[1] }
      if 2 < len(pathParts) { theRun     = pathParts[2] }
      if 3 < len(pathParts) { outputFile = pathParts[3] }

      if len(pathParts) < 2 {
        //
        // Delete all runs associated with any action
        //
        interfaceImpl.ActionDeleteAll()
        http.Redirect(w, r, "/action", http.StatusSeeOther)
      } else if len(pathParts) < 3 {
        //
        // Delete all runs associated with <anAction>
        //
        interfaceImpl.ActionDeleteRunsFor(theAction)
        http.Redirect(w, r, "/action/output", http.StatusSeeOther)
      } else if len(pathParts) < 4 {
        //
        // Delete all output file associated with <aRun> of <anAction>
        //
        interfaceImpl.ActionDeleteOutputFilesFor(theAction, theRun)
        http.Redirect(w, r, "/action/output/"+theAction, http.StatusSeeOther)
      } else {
        //
        // Delete <outputFile> assocaited with <aRun> of <anAction>
        //
        interfaceImpl.ActionDeleteOutputFile(theAction, theRun, outputFile)
        http.Redirect(
          w, r,
          "/action/output/"+theAction+"/"+theRun,
          http.StatusSeeOther,
        )
      }
    },
  )
  ws.Log.MayBeError("Could not add GET handler for [/action]", err)

}

