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
  "os"
  "strconv"
  "strings"
  "text/template"
)

// The NurseryType contains the information required to:
//
//   1. Create x509 Server Certificates as well as associated 
//      public/private RSA keys.
//
//   2. Write out the YAML configuration files used by each cnNursery. 
//
type NurseryType struct {
  Name          string
  Host          string
  Hosts         []string
  Interface     string
  Port          uint
  Html_Dir      string
  Cert_Dir      string
  Ca_Cert_Path  string
  Cert_Path     string
  Key_Path      string
  Is_Primary    bool
  Base_Url      string
  Primary_Url   string
  Config_Path   string
  Work_Dir      string
  Actions_Dir   string
  Serial_Number int64
  Key_Size      uint
}

var (
  NurseryDefaults = NurseryType{
    "",                       // Name
    "",                       // Host
    []string{},               // Hosts
    "0.0.0.0",                // Interface
    8989,                     // Port
    "/var/www/html",          // Html_Dir
    "",                       // Cert_Dir
    "",                       // Ca_Cert_Path
    "",                       // Cert_Path
    "",                       // Key_Path
    false,                    // Is_Primary
    "https://localhost:8989", // Base_Url
    "",                       // Primary_Url
    "",                       // Config_Path
    "workDir",                // Word_Dir
    "actionsDir",             // Actions_Dir
    0,                        // Serial_Number
    0,                        // Key_Size
  }
)

// Compute the (control) URL associated with a given Nursery.
//
func (nursery *NurseryType) ComputeUrl() string {
  return "https://"+nursery.Hosts[0]+":"+strconv.Itoa(int(nursery.Port))
}

// Normalize the fields of a given NurseryType (using the defaults 
// provided). 
//
func (nursery *NurseryType) NormalizeConfig(
  nurseryNum int,
  defaults  *NurseryType,
  config    *ConfigType,
) {
  if nursery.Interface    == "" { nursery.Interface    = defaults.Interface }
  if nursery.Port         == 0  { nursery.Port         = defaults.Port }
  if nursery.Html_Dir     == "" { nursery.Html_Dir     = defaults.Html_Dir }
  if nursery.Cert_Dir     == "" { nursery.Cert_Dir     = defaults.Cert_Dir }
  if nursery.Ca_Cert_Path == "" { nursery.Ca_Cert_Path = defaults.Ca_Cert_Path }
  if nursery.Cert_Path    == "" { nursery.Cert_Path    = defaults.Cert_Path }
  if nursery.Key_Path     == "" { nursery.Key_Path     = defaults.Key_Path }
  if nursery.Work_Dir     == "" { nursery.Work_Dir     = defaults.Work_Dir }
  if nursery.Actions_Dir  == "" { nursery.Actions_Dir  = defaults.Actions_Dir }
  if nursery.Key_Size     == 0  { nursery.Key_Size     = defaults.Key_Size }
  
  if nursery.Host == "" { nursery.Host = defaults.Host }
  if nursery.Host != "" {
    hosts := strings.Split(nursery.Host, ",")
    for _, aString := range hosts {
      nursery.Hosts = append(nursery.Hosts, strings.TrimSpace(aString))
    }
    if nursery.Name == "" { nursery.Name = nursery.Hosts[0] }
    if nursery.Cert_Dir     == "" {
      nursery.Cert_Dir = "servers/"+nursery.Name
    }
    nPathPrefix := nursery.Cert_Dir + "/" + nursery.Name
    if nursery.Ca_Cert_Path == "" { nursery.Ca_Cert_Path = nPathPrefix+"-ca-crt.pem" }
    if nursery.Cert_Path    == "" { nursery.Cert_Path    = nPathPrefix+"-crt.pem" }
    if nursery.Key_Path     == "" { nursery.Key_Path     = nPathPrefix+"-key.pem" }
    if nursery.Config_Path  == "" { nursery.Config_Path  = nPathPrefix+"-config.yaml" }
    if nursery.Base_Url     == "" { 
      nursery.Base_Url  = nursery.ComputeUrl()
    }
  }
  
  // we need to use DIFFERENT serial numbers for each of CA (1<<32),
  //  C/S  ((1<<5 + nurseryNum)<<33) and
  //  User ((2<<5 + userNum)<<33)
  //
  if nursery.Serial_Number == 0 {
    nursery.Serial_Number = int64(1<<5 + nurseryNum)
  }
  if nursery.Key_Size == 0 {
    nursery.Key_Size = config.Key_Size
  }
}

// Set the Nursery's Primary Url (of the whole federation)
//
func (nursery *NurseryType) SetPrimaryUrl(primaryUrl string) {
  nursery.Primary_Url = primaryUrl
}

// Write out the YAML configuration file requred to run a given cnNursery 
// command. 
//
func (nursery *NurseryType) WriteConfiguration(
  config *ConfigType,
) {

  fmt.Printf("\n\nCreating configuration for the [%s] Nursery\n", nursery.Name)

  yamlTemplateStr := `
# This is the configuration for the {{.Name}} Nursery
#
# It has been automatically generated by the cnSetup tool
#
# DO NOT EDIT THIS FILE (any changes will be lost)

name:         "{{.Name}}"
host:         "{{.Host}}"
interface:    "{{.Interface}}"
port:          {{.Port}}
html_dir:     "{{.Html_Dir}}"
base_url:     "{{.Base_Url}}"
primary_url:  "{{.Primary_Url}}"
ca_cert_path: "{{.Ca_Cert_Path}}"
cert_path:    "{{.Cert_Path}}"
key_path:     "{{.Key_Path}}"
work_dir:     "{{.Work_Dir}}"
actions_dir:  "{{.Actions_Dir}}"
`

  yamlTemplate, err := template.New("yamlTemplate").Parse(yamlTemplateStr)
  config.CSLog.MayBeFatal("Could not parse the yaml template", err)

  yamlFile, err := os.Create(nursery.Config_Path)
  config.CSLog.MayBeFatal("Could not open the config file for writing", err)

  err = yamlTemplate.Execute(yamlFile, nursery)

  err = yamlFile.Close()
  config.CSLog.MayBeFatal("Could not close the config file", err)
}
