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
  "math/rand"
  "os"
  "time"
)

//////////////////////////
// Configuration variables
//

var config = struct {
  Name         string
  Host         string
  Interface    string
  Port         uint
  Html_Dir     string
  Base_Url     string
  Primary_Url  string
  Ca_Cert_Path string
  Cert_Path    string
  Key_Path     string
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
    configFileNameDefault =  "nursery.yaml"
    configFileNameUsage   =  "The configuration file to load"
    showConfigDefault     =  false
    showConfigUsage       =  "Show the loaded configuration"
  )
  flag.StringVar(&configFileName, "config", configFileNameDefault, configFileNameUsage)
  flag.StringVar(&configFileName, "c", configFileNameDefault, configFileNameUsage)
  flag.BoolVar(&showConfig, "show", showConfigDefault, showConfigUsage)
  flag.BoolVar(&showConfig, "s", showConfigDefault, showConfigUsage)
  flag.Parse()

  configor.Load(&config, configFileName)

  if showConfig {
    configStr, _ := json.MarshalIndent(config, "", "  ")
    fmt.Printf("%s\n", string(configStr))
    os.Exit(0)
  }

  // seed the math/rand random number generator with a "random" seed
  rand.Seed(time.Now().Unix())

  go sendPeriodicHeartBeats()

  runWebServer()

}
