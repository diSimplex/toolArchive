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

package main

import (
  "crypto/tls"
  "crypto/x509"
  "flag"
  "fmt"
  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/interfaces/control"
  "github.com/diSimplex/ConTeXtNursery/interfaces/discovery"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "math/rand"
  "os"
  "time"
)

var configFileName string
var showConfig     bool
var serverCert     tls.Certificate
var caCertPool     *x509.CertPool
var cnInfoMap      *CNInfoMap
var cnState        *CNState

var cnLog = logger.CreateLogger("cnNursery")

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

  loadConfiguration(configFileName)

  if showConfig {
    configBytes, _ := configToJsonBytes()
    fmt.Printf("%s\n", string(configBytes))
    os.Exit(0)
  }

  // seed the math/rand random number generator with a "random" seed
  rand.Seed(time.Now().Unix())

  cnLog.SetPrintStack(true)

  ////////////////////////////////
  // initialize interfaces
  //   BEFORE we start any threads
  lConfig := getConfig()

  cc := clientConnection.CreateClientConnection(
    lConfig.Ca_Cert_Path, lConfig.Cert_Path, lConfig.Key_Path,
    cnLog,
  )

  ws := webserver.CreateWebServer(
    lConfig.Interface, lConfig.Port, `

The cnNursery process provides a RESTful interface to the federation of
ConTeXt Nurseries.

Each ConTeXt Nursery in the federation is capable of managing the type
setting of one or more ConTeXt based (sub)documents in parallel.

`,
  lConfig.Ca_Cert_Path, lConfig.Cert_Path, lConfig.Key_Path,
  cnLog,
  )

//  handleControl()

  cnInfoMap = CreateCNInfoMap()
  discovery.AddDiscoveryInterface(ws, cnInfoMap)

  cnState = CreateCNState(ws, cc)
  control.AddControlInterface(ws, cnState)

  /////////////////////////////////////
  // Start client and webServer threads

  go sendPeriodicHeartBeats(cc)

  ws.RunWebServer()

}
