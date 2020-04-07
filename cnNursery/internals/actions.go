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
  "github.com/diSimplex/ConTeXtNursery/logger"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "github.com/jinzhu/configor"
  "io"
  "path/filepath"
  "sync"
)

// ActionsState contains the (essentially global) state required to 
// implement the Actions RESTful interface.
//
// CONSTRAINTS: Once created, the values in this structure SHOULD only be 
// altered by structure methods.
//
type ActionsState struct {
  Mutex      sync.RWMutex
  State      control.NurseryState
  ActionsDir string
  WorkDir    string
  Actions    action.ActionList
  Ws        *webserver.WS
  Cc        *clientConnection.CC
  CNLog     *logger.LoggerType
}

// Create an ActionsState structure
//
// READS config;
// FIELD ws (ActionState);
// FIELD cc (ActionState);
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
    ActionsDir: config.Actions_Dir,
    WorkDir:    config.Work_Dir,
    Actions:    make(action.ActionList, 0),
    Ws:         ws,
    Cc:         cc,
    CNLog:      config.CNLog,
  }
}

// (re)Scan for actions in the configured Actions_Dir.
//
// Look for each *.yaml, *.toml, or *.json file, read it and store the 
// associated Action Description in the ActionList. 
//
func (aState *ActionsState) ScanForActions() {

  aState.CNLog.Logf("Looking for actions in [%s]", aState.ActionsDir)

  // walk the ActionsDirectory looking for *.config files 
  actionDescriptions, err := filepath.Glob(aState.ActionsDir+"/*.config")
  aState.CNLog.MayBeErrorf(
    err, "Could not collect *.config files from [%s]", aState.ActionsDir,
  )
  
  // for each file found we use configor to load it and capture the 
  // description storing the description in the ActionsList 
  for _, aPath := range actionDescriptions {
    aState.CNLog.Logf("Loading action description from [%s]", aPath)
    var anActionDesc action.ActionDescription
    err := configor.Load(&anActionDesc, aPath)
    if err == nil && anActionDesc.Name != "" {
      aState.Actions[anActionDesc.Name] = anActionDesc
    } else {
      aState.CNLog.MayBeErrorf(
        err, "Could not load action description from [%s]", aPath,
      )
    }
  }
}

// Returns the mapping of the currently registered actions together with a 
// brief description of each action. 
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ResponseListActionsJSON() action.ActionList {
  aState.ScanForActions()
  return aState.Actions
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ActionRunAction(string, *action.ActionConfig) string {
  return ""
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ResponseDescribeActionJSON(
  actionName string,
) action.ActionDescription {
  aState.ScanForActions()
  actionDesc := aState.Actions[actionName]
  if actionDesc.Name == "" {
    actionDesc = action.ActionDescription{
      Name: "not found",
      Desc: "could not find the action ["+actionName+"]",
    }
  }
  return actionDesc
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ResponseListActionsWithRunsJSON() map[string]string {
  return make(map[string]string, 0)
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ResponseListRunsForActionJSON(string) map[string]string {
  return make(map[string]string, 0)
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ResponseListOutputsForActionRunJSON(
  string, string,
) map[string]string {
  return make(map[string]string, 0)
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ResponseOutputFileForActionRunReader(
  string, string, string,
) (
  io.Reader, string, error,
) {
  return nil, "", nil
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ActionDeleteAll() {
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ActionDeleteRunsFor(string) {
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ActionDeleteOutputFilesFor(string, string) {
}

// TODO
//
// Part of the action.ActionImpl interface.
//
func (aState *ActionsState) ActionDeleteOutputFile(string, string, string) {
}
