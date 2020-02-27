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
  "bytes"
  "crypto/rand"
  "crypto/rsa"
  "crypto/x509"
  "crypto/x509/pkix"
  "encoding/json"
  "encoding/pem"
  "flag"
  "fmt"
  "github.com/jinzhu/configor"
  "log"
  "math/big"
  "os"
  "time"
)

//////////////////////////
// Configuration variables
//
var config = struct {

  Federation_Name string `default:"nurseries"`

  Key_Size uint `default:"4096"`

  Certificate_Authority struct {
    Serial_Number  uint    `default:"1"`
    Organization   string
    Country        string
    Province       string
    Locality       string
    Street_Address string
    Postal_Code    string

    Valid_For struct {
      Years  uint `default:"10"`
      Months uint `default:"0"`
      Days   uint `default:"0"`
    }
  }

  Default_Port uint `default:"0"`

  Nurseries []struct {
    Host    string `required`
    Port    uint   `default:"0"`
    Primary bool   `default:"false"`
  }

  Users []string `required`
}{}

var configFileName string
var showConfig     bool

/////////////////////////////
// Logging and Error handling
//
func configMayBeFatal(logMessage string, err error) {
  if err != nil {
    log.Fatalf("cnConfig(FATAL): %s ERROR: %s\n", logMessage, err)
  }
}

///////////////////////
// Certificate creation
//
func createCA() {
  fmt.Print("\nCreating the Certificate Authority\n")
  ca := &x509.Certificate {
    SerialNumber: big.NewInt(int64(config.Certificate_Authority.Serial_Number)),
    Subject: pkix.Name {
      Organization:  []string{config.Certificate_Authority.Organization},
      Country:       []string{config.Certificate_Authority.Country},
      Province:      []string{config.Certificate_Authority.Province},
      Locality:      []string{config.Certificate_Authority.Locality},
      StreetAddress: []string{config.Certificate_Authority.Street_Address},
      PostalCode:    []string{config.Certificate_Authority.Postal_Code},
    },
    NotBefore: time.Now(),
    NotAfter:  time.Now().AddDate(int(config.Certificate_Authority.Valid_For.Years),
                                  int(config.Certificate_Authority.Valid_For.Months),
                                  int(config.Certificate_Authority.Valid_For.Days)),
    IsCA:        true,
    ExtKeyUsage: []x509.ExtKeyUsage{
      x509.ExtKeyUsageClientAuth,
      x509.ExtKeyUsageServerAuth,
    },
    KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign ,
    BasicConstraintsValid: true,
  }

  caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
  configMayBeFatal("could not generate rsa key for CA", err)

  caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
  configMayBeFatal("could not create the CA certificate", err)

  caPEM := new(bytes.Buffer)
  pem.Encode(caPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: caBytes,
  })

  caPrivateKeyPEM := new(bytes.Buffer)
  pem.Encode(caPrivateKeyPEM, &pem.Block {
    Type: "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(caPrivateKey),
  })
}


func main() {
  const (
    configFileNameDefault =  "nurseries.yaml"
    configFileNameUsage   =  "The configuration file to load"
    showConfigDefault     =  false
    showConfigUsage       =  "Show the loaded configuration"
  )
  flag.StringVar(&configFileName, "config", configFileNameDefault, configFileNameUsage)
  flag.StringVar(&configFileName, "c", configFileNameDefault, configFileNameUsage)
  flag.BoolVar(&showConfig, "show", showConfigDefault, showConfigUsage)
  flag.BoolVar(&showConfig, "s", showConfigDefault, showConfigUsage)
  flag.Parse()

  configor.Load(&config, configFileName)

  if showConfig {
    configStr, _ := json.MarshalIndent(config, "", "  ")
    fmt.Print(string(configStr))
    os.Exit(0)
  }

  createCA()
}
