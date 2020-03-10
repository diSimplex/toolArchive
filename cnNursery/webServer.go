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
  "encoding/json"
  "html/template"
  "net/http"
  "strconv"
  "strings"
)

////////////////////////////////////////////////////////////////////////
// Web Server

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

  handleBasePage()

//  cnNurseryJson("tlsConfig.Certificates ", "tlsConfig,Certficates", tlsConfig.Certificates)
//  cnNurseryJson("tlsConfig.RootCAs ", "tlsConfig,RootCAs", tlsConfig.RootCAs)
//  cnNurseryJson("tlsConfig.ClientCAs ", "tlsConfig,ClientCAs", tlsConfig.ClientCAs)

  lConfig := getConfig()
  hostPort := lConfig.Interface + ":" + strconv.Itoa(int(lConfig.Port))

  cnNurseryLogf("listening at [%s]\n", hostPort)
  listener, err := tls.Listen("tcp",  hostPort, tlsConfig)
  cnNurseryMayBeFatal("Could not create listener", err)

  server := &http.Server{TLSConfig: tlsConfig }
  server.Serve(listener)
}

////////////////////////////////////////////////////////////////////////
// Base page
//

type BasePageMapType map[string]string

var basePageMap = make(BasePageMapType,0)

func addBasePagePath(path string, description string) {
  if path != "" {
    basePageMap[path] = description
  }
}

type BasePageDataType struct {
  Name        string
  BasePageMap BasePageMapType
}

func handleBasePage() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    lConfig := getConfig()

    cnNurseryLogf("url: [%s] method: [%s]", r.URL.Path, r.Method)

    if r.Method == http.MethodGet {

      // we are replying to a (human) browser

      bpTemplate := basePageTemplate()
      err := bpTemplate.Execute(w, BasePageDataType{
        Name: lConfig.Name,
        BasePageMap: basePageMap,
      })
      if err != nil {
        cnNurseryMayBeError("Could not execute base page template", err)
        w.Write([]byte("Could not provide any ConTeXt Nursery information\nPlease try again!"))
      }
      return
    }
    http.Error(w, "Incorrect Request Method for /", http.StatusMethodNotAllowed)
  })

}

func basePageTemplate() *template.Template {
  basePageTemplateStr := `
  <body>
    <h1>ConTeXt Nursery on {{.Name}}</h1>
    <ul>
{{ range $path, $desc := .BasePageMap }}
      <li>
        <strong><a href="/{{$path}}">{{$path}}</a></strong>
        <p>{{$desc}}</p>
      </li>
{{ end }}
    </ul>
 </body>

`

  theTemplate := template.New("body")

  theTemplate, err := theTemplate.Parse(basePageTemplateStr)
  cnNurseryMayBeFatal("Could not parse the internal base page template", err)

  return theTemplate
}
