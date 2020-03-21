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
  "flag"
  "fmt"

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


var userPasswords = map[string]string{}

var createCA       bool
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
    wg sync.WaitGroup
  

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

  config := LoadConfiguration(configFileName, showConfig)

  wg.Add(1)
  
  loadCA()

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
  
  wg.Add(-1)
  wg.Wait()
}
