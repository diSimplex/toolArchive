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
  "sync"
)

// ConfigType contains the configuration for the whole cnSetup command. 
//
//
// Its associated methods are responsible for loading configuration for a 
// Federation of ConTeXt Nurseries, descriptions of all nurseries and 
// users, as well as maintaining internal copies of the CA 
//
type ConfigType struct {

  // Name of the Federation of ConTeXt Nurseries
  //
  Federation_Name       string `default:"nurseries"`

  // Certificate information
  //
  Key_Size              uint `default:"4096"`
  Certificate_Authority CAType

  // Nurseries
  //
  Nursery_Defaults      Nursery
  Nurseries           []Nursery
  Primary_Nursery      *Nursery
  Primary_Nursery_Url   string

  // Users
  //
  User_Defaults         User
  Users               []User
  
  // Auxilary fields for access and logging
  //
  Mutex                 sync.RWMutex
  CSLog                *logger.LoggerType
}

//type Config struct {
//  ConfigPriv Configuration
//}

func CreateConfiguration(csLog *logger.LoggerType) *ConfigType {
  return &ConfigType{}
}

func (config *ConfigType) LoadConfiguration(
  configFileName string,
  showConfig     bool,
) {
  config.Mutex.Lock()
  defer config.Mutex.Unlock()
  
  configor.Load(&config, configFileName)
  
  config.Certificate_Authority.NormalizeCA(config)
  
  if config.Federation_Name == "" {
    config.csLog.Logf("You MUST specify a Federation Name")
    os.Exit(-1)
  }

  // locate the primary Nursery
  config.Primary_Nursery = &config.Nurseries[0]
  for i, _ := range config.Nurseries {
    if config.Nurseries[i].Is_Primary {
      if ! config.Primary_Nursery.Is_Primary {
         config.Primary_Nursery = &config.Nurseries[i]
      }
    }
  }
  config.Primary_Nursery_Url = config.Primary_Nursery.ComputeUrl()

  // now normalize the Nursery defaults
  config.Nursery_Defaults.NormalizeConfig(0, &NurseryDefaults, config)
  config.Primary_Nursery = &config.Nurseries[0]
  for i, _ := range config.Nurseries {
    if config.Nurseries[i].Is_Primary {
      if ! config.Primary_Nursery.Is_Primary {
         config.Primary_Nursery = &config.Nurseries[i]
      }
    }
    config.Nurseries[i].NormalizeConfig(i, &config.Nursery_Defaults, config)
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
    
  if showConfig {
    configStr, _ := json.MarshalIndent(config, "", "  ")
    fmt.Print(string(configStr))
    os.Exit(0)
  }
}
