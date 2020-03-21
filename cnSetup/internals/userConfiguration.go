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

package CNSetup

import (
  "fmt"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "os"
  "strings"
  "sync"
  "text/template"
)

type User struct {
  Name          string
  Ca_Cert_Path  string
  Cert_Path     string
  Key_Path      string
  Primary_Url   string
  Config_Path   string
  Serial_Number int64
}

var (
  UserDefaults = User{
    "", // Name
    "", // Ca_Cert_Path
    "", // Cert_Path
    "", // Key_Path
    "", // Primary_Url
    "", // Config_Path
    0,  // Serial_Number
  }
)

func (user *User) NormalizeConfiguration(
  userNum    int,
  defaults   User,
  primaryUrl string,
) {
  user.Primary_Url = primaryUrl

  if user.Ca_Cert_Path == "" { user.Ca_Cert_Path = defaults.Ca_Cert_Path }
  if user.Cert_Path    == "" { user.Cert_Path    = defaults.Cert_Path }
  if user.Key_Path     == "" { user.Key_Path     = defaults.Key_Path }

  nPathPrefix := "users/"+user.Name +"/"+strings.ReplaceAll(user.Name, ".", "-")

  if user.Ca_Cert_Path == "" { user.Ca_Cert_Path = nPathPrefix+"-ca-crt.pem" }
  if user.Cert_Path    == "" { user.Cert_Path    = nPathPrefix+"-crt.pem" }
  if user.Key_Path     == "" { user.Key_Path     = nPathPrefix+"-key.pem" }
  if user.Config_Path  == "" { user.Config_Path  = nPathPrefix+"-config.yaml" }

  // we need to use DIFFERENT serial numbers for each of CA (1<<32),
  //  C/S  ((1<<5 + nurseryNum)<<33) and
  //  User ((2<<5 + userNum)<<33)
  //
  if user.Serial_Number == 0 { user.Serial_Number = int64(2<<5 + userNum)  }
}

// Write out a user's configuration YAML file which is required for them 
// to use the cnTypeSetter command to type set one or more of their 
// ConTeXt documents. 
//
// We provide the user, the user defaults as well as the primaryUrl of 
// this federation of Nurseries. 
//
// We also provide an optional WaitGroup which, if not nil, is used to 
// allow this function to be called asynchronously as a go routine. 
//
func (user *User) WriteUserConfiguration(
  wg        *sync.WaitGroup,
  log       *logger.LoggerType,
) {
  if wg != nil {
    wg.Add(1)
    defer wg.Done()
  }

  fmt.Printf("\n\nCreating configuration file for the user [%s]\n", user)

  yamlTemplateStr := `
# This is the configuration for the {{.Name}} User
#
# It has been automatically generated by the cnSetup tool
#
# DO NOT EDIT THIS FILE (any changes will be lost)

name:         "{{.Name}}"
primary_url:  "{{.Primary_Url}}"
ca_cert_path: "{{.Ca_Cert_Path}}"
cert_path:    "{{.Cert_Path}}"
key_path:     "{{.Key_Path}}"
`
  yamlTemplate, err := template.New("yamlTemplate").Parse(yamlTemplateStr)
  log.MayBeFatal("Could not parse the yaml template", err)

  yamlFile, err := os.Create(user.Config_Path)
  log.MayBeFatal("Could not open the config file for writing", err)

  err = yamlTemplate.Execute(yamlFile, user)

  err = yamlFile.Close()
  log.MayBeFatal("Could not close the config file", err)
}
