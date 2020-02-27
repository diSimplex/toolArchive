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

// This code has been inspired by: Shane Utt's excellent article:
//   https://shaneutt.com/blog/golang-ca-and-signed-cert-go/

package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "github.com/jinzhu/configor"
  "log"
  "os"
)

//////////////////////////
// Configuration variables
//
var config = struct {

  Federation_Name string `default:"nurseries"`

  Key_Size uint `default:"4096"`

  Certificate_Authority struct {
    Serial_Number  uint    `default:"1"`
    Organization   string
    Country        string
    Province       string
    Locality       string
    Street_Address string
    Postal_Code    string

    Valid_For struct {
      Years  uint `default:"10"`
      Months uint `default:"0"`
      Days   uint `default:"0"`
    }
  }

  Default_Port uint `default:"0"`

  Nurseries []struct {
    Host    string `required`
    Port    uint   `default:"0"`
    Primary bool   `default:"false"`
  }

  Users []string `required`
}{}

var configFileName string
var showConfig     bool


/////////////////////////////
// Logging and Error handling
//
func configMayBeFatal(logMessage string, err error) {
  if err != nil {
    log.Fatalf("cnConfig(FATAL): %s ERROR: %s\n", logMessage, err)
  }
}

func main() {
  const (
    configFileNameDefault =  "nurseries.yaml"
    configFileNameUsage   =  "The configuration file to load"
    showConfigDefault     =  false
    showConfigUsage       =  "Show the loaded configuration"
  )
  flag.BoolVar(&createCA, "createCA", false, "Should the Certificate Authority be created if the crt and key files can't be loaded?")
  flag.StringVar(&configFileName, "config", configFileNameDefault, configFileNameUsage)
  flag.StringVar(&configFileName, "c", configFileNameDefault, configFileNameUsage)
  flag.BoolVar(&showConfig, "show", showConfigDefault, showConfigUsage)
  flag.BoolVar(&showConfig, "s", showConfigDefault, showConfigUsage)
  flag.Parse()

  configor.Load(&config, configFileName)

  if showConfig {
    configStr, _ := json.MarshalIndent(config, "", "  ")
    fmt.Print(string(configStr))
    os.Exit(0)
  }

  loadCA()
}
