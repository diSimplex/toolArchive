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

//go:generate ./buildBrowserApp

//go:generate esc -o browserApp.go ./browserApp/static

package main

import (
  crand "crypto/rand"
  "encoding/binary"
  "flag"
  "github.com/diSimplex/ConTeXtNursery/cnTypeSetter/internals"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "github.com/diSimplex/ConTeXtNursery/natsServer"
  "github.com/diSimplex/ConTeXtNursery/webserver"
  "log"
  "math/rand"
  "os"
  "runtime"
)

//////////////////////////
// Configuration variables
//
var configFileName string
var configDir      string
var showConfig     bool
var browserAppDir  string

/////////////////////////////
// Logging and Error handling
//
func typeSetterMayBeFatal(logMessage string, err error) {
  if err != nil {
    log.Fatalf("cnTypeSetter(FATAL): %s ERROR: %s\n", logMessage, err)
  }
}

func main() {
  const (
    browserAppDirDefault  =  ""
    browserAppDirUsage    = "An on disk directory in which to find the browser application "
    configDirDefault      =  "/.config/cnNursery"
    configDirUsage        =  "The configuration directory"
    configFileNameDefault =  "cnTypeSetter.yaml"
    configFileNameUsage   =  "The configuration file to load"
    showConfigDefault     =  false
    showConfigUsage       =  "Show the loaded configuration"
  )
  flag.StringVar(&browserAppDir,  "browserApp", browserAppDirDefault,               browserAppDirUsage)
  flag.StringVar(&browserAppDir,  "b",          browserAppDirDefault,               browserAppDirUsage)
  flag.StringVar(&configFileName, "config",     configFileNameDefault,              configFileNameUsage)
  flag.StringVar(&configFileName, "c",          configFileNameDefault,              configFileNameUsage)
  flag.StringVar(&configDir,      "dir",        os.Getenv("HOME")+configDirDefault, configDirUsage)
  flag.StringVar(&configDir,      "d",          os.Getenv("HOME")+configDirDefault, configDirUsage)
  flag.BoolVar(&showConfig,       "show",       showConfigDefault,                  showConfigUsage)
  flag.BoolVar(&showConfig,       "s",          showConfigDefault,                  showConfigUsage)
  flag.Parse()

  cnLog := logger.CreateLogger("cnTypeSetter")
  cnLog.SetPrintStack(true)

  // seed the math/rand random number generator with a "random" seed
  // see: https://stackoverflow.com/a/54491783
  // (random is used when loading the configuration to ensure the NATS 
  // routes are randomized) 
  //
  var randomSeed [8]byte
  _, err := crand.Read(randomSeed[:])
  cnLog.MayBeFatal("Could not read a random seed from the system random stream", err)
  rand.Seed(int64(binary.LittleEndian.Uint64(randomSeed[:])))

  config := CNTypeSetter.CreateConfiguration(cnLog)
  config.LoadConfiguration(configDir, configFileName, browserAppDir, showConfig)

  cnLog.Logf("cnTypeSetter: %s started", config.Name)
  cnLog.Logf("numCPU: %d\n", runtime.NumCPU())
  cnLog.Logf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(-1))

  _ = natsServer.ConnectServer(config.Nats_Routes, cnLog)

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

  err = ws.AddStaticFileHandlers(
    config.Browser_App_Dir + "static/index.html",
    config.Browser_App_Dir + "static/images/TeddyBear.ico",
    "/static",
    config.Browser_App_Dir + "static",
    FSMustByte,
  )
  cnLog.MayBeError("Could not add static file handlers", err)
  
//  ws.RunWebServer()

}
