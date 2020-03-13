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

// This file collects all of the globals required for the cnNursery
//  process.
//
// Since cnNursery makes essential use of multi-threading, we need to
// ensure all globals are thread safe. To do this we make liberal use
// of the sync.RWMutexes, one for each global.
//
// In this file we manage the global singleton for configuration.
//

package main

import (
  "encoding/json"
  "github.com/jinzhu/configor"
  "sync"
)

//////////////////////////
// Configuration variables
//

type ConfigType struct {
  Name         string
  Host         string
  Interface    string
  Port         string
  Html_Dir     string
  Base_Url     string
  Primary_Url  string
  Ca_Cert_Path string
  Cert_Path    string
  Key_Path     string
}

var configSync sync.RWMutex
var configPriv ConfigType

//////////////////////////
// Configuration functions

func loadConfiguration(configFileName string) {
  configSync.Lock()
  defer configSync.Unlock()

  configor.Load(&configPriv, configFileName)
}

func getConfig() ConfigType {
  configSync.RLock()
  defer configSync.RUnlock()

  return configPriv
}

func configToJsonBytes() ([]byte, error) {
  configSync.RLock()
  defer configSync.RUnlock()

  return json.MarshalIndent(configPriv, "", "  ")
}
