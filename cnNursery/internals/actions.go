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
  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/interfaces/action"
  "github.com/diSimplex/ConTeXtNursery/interfaces/control"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "html/template"
  "io"
  "sync"
)

// ActionsState contains the (essentially global) state required to 
// implement the Actions RESTful interface.
//
// CONSTRAINTS: Once created, the values in this structure SHOULD only be 
// altered by structure methods.
//
type ActionsState struct {
  Mutex sync.RWMutex
  State control.NurseryState
  Ws    *webserver.WS
  Cc    *clientConnection.CC
}

// Create an ActionsState structure
//
// READS config;
// READS ws;
// READS cc;
//
func CreateActionsState(
  config *ConfigType,
  ws     *webserver.WS,
  cc     *clientConnection.CC,
) *ActionsState {
  return &ActionsState{
    State: control.NurseryState{
      Base_Url:     config.Base_Url,
      Url_Modifier: "",
      State:        "up",
      Processes:    0,
    },
    Ws: ws,
    Cc: cc,
  }
}

// (re)Scan for actions in the configured Actions_Dir.
//
// Look for each *.yaml, *.toml, or *.json file, read it and store the 
// associated Action Description in the ActionList. 
//
func (aState *ActionsState) ScanForActions() {

}

func (aState *ActionsState) ResponseListActionsJSON() action.ActionList {
  return nil
}

func (aState *ActionsState) ResponseListActionsTemplate() *template.Template {
  return nil
}

func (aState *ActionsState) ActionRunAction(string, *action.ActionConfig) string {
  return ""
}

func (aState *ActionsState) ResponseDescribeActionJSON() *action.ActionConfig {
  return nil
}

func (aState *ActionsState) ResponseDescribeActionTemplate() *template.Template {
  return nil
}

func (aState *ActionsState) ResponseListActionsWithRunsJSON() map[string]string {
  return make(map[string]string, 0)
}

func (aState *ActionsState) ResponseListActionsWithRunsTemplate() *template.Template {
  return nil
}

func (aState *ActionsState) ResponseListRunsForActionJSON(string) map[string]string {
  return make(map[string]string, 0)
}

func (aState *ActionsState) ResponseRunsForActionTemplate() *template.Template {
  return nil
}

func (aState *ActionsState) ResponseListOutputsForActionRunJSON(
  string, string,
) map[string]string {
  return make(map[string]string, 0)
}

func (aState *ActionsState) ResponseOutputsForActionRunTemplate() *template.Template {
  return nil
}

func (aState *ActionsState) ResponseOutputFileForActionRunReader(
  string, string, string,
) (
  io.Reader, string, error,
) {
  return nil, "", nil
}

func (aState *ActionsState) ResponseOutputFileTemplate() *template.Template {
  return nil
}

func (aState *ActionsState) ActionDeleteAll() {
}

func (aState *ActionsState) ActionDeleteRunsFor(string) {
}

func (aState *ActionsState) ActionDeleteOutputFilesFor(string, string) {
}

func (aState *ActionsState) ActionDeleteOutputFile(string, string, string) {
}
