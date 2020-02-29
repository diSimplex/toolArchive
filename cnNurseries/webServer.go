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
  "log"
  "net/http"
  "strconv"
)

func WebserverMaybeFatal(logMessage string, err error) {
  if err != nil {
    log.Fatalf("Webserver(FATAL): %s ERROR: %s", logMessage, err)
  }
}

func WebserverMaybeError(logMessage string, err error) {
  if err != nil {
    log.Printf("Webserver(error): %s error: %s",logMessage, err)
  }
}

func WebserverLog(logMesg string) {
  log.Printf("Webserver(info): %s", logMesg)
}

func WebserverLogf(logFormat string, v ...interface{}) {
  log.Printf("Webserver(info): "+logFormat, v...)
}

func runWebServer() {

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    WebserverLogf("url: [%s]", r.URL.Path)
    w.Write([]byte("Hello from the webServer!"))
  })

  config.Port = uint(8989)
  hostPort := config.Host + ":" + strconv.Itoa(int(config.Port))
  WebserverLogf("listening on [%s]\n", hostPort)
  http.ListenAndServe(hostPort, nil)
}
