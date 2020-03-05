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
  "github.com/jinzhu/configor"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "time"
)

//////////////////////////
// Configuration variables
//

var config = struct {
  Name         string
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
func typeSetterMayBeFatal(logMessage string, err error) {
  if err != nil {
    log.Fatalf("cnTypeSetter(FATAL): %s ERROR: %s\n", logMessage, err)
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

  cert, err := tls.LoadX509KeyPair(config.Cert_Path, config.Key_Path)
  typeSetterMayBeFatal("Could not load cert/key pair", err)

  caCert, err := ioutil.ReadFile(config.Ca_Cert_Path)
  typeSetterMayBeFatal("Could not load the CA certificate", err)

  caCertPool := x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caCert)

  // Setup HTTPS client
  tlsConfig := &tls.Config{
    ClientAuth:     tls.RequireAndVerifyClientCert,
    Certificates: []tls.Certificate{cert},
    RootCAs:        caCertPool,
    ClientCAs:      caCertPool,
  }

  transport := &http.Transport{
    TLSClientConfig:    tlsConfig,
    ForceAttemptHTTP2:  true,
    MaxIdleConns:       10,
    IdleConnTimeout:    30 * time.Second,
    DisableCompression: true,
  }

  client := &http.Client{
    Transport: transport,
  }

  resp, err := client.Get(config.Primary_Url)
  typeSetterMayBeFatal("Could not access the primary Nursery", err)
  defer resp.Body.Close()
  
  respBody, err := ioutil.ReadAll(resp.Body)
  typeSetterMayBeFatal("Could not read the body of the response", err)

  fmt.Printf("%s\n", string(respBody))

}
