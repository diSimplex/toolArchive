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
  Nursery_Defaults      NurseryType
  Nurseries           []NurseryType
  Primary_Nursery      *NurseryType
  Primary_Nursery_Url   string

  // Users
  //
  User_Defaults         UserType
  Users               []UserType
  
  // Auxilary fields for access and logging
  //
  Mutex                 sync.RWMutex
  CSLog                *logger.LoggerType
}

// Create an (empty) configuration structure.
//
// Typically this empty configuration structure will be used to LoadConfiguration.
//
func CreateConfiguration(csLog *logger.LoggerType) *ConfigType {
  return &ConfigType{}
}

// Start changing the configuration by obtaining a (Write) lock on the 
// configuration Mutex. 
// 
// This requests a complete lock on the configuration values.
// 
// If you only need to READ these values, consider using the StartReading 
// method instead. 
//
func (config *ConfigType) StartChanging() {
  config.Mutex.Lock()
}

// Stop changing the configuration by releasing the (Write) lock on the 
// configuration Mutex. 
//
func (config *ConfigType) StopChanging() {
  config.Mutex.Unlock()
}

// Start reading the values in the configuration by obtaining a (Read) 
// lock on the configuration Mutex. 
//
// IT IS CRITICAL THAT NO VALUES ARE CHANGED. 
//
// If you need to CHANGE a value, then use the StartChanging method 
// instead. 
//
func (config *ConfigType) StartReading() {
  config.Mutex.RLock()
}

// Stop reading the values in the configuraiton by releasing the (Read) 
// lock on the configuration Mutex. 
//
func (config *ConfigType) StopReading() {
  config.Mutex.RUnlock()
}

// Load and normalize a configuration from the configFileName file.
//
// If showConfig is true, show the normalized configuration and exit.
//
func (config *ConfigType) LoadConfiguration(
  configFileName string,
  showConfig     bool,
) {
  config.Mutex.Lock()
  defer config.Mutex.Unlock()
  
  configor.Load(&config, configFileName)

    
  config.Certificate_Authority.NormalizeCA(config)
  
  if config.Federation_Name == "" {
    config.CSLog.Logf("You MUST specify a Federation Name")
    os.Exit(-1)
  }

  // locate the primary Nursery and normalize each Nursery structure 
  //
  config.Primary_Nursery = &config.Nurseries[0]
  config.Nursery_Defaults.NormalizeConfig(0, &NurseryDefaults, config)
  for i, _ := range config.Nurseries {
    if config.Nurseries[i].Is_Primary {
      if ! config.Primary_Nursery.Is_Primary {
         config.Primary_Nursery = &config.Nurseries[i]
      }
    }
    config.Nurseries[i].NormalizeConfig(i, &config.Nursery_Defaults, config)
  }

  config.Primary_Nursery_Url = config.Primary_Nursery.ComputeUrl()

  // now explicitly set the primary url for each Nursery.
  //
  for i, _ := range config.Nurseries {
    config.Nurseries[i].SetPrimaryUrl(config.Primary_Nursery_Url)
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
