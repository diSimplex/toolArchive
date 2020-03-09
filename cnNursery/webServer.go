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
//  "crypto/x509"
  "encoding/json"
//  "io/ioutil"
//  "log"
//  "net"
  "net/http"
  "strconv"
  "strings"
)

func repliedInJson(w http.ResponseWriter, r *http.Request, value interface{}) bool {
  //
  // determine if we are replying in JSON
  //
  replyInJson := false
  for _, anAcceptValue := range r.Header["Accept"] {
    if strings.Contains(strings.ToLower(anAcceptValue), "json") {
      replyInJson = true
      break
    }
  }

  if replyInJson {
    jsonBytes, err := json.Marshal(value)
    if err != nil {
      cnNurseryMayBeError("Could not json.marshal value in repliedInJson", err)
      jsonBytes = []byte{}
    }
    w.Write(jsonBytes)
  }
  return replyInJson
}

func runWebServer() {

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    cnNurseryLogf("url: [%s]", r.URL.Path)
    w.Write([]byte("Hello from the webServer!"))
  })

  // Setup HTTPS client
  tlsConfig := &tls.Config{
    ClientAuth:     tls.RequireAndVerifyClientCert,
    Certificates: []tls.Certificate{serverCert},
    RootCAs:        caCertPool,
    ClientCAs:      caCertPool,
  }

  lConfig := getConfig()
  hostPort := lConfig.Interface + ":" + strconv.Itoa(int(lConfig.Port))

  cnNurseryLogf("listening at [%s]\n", hostPort)
  listener, err := tls.Listen("tcp",  hostPort, tlsConfig)
  cnNurseryMayBeFatal("Could not create listener", err)

  server := &http.Server{TLSConfig: tlsConfig }
  server.Serve(listener)
}
