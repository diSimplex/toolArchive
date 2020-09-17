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
  "math/rand"
  "os"
  "path"
  "strings"
  "text/template"
)

// The UserType contains the information required to:
//
//   1. Create x509 Client Certificates as well as associated
//      public/private RSA keys.
//
//   2. Write out the YAML configuration files used by each cnTypeSetter.
//
// CONSTRAINTS: Once created, the values in this structure SHOULD only be 
// altered by structure methods. 
//
type UserType struct {
  Name                  string
  Cert_Dir              string
  Ca_Cert_Path          string
  Ca_Cert_File          string
  Cert_Path             string
  Cert_File             string
  Key_Path              string
  Key_File              string
  Pkcs12_Path           string
  NATS_Message_Routes []string
  Config_Path           string
  Serial_Number         uint
  Key_Size              uint
  Password              string
}

var (
  UserDefaults = UserType{
    "",         // Name
    "",         // Cert_Dir
    "",         // Ca_Cert_Path
    "",         // Ca_Cert_File
    "",         // Cert_Path
    "",         // Cert_File
    "",         // Key_Path
    "",         // Key_File
    "",         // Pkcs12_Path
    []string{}, // Primary_Host
    "",         // Config_Path
    0,          // Serial_Number
    0,          // Key_Size
    "",         // Password
  }
)

// Normalize the fields of a given UserType (using the defaults provided). 
//
// READS defaults;
// READS config;
// ALTERS user;
// NOT THREAD-SAFE;
// CALLED BY: LoadConfiguraiton ONLY;
//
func (user *UserType) NormalizeConfig(
  userNum   int,
  defaults *UserType,
  config   *ConfigType,
) {
  if user.Name == "" && -1 < userNum {
    config.CSLog.Logf("You MUST supply a name for all users!")
    os.Exit(-1)
  }

  if user.Cert_Dir     == "" { user.Cert_Dir     = defaults.Cert_Dir }
  if user.Ca_Cert_Path == "" { user.Ca_Cert_Path = defaults.Ca_Cert_Path }
  if user.Cert_Path    == "" { user.Cert_Path    = defaults.Cert_Path }
  if user.Key_Path     == "" { user.Key_Path     = defaults.Key_Path }
  if user.Pkcs12_Path  == "" { user.Pkcs12_Path  = defaults.Pkcs12_Path }
  
  if -1 < userNum {
    if user.Cert_Dir     == "" { user.Cert_Dir     =  "users/"+user.Name }
  
    nPathPrefix := user.Cert_Dir + "/" + strings.ReplaceAll(user.Name, ".", "-")
    if user.Ca_Cert_Path == "" { user.Ca_Cert_Path = nPathPrefix+"-ca-crt.pem" }
    if user.Ca_Cert_File == "" { user.Ca_Cert_File = path.Base(user.Ca_Cert_Path) }
    if user.Cert_Path    == "" { user.Cert_Path    = nPathPrefix+"-crt.pem" }
    if user.Cert_File    == "" { user.Cert_File    = path.Base(user.Cert_Path) }
    if user.Key_Path     == "" { user.Key_Path     = nPathPrefix+"-key.pem" }
    if user.Key_File     == "" { user.Key_File     = path.Base(user.Key_Path) }
    if user.Pkcs12_Path  == "" { user.Pkcs12_Path  = nPathPrefix+"-pkcs12.p12" }
    if user.Config_Path  == "" { user.Config_Path  = user.Cert_Dir+"/cnTypeSetter.yaml" }
  }

  // we need to use DIFFERENT serial numbers for each of CA (1<<32),
  //  C/S  ((1<<5 + nurseryNum)<<33) and
  //  User ((2<<5 + userNum)<<33)
  //
  if user.Serial_Number == 0 { user.Serial_Number = uint(2<<5 + userNum)  }

  if user.Key_Size == 0 { user.Key_Size = config.Key_Size   }
}


// Set the Nursery's NATS routes (of the whole federation)
//
// ALTERS nursery;
// NOT THREAD-SAFE;
// CALLED BY: LoadConfiguration ONLY;
//
func (user *UserType) SetNatsRoutes(
  natsMessageRoutes    *[]string,
  natsShuffle          *[]int,
) {
  rand.Shuffle(len(*natsShuffle), func(i, j int) {
    (*natsShuffle)[i], (*natsShuffle)[j] =
      (*natsShuffle)[j], (*natsShuffle)[i]
  })
  user.NATS_Message_Routes = make([]string, len(*natsShuffle))
  for i, j := range *natsShuffle {
  	user.NATS_Message_Routes[i] = (*natsMessageRoutes)[j]
  }
}


// Write out a user's configuration YAML file which is required for them 
// to use the cnTypeSetter command to type set one or more of their 
// ConTeXt documents. 
//
// We provide the user, the user defaults as well as the primaryUrl of 
// this federation of Nurseries. 
//
// READS user;
//
func (user *UserType) WriteConfiguration() error {

  fmt.Printf("\n\nCreating configuration file for the user [%s]\n", user.Name)

  yamlTemplateStr := `
# This is the configuration for the {{.Name}} User
#
# It has been automatically generated by the cnSetup tool
#
# DO NOT EDIT THIS FILE (any changes will be lost)

name:         "{{.Name}}"
ca_cert_path: "{{.Ca_Cert_File}}"
cert_path:    "{{.Cert_File}}"
key_path:     "{{.Key_File}}"
nats_routes:{{ range .NATS_Message_Routes }}
  - {{ . }}{{ end }}
`
  yamlTemplate, err := template.New("yamlTemplate").Parse(yamlTemplateStr)
  if err != nil {
    return fmt.Errorf("Could not parse the yaml template: %w", err)
  }

  yamlFile, err := os.Create(user.Config_Path)
  if err != nil {
    return fmt.Errorf("Could not open the config file for writing: %w", err)
  }

  err = yamlTemplate.Execute(yamlFile, user)
  if err != nil {
    yamlFile.Close()
    return fmt.Errorf("Could not run user configuration YAML template: %w", err)
  }
  
  err = yamlFile.Close()
  if err != nil {
    return fmt.Errorf("Could not close the user config file: %w", err)
  }
  
  return nil
}
