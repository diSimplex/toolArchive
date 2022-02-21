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
  crand "crypto/rand"
  "encoding/binary"
  "flag"
  "fmt"
  "github.com/diSimplex/ConTeXtNursery/cnSetup/internals"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "os"
  "math/rand"
  "runtime"
  "strings"
  "sync"
)

// A User.Name->User.Password mapping used to write out the 
// "users/password" file. 
//
var userPasswords = map[string]string{}

// Flags for use by the cnSetup commands command line options.
//
var createCA       bool
var configFileName string
var showConfig     bool

// Flag descriptions and defaults as used by the cnSetup command line 
// options. 
//
const (
  configFileNameDefault =  "nurseries.yaml"
  configFileNameUsage   =  "The configuration file to load"
  showConfigDefault     =  false
  showConfigUsage       =  "Show the loaded configuration"
)

func WorkOnNursery(
  i         int,
  aNursery *CNSetup.NurseryType,
  ca       *CNSetup.CAType,
  config   *CNSetup.ConfigType,
  wg       *sync.WaitGroup,
) {
  defer wg.Done()
  
  config.CSLog.DebugLockf("(%d)started on nursery: [%s]\n", i, aNursery.Name)
  err := aNursery.CreateNurseryCertificateToFiles(i, ca, config.Federation_Name)
  config.CSLog.MayBeErrorf(
    err,
    "Could not create nurseryCertificate for [%s]",
    aNursery.Name,
  )
  if err == nil {
    err = aNursery.WriteConfiguration()
    config.CSLog.MayBeErrorf(
      err,
      "Could not write nursery [%s] configuration file", 
      aNursery.Name, 
    )
  }
  if err == nil {
    err = aNursery.WriteNATSConfiguration()
    config.CSLog.MayBeErrorf(
      err,
      "Could not write cnMessages(NATS) nursery [%s] configuration file", 
      aNursery.Name, 
    )
  }
  if err == nil {
    err = aNursery.WriteENVSConfiguration()
    config.CSLog.MayBeErrorf(
      err,
      "Could not write pod environment variables file for the [%s] nursery", 
      aNursery.Name, 
    )
  }
  config.CSLog.DebugLockf("(%d)finished on nursery: [%s]\n", i, aNursery.Name)
}

func WorkOnUser(
  i       int,
  aUser  *CNSetup.UserType,
  ca     *CNSetup.CAType,
  config *CNSetup.ConfigType,
  wg     *sync.WaitGroup,
) {
  defer wg.Done()
  
  config.CSLog.DebugLockf("(%d)started on user: [%s]\n", i, aUser.Name)
  aUser.Password = userPasswords[aUser.Name]
  err := aUser.CreateUserCertificate(i, ca, config.Federation_Name) 
  config.CSLog.MayBeErrorf(
    err,
    "Could not create userCertificate for [%s]",
    aUser.Name,
  )
  if err == nil {
    err = aUser.WriteConfiguration()
    config.CSLog.MayBeErrorf(
      err,
      "Could not write user [%s] configuration file",
      aUser.Name,
    )
  }
  config.CSLog.DebugLockf("(%d)finished on user: [%s]\n", i, aUser.Name)
}

// Orchestrate the (optional) (re)creation of a (self-signed) Certificate 
// Authority, as well as Certificates and Configuration for each Nursery 
// and User. 
//
// After (optionally) (re)creating the Certificate Authority, we use 
// sync.WaitGroups to allow the creation of the Certificates (and 
// configuration) for each Nursery and User to occur in parallel. 
//
func main() {
  var (
    wg sync.WaitGroup
  )

  // Setup the command line options (using the GoLang flag package)
  //
  flag.BoolVar(&createCA, "createCA", false, "Should the Certificate Authority be created if the crt and key files can't be loaded?")
  flag.StringVar(&configFileName, "config", configFileNameDefault, configFileNameUsage)
  flag.StringVar(&configFileName, "c", configFileNameDefault, configFileNameUsage)
  flag.BoolVar(&showConfig, "show", showConfigDefault, showConfigUsage)
  flag.BoolVar(&showConfig, "s", showConfigDefault, showConfigUsage)
  flag.Parse()

  // Setup logging and load the configuration.
  //
  csLog  := logger.CreateLogger("cnSetup")

  // seed the math/rand random number generator with a "random" seed
  // see: https://stackoverflow.com/a/54491783
  // (random is used when loading the configuration to ensure the NATS 
  // routes are randomized) 
  //
  var randomSeed [8]byte
  _, err := crand.Read(randomSeed[:])
  csLog.MayBeFatal("Could not read a random seed from the system random stream", err)
  rand.Seed(int64(binary.LittleEndian.Uint64(randomSeed[:])))
  config := CNSetup.CreateConfiguration(csLog)
  config.LoadConfiguration(configFileName, showConfig)
  
  // load the configuration.
  //
  fmt.Printf("numCPU: %d\n", runtime.NumCPU())
  fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(-1))
  
  
  // Load or (re)Create the CA...
  // (this MUST be done synchronously)
  //
  ca  := CNSetup.CreateCA(config)
  err  = ca.LoadCAFromFiles()
  if err != nil {
    if createCA {
      err = ca.CreateNewCA()
      if err != nil {
        csLog.MayBeFatal("Could not create a new CA", err)
      } else {
        err = ca.WriteCAFiles(config)
        csLog.MayBeFatal("Could not write CA files", err)
      }
    } else {
      csLog.MayBeFatal("Could not load existing CA from files\n\tDid you mean to use the -createCA command line switch?\n", err)
    }
  }

  // The creation of the Nursery and User certificates and configuration
  // can take place asynchronously...
  //
  // We asynchronously create each Nursery's certificates as well as 
  // configuration 
  //
  for i, aNursery := range config.Nurseries {
    fmt.Printf("(%d)working on nursery: [%s]\n", i, aNursery.Name)
    wg.Add(1)
    go WorkOnNursery(i, &config.Nurseries[i], ca, config, &wg)
  }

  // Now deal with the users...
  //
  // ... start by loading in the existing user passwords
  //
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

  // Now create each User's certificates and cnTypeSetter configuration
  //
  for i, aUser := range config.Users {
    fmt.Printf("(%d)working on user: [%s]\n", i, aUser.Name)
    wg.Add(1)
    go WorkOnUser(i, &config.Users[i], ca, config, &wg)
  } 
  
  // Wait for all go routines
  wg.Wait()

  // Now write out the file of user passwords
  //
  passwordFile, err = os.Create("users/passwords")
  csLog.MayBeFatal("Could not open [users/passwords] file", err)
  for _, aUser := range config.Users {
    passwordFile.WriteString(aUser.Name+"\t"+aUser.Password+"\n")
  }
  passwordFile.Close()
  os.Chmod("users/passwords", 0600)
  fmt.Printf("\nThe automatically generated passwords for each user's PKCS#12 file\n")
  fmt.Printf("  can be found in the file [users/passwords]\n\n")
}
