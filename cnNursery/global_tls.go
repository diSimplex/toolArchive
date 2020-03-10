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

// This file collects all of the globals required for the cnNursery
//  process.
//
// Since cnNursery makes essential use of multi-threading, we need to
// ensure all globals are thread safe. To do this we make liberal use
// of the sync.RWMutexes, one for each global.
//
// In this file we manage the global singleton for configuration.
//

package main

import (
  "crypto/tls"
  "crypto/x509"
  "io/ioutil"
)

///////////////////////////////////////////////
// Transport Layer Security (TLS) configuration

var tlsConfig *tls.Config

// This MUST BE CALLED ONLY ONCE BEFORE any threads are started
//   (i.e. before any client threads as well as the webserver)
//
func initializeTLS() {

  // load the server and ca certificates for use by all client/servers in
  // this Nursery
  //
  var err error
  lConfig := getConfig()
  serverCert, err := tls.LoadX509KeyPair( lConfig.Cert_Path, lConfig.Key_Path )
  cnNurseryMayBeFatal("Could not load cert/key pair", err)
  //
  caCert, err := ioutil.ReadFile(lConfig.Ca_Cert_Path)
  cnNurseryMayBeFatal("Could not load the CA certificate", err)
  //
  caCertPool := x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caCert)

  // Setup HTTPS client
  tlsConfig = &tls.Config{
    ClientAuth:     tls.RequireAndVerifyClientCert,
    Certificates: []tls.Certificate{serverCert},
    RootCAs:        caCertPool,
    ClientCAs:      caCertPool,
  }

}
