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

package CNTypeSetter

import (
  "encoding/json"
  "fmt"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "github.com/jinzhu/configor"
  "os"
  "strings"
)

// The TypeSetter configuration
//
type ConfigType struct {
  Name            string
  Interface       string
  Port            string
  Config_Dir      string
  Browser_App_Dir string
  Ca_Cert_Path    string
  Cert_Path       string
  Key_Path        string
  Nats_Routes   []string

  // Auxilary fields for logging
  //
  CNLog                *logger.LoggerType

}

// Create an (empty) configuration structure.
//
// Typically this empty configuration structure will be used to LoadConfiguration.
//
// CREATES config;
//
func CreateConfiguration(cnLog *logger.LoggerType) *ConfigType {
  return &ConfigType{CNLog: cnLog}
}

// Load and normalize a configuration from the configFileName.
//
// If showConfig is true, show the normalized configuration and exit.
//
// ALTERS config;
// NOT THREAD-SAFE;
// USES various NormalizeXXX methods;
//
func (config *ConfigType) LoadConfiguration(
  configDir      string,
  configFilePath string,
  browserAppDir  string,
  showConfig     bool,
) {
  if ! strings.HasSuffix(configDir, "/") {
  	configDir = configDir + "/"
  }

  config.Config_Dir = configDir
  
  if ! strings.HasPrefix(configFilePath, "/") && 
     ! strings.HasPrefix(configFilePath, ".") {
  	configFilePath = configDir + configFilePath
  }

  configor.Load(config, configFilePath)

  
  if ! strings.HasPrefix(config.Ca_Cert_Path, "/") {
  	config.Ca_Cert_Path = configDir + config.Ca_Cert_Path
  }
  if ! strings.HasPrefix(config.Cert_Path, "/") {
  	config.Cert_Path = configDir + config.Cert_Path
  }
  if ! strings.HasPrefix(config.Key_Path, "/") {
  	config.Key_Path = configDir + config.Key_Path
  }

  
  if browserAppDir != "" { config.Browser_App_Dir = browserAppDir }
  if config.Browser_App_Dir != "" &&
     ! strings.HasPrefix(config.Browser_App_Dir, "/") && 
     ! strings.HasPrefix(config.Browser_App_Dir, ".") {
    config.Browser_App_Dir = configDir + config.Browser_App_Dir
  }
  if config.Browser_App_Dir != "" &&
     ! strings.HasSuffix(config.Browser_App_Dir, "/") {
    config.Browser_App_Dir = config.Browser_App_Dir + "/"
  }

  if config.Interface == "" { config.Interface = "0.0.0.0" }
  if config.Port      == "" { config.Port      = "4224"}

  if showConfig {
    configBytes, _ := json.MarshalIndent(config, "", "  ")
    fmt.Printf("%s\n", string(configBytes))
    os.Exit(0)
  }
}
