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
  "encoding/json"
  "fmt"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "github.com/jinzhu/configor"
  "os"
)

// ConfigType contains the configuration for the whole cnNursery command.
//
// CONSTRAINTS: Once created, the values in this structure SHOULD only be 
// altered by structure methods.
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
  Work_Dir     string
  Actions_Dir  string
  CNLog       *logger.LoggerType
}

// Create an (empty) configuration structure
//
// Typically this empty configuration structure will be used to 
// LoadConfiguration. 
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
  configFileName string,
  showConfig     bool,
) {
  configor.Load(config, configFileName)
  
  if showConfig {
    configBytes, _ := json.MarshalIndent(config, "", "  ")
    fmt.Printf("%s\n", string(configBytes))
    os.Exit(0)
  }
}

// Returns true if this cnNursery is configured as the Primary cnNursery 
// of this federation. 
//
// READS config;
//
func (config *ConfigType) IsPrimary() bool {
  return config.Base_Url == config.Primary_Url
}
