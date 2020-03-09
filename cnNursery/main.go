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
  "crypto/tls"
  "crypto/x509"
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "math/rand"
  "os"
  "time"
)

var configFileName string
var showConfig     bool
var serverCert     tls.Certificate
var caCertPool     *x509.CertPool

/////////////////////////////
// Logging and Error handling
//
func cnNurseryMayBeFatal(logMessage string, err error) {
  if err != nil {
    log.Fatalf("cnNursery(FATAL): %s ERROR: %s", logMessage, err)
  }
}

func cnNurseryMayBeError(logMessage string, err error) {
  if err != nil {
    log.Printf("cnNursery(error): %s error: %s",logMessage, err)
  }
}

func cnNurseryLog(logMesg string) {
  log.Printf("cnNursery(info): %s", logMesg)
}

func cnNurseryLogf(logFormat string, v ...interface{}) {
  log.Printf("cnNursery(info): "+logFormat, v...)
}

func cnNurseryJson(logMesg string, valName string, aValue interface{}) {
  jsonBytes, err := json.MarshalIndent(aValue, "", "  ")
  if err != nil {
    cnNurseryMayBeError("Could not marshal "+valName+" into json", err)
    jsonBytes = make([]byte, 0)
  }
  log.Printf("cnNursery(json): %s", string(jsonBytes))
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

  loadConfiguration(configFileName)

  if showConfig {
    configBytes, _ := configToJsonBytes()
    fmt.Printf("%s\n", string(configBytes))
    os.Exit(0)
  }

  // seed the math/rand random number generator with a "random" seed
  rand.Seed(time.Now().Unix())

  // load the server and ca certificates for use by all client/servers in
  // this Nursery
  //
  var err error
  lConfig := getConfig()
  serverCert, err = tls.LoadX509KeyPair( lConfig.Cert_Path, lConfig.Key_Path )
  cnNurseryMayBeFatal("Could not load cert/key pair", err)
  //
  caCert, err := ioutil.ReadFile(lConfig.Ca_Cert_Path)
  cnNurseryMayBeFatal("Could not load the CA certificate", err)
  //
  caCertPool = x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caCert)

  go sendPeriodicHeartBeats()

  handleHeartBeats()

  runWebServer()

}
