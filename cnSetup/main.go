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
  "bufio"
  "encoding/json"
  "flag"
  "fmt"
  "github.com/jinzhu/configor"
  "log"
  "os"
  "time"
  "strings"
// temp
//  "bytes"
//  "crypto/x509"
//  "encoding/pem"
//  "io/ioutil"
)

//////////////////////////
// Configuration variables
//

type Nursery struct {
  Name         string
  Host         string
  Hosts        []string
  Interface    string
  Port         uint
  Html_Dir     string
  Ca_Cert_Path string
  Cert_Path    string
  Key_Path     string
  Is_Primary   bool
  Base_Url     string
  Primary_Url  string
  Config_Path  string
}

type User struct {
  Name         string
  Ca_Cert_Path string
  Cert_Path    string
  Key_Path     string
  Primary_Url  string
  Config_Path  string
}

var config = struct {

  Federation_Name string `default:"nurseries"`

  Key_Size uint `default:"4096"`

  Certificate_Authority struct {
    Serial_Number  uint
    Organization   string
    Country        string
    Province       string
    Locality       string
    Street_Address string
    Postal_Code    string
    Email_Address  string
    Common_Name    string

    Valid_For struct {
      Years  uint `default:"10"`
      Months uint `default:"0"`
      Days   uint `default:"0"`
    }
  }

  Nursery_Defaults Nursery

  Nurseries []Nursery

  Users []User
}{}

var userPasswords = map[string]string{}

var configFileName string
var showConfig     bool

/////////////////////////////
// Logging and Error handling
//
func setupMayBeFatal(logMessage string, err error) {
  if err != nil {
    log.Fatalf("cnSetup(FATAL): %s ERROR: %s\n", logMessage, err)
  }
}

func main() {
var (
    nurseryDefaults = Nursery{
      "",                       // Name
      "",                       // Host
      []string{},               // Hosts
      "0.0.0.0",                // Interface
      8989,                     // Port
      "/var/www/html",          // Html_Dir
      "",                       // Ca_Cert_Path
      "",                       // Cert_Path
      "",                       // Key_Path
      false,                    // Is_Primary
      "https://localhost:8989", // Base_Url
      "",                       // Primary_Url
      "",                       // Config_Path
    }
    userDefaults = User{
      "", // Name
      "", // Ca_Cert_Path
      "", // Cert_Path
      "", // Key_Path
      "", // Primary_Url
      "", // Config_Path
    }
  )

  const (
    configFileNameDefault =  "nurseries.yaml"
    configFileNameUsage   =  "The configuration file to load"
    showConfigDefault     =  false
    showConfigUsage       =  "Show the loaded configuration"
  )
  flag.BoolVar(&createCA, "createCA", false, "Should the Certificate Authority be created if the crt and key files can't be loaded?")
  flag.StringVar(&configFileName, "config", configFileNameDefault, configFileNameUsage)
  flag.StringVar(&configFileName, "c", configFileNameDefault, configFileNameUsage)
  flag.BoolVar(&showConfig, "show", showConfigDefault, showConfigUsage)
  flag.BoolVar(&showConfig, "s", showConfigDefault, showConfigUsage)
  flag.Parse()

  configor.Load(&config, configFileName)

  // make sure the Serial_Number is constantly increasing...
  //
  if config.Certificate_Authority.Serial_Number == 0 {
    config.Certificate_Authority.Serial_Number = uint(time.Now().Unix())
  }

  if config.Federation_Name == "" {
    config.Federation_Name = "ConTeXt Nurseries"
  }

  if showConfig {
    configStr, _ := json.MarshalIndent(config, "", "  ")
    fmt.Print(string(configStr))
    os.Exit(0)
  }

  loadCA()


  // locate the primary Nursery
  normalizeNurseryConfig(&config.Nursery_Defaults, nurseryDefaults)
  primaryNursery := &config.Nurseries[0]
  for i, _ := range config.Nurseries {
    if config.Nurseries[i].Is_Primary {
      if ! primaryNursery.Is_Primary {
         primaryNursery = &config.Nurseries[i]
      }
    }
    normalizeNurseryConfig(&config.Nurseries[i], config.Nursery_Defaults)
  }
  primaryNurseryUrl := computePrimaryNurseryUrl(primaryNursery)

  // now create each Nursery's certificates as well as configuration
  for i, aNursery := range config.Nurseries {
    createNurseryCertificate(&aNursery, i)
    writeNurseryConfiguration(&aNursery, primaryNurseryUrl)
  }

  // start by loading in the existing user passwords
  passwordFile, err := os.Open("users/passwords")
  if err == nil {
    scanner := bufio.NewScanner(passwordFile)
    scanner.Split(bufio.ScanLines)
    for scanner.Scan() {
      aLine := scanner.Text()
      fields    := strings.Split(aLine, "\t")
      aUser     := fields[0]
      aPassword := fields[1]
      userPasswords[aUser] = aPassword
    }
    passwordFile.Close()
  }

  // now create each User's certificates
  for i, aUser := range config.Users {
    createUserCertificate(aUser.Name, i)
    writeUserConfiguration(aUser, userDefaults, primaryNurseryUrl)
  }

  // now write out the file of user passwords
  passwordFile, err = os.Create("users/passwords")
  setupMayBeFatal("Could not open [users/passwords] file", err)
  for aUser, aPassword := range userPasswords {
    passwordFile.WriteString(aUser+"\t"+aPassword+"\n")
  }
  passwordFile.Close()
  os.Chmod("users/passwords", 0600)
  fmt.Printf("\nThe automatically generated passwords for each user's PKCS#12 file\n")
  fmt.Printf("  can be found in the file [users/passwords]\n\n")
}
