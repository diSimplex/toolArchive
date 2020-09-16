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
  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/cnNursery/internals"
  "github.com/diSimplex/ConTeXtNursery/interfaces/action"
  "github.com/diSimplex/ConTeXtNursery/interfaces/control"
  "github.com/diSimplex/ConTeXtNursery/interfaces/discovery"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "math/rand"
  "runtime"
  "time"
)

var configFileName string
var showConfig     bool
var serverCert     tls.Certificate
var caCertPool     *x509.CertPool

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

  cnLog  := logger.CreateLogger("cnNursery")
  cnLog.SetPrintStack(true)
  
  config := CNNurseries.CreateConfiguration(cnLog)
  config.LoadConfiguration(configFileName, showConfig)

  cnLog.Logf("cnNursery: %s started", config.Name)
  cnLog.Logf("numCPU: %d\n", runtime.NumCPU())
  cnLog.Logf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(-1))
  
  // seed the math/rand random number generator with a "random" seed
  rand.Seed(time.Now().Unix())

  ////////////////////////////////
  // initialize interfaces
  //   BEFORE we start any threads
  cc := clientConnection.CreateClientConnection(
    config.Ca_Cert_Path, config.Cert_Path, config.Key_Path,
    cnLog,
  )

  ws := webserver.CreateWebServer(
    config.Interface, config.Port, `

The cnNursery process provides a RESTful interface to the federation of
ConTeXt Nurseries.

Each ConTeXt Nursery in the federation is capable of managing the type
setting of one or more ConTeXt based (sub)documents in parallel.

`,
    config.Ca_Cert_Path, config.Cert_Path, config.Key_Path,
    cnLog,
  )

  err := ws.AddStaticFileHandlers(
    "static/index.html",
    "static/images/TeddyBear.ico",
    "/static",
    "static",
  )
  cnLog.MayBeError("Could not add static file handlers", err)
  
  cnActions := CNNurseries.CreateActionsState(config, ws, cc)
  action.AddActionInterface(ws, cnActions)
  
  cnInfoMap := CNNurseries.CreateCNInfoMap(config)
  discovery.AddDiscoveryInterface(ws, cnInfoMap)

  cnState := CNNurseries.CreateCNState(config, cnInfoMap, ws, cc)
  control.AddControlInterface(ws, cnState)

  /////////////////////////////////////
  // Start client and webServer threads

  // periodically send out a heart beat message the the federation's
  // primary cnNursery 
  go CNNurseries.SendPeriodicHeartBeats(config, cnState, cnInfoMap, cc)

  // periodically cull Nurseries to which we can no longer connect to
  go CNNurseries.GrimReaper(config, cnInfoMap, cc)
  
  ws.RunWebServer()
}
