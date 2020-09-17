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

package CNSetup

import (
  "encoding/json"
  "fmt"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "github.com/jinzhu/configor"
  "os"
)

// ConfigType contains the configuration for the whole cnSetup command. 
//
//
// Its associated methods are responsible for loading configuration for a 
// Federation of ConTeXt Nurseries, descriptions of all nurseries and 
// users, as well as maintaining internal copies of the CA 
//
// CONSTRAINTS: Once created, the values in this structure SHOULD only be 
// altered by structure methods. 
//
type ConfigType struct {

  // Name of the Federation of ConTeXt Nurseries
  //
  Federation_Name        string `default:"nurseries"`

  // Certificate information
  //
  Key_Size                uint `default:"4096"`
  Certificate_Authority   CAType

  // Nurseries
  //
  Nursery_Defaults         NurseryType
  Nurseries              []NurseryType
  NATS_Federation_Routes []string
  NATS_Message_Routes    []string

  // Users
  //
  User_Defaults            UserType
  Users                  []UserType
  
  // Auxilary fields for logging
  //
  CSLog                   *logger.LoggerType
}

// Create an (empty) configuration structure.
//
// Typically this empty configuration structure will be used to LoadConfiguration.
//
// CREATES config;
//
func CreateConfiguration(csLog *logger.LoggerType) *ConfigType {
  return &ConfigType{CSLog: csLog}
}

// Load and normalize a configuration from the configFileName file.
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
  configor.Load(&config, configFileName)

    
  config.Certificate_Authority.NormalizeCA(config)
  
  if config.Federation_Name == "" {
    config.CSLog.Logf("You MUST specify a Federation Name")
    os.Exit(-1)
  }

  if len(config.Nurseries) < 1 {
    config.CSLog.Logf("You MUST specify at least ONE Nursery")
    os.Exit(-1)
  }

  // locate the primary Nursery and normalize each Nursery structure 
  //
  config.Nursery_Defaults.NormalizeConfig(0, &NurseryDefaults, config)
  config.NATS_Federation_Routes = make([]string, len(config.Nurseries))
  config.NATS_Message_Routes    = make([]string, len(config.Nurseries))
  natsShuffle := make([]int, len(config.Nurseries))
  for i, _ := range config.Nurseries {
    natsShuffle[i] = i
    config.Nurseries[i].NormalizeConfig(i, &config.Nursery_Defaults, config)
    config.NATS_Federation_Routes[i] = config.Nurseries[i].ComputeFederationNATS()
    config.NATS_Message_Routes[i]    = config.Nurseries[i].ComputeMessageNATS()
  }

  // now explicitly reshuffle the NATS routes for each cnNursery
  //
  for i, _ := range config.Nurseries {
    config.Nurseries[i].SetNatsRoutes(
      &config.NATS_Federation_Routes,
      &config.NATS_Message_Routes,
      &natsShuffle,
    )
  }

  config.User_Defaults.NormalizeConfig(
    -1,
    &UserDefaults,
    config,
  )
  for i, _ := range config.Users {
    config.Users[i].NormalizeConfig(
      i,
      &config.User_Defaults,
      config,
    )
  }

  // now explicitly reshuffle the NATS routes for each user
  //
  for i, _ := range config.Users {
    config.Users[i].SetNatsRoutes(
      &config.NATS_Message_Routes,
      &natsShuffle,
    )
  }
    
  if showConfig {
    configStr, _ := json.MarshalIndent(config, "", "  ")
    fmt.Print(string(configStr))
    os.Exit(0)
  }
}
