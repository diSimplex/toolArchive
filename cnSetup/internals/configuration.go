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
  "github.com/jinzhu/configor"
  "os"
  "sync"
  "time"
)

//////////////////////////
// Configuration variables
//

type Configuration struct {
  Mutex           sync.RWMutex
  
  Federation_Name string `default:"nurseries"`

  Key_Size uint `default:"4096"`

  Certificate_Authority struct {
    Serial_Number  uint
    Organization   string
    Country        string
    Province       string
    Locality       string
    Street_Address string
    Postal_Code    string
    Email_Address  string
    Common_Name    string

    Valid_For struct {
      Years  uint `default:"10"`
      Months uint `default:"0"`
      Days   uint `default:"0"`
    }
  }

  Nursery_Defaults    Nursery

  Primary_Nursery    *Nursery
  Primary_Nursery_Url string
  
  Nurseries         []Nursery

  User_Defaults       User
  
  Users             []User
}

//type Config struct {
//  ConfigPriv Configuration
//}

func CreateConfiguration() *Configuration {
  return Configuration{}
}

func (config *Configuration) LoadConfiguration(
  configFileName string,
  showConfig     bool,
) {
  configor.Load(&config, configFileName)
  
    // make sure the Serial_Number is constantly increasing...
  //
  if config.Certificate_Authority.Serial_Number == 0 {
    config.Certificate_Authority.Serial_Number = uint(time.Now().Unix())
  }

  if config.Federation_Name == "" {
    config.Federation_Name = "ConTeXt Nurseries"
  }

  // locate the primary Nursery
  config.Nursery_Defaults.NormalizeNurseryConfig(NurseryDefaults)
  config.Primary_Nursery = &config.Nurseries[0]
  for i, _ := range config.Nurseries {
    if config.Nurseries[i].Is_Primary {
      if ! config.Primary_Nursery.Is_Primary {
         config.Primary_Nursery = &config.Nurseries[i]
      }
    }
    config.Nurseries[i].NormalizeNurseryConfig(config.Nursery_Defaults)
  }
  config.Primary_Nursery_Url = config.Primary_Nursery.ComputeUrl()

  config.User_Defaults.NormalizeConfiguration(
    UserDefaults,
    config.Primary_Nursery_Url,
  )
  for i, _ := range config.Users {
    config.Users[i].NormalizeConfiguration(
      config.User_Defaults,
      config.Primary_Nursery_Url,
    )
  }
  
  if showConfig {
    configStr, _ := json.MarshalIndent(config, "", "  ")
    fmt.Print(string(configStr))
    os.Exit(0)
  }
}
