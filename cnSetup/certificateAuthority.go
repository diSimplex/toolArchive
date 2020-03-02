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
  "encoding/pem"
  "fmt"
  "io/ioutil"
  "math/big"
  "os"
  "time"
)

////////////////////////
// Certificate Authority
//
var createCA                 bool
var caDir                 = "ca"
var caCertFileName        = "ca/certificateAuthority.crt"
var caPrivateKeyFileName  = "ca/certificateAuthority.key"
var caCert                  *x509.Certificate
var caPrivateKey            *rsa.PrivateKey

func createCertificateAuthorityFiles() {
  fmt.Print("\nCreating a new Certificate Authority\n")
  lcaCert := &x509.Certificate {
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

  lcaPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
  setupMayBeFatal("could not generate rsa key for CA", err)

  caBytes, err := x509.CreateCertificate(rand.Reader, lcaCert, lcaCert, &lcaPrivateKey.PublicKey, lcaPrivateKey)
  setupMayBeFatal("could not create the CA certificate", err)

  caSubject := "ConTeXt Nursery " + config.Federation_Name + " Certificate Authority"
  caDate    := time.Now().String()

  os.MkdirAll(caDir, 0755)

  caPEM := new(bytes.Buffer)
  pem.Encode(caPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Headers: map[string]string {
      "Subject": caSubject,
      "Date":    caDate,
    },
    Bytes: caBytes,
  })
  err = ioutil.WriteFile(caCertFileName, caPEM.Bytes(), 0644)
  setupMayBeFatal("could not write the certificateAuthority.crt file", err)

  caPrivateKeyPEM := new(bytes.Buffer)
  pem.Encode(caPrivateKeyPEM, &pem.Block {
    Type: "RSA PRIVATE KEY",
    Headers: map[string]string {
      "Subject": caSubject,
      "Date":    caDate,
    },
    Bytes: x509.MarshalPKCS1PrivateKey(lcaPrivateKey),
  })
  err = ioutil.WriteFile(caPrivateKeyFileName, caPrivateKeyPEM.Bytes(), 0644)
  setupMayBeFatal("could not write the certificateAuthority.key file", err)

  // since we have made it this far... both the cert and key are OK...
  // so store the local copies in the global variables...
  caCert       = lcaCert
  caPrivateKey = lcaPrivateKey
}

func loadCA() {
  caCertBytes, err := ioutil.ReadFile(caCertFileName)
  if err != nil {
    if !createCA {
      setupMayBeFatal("could not load the certificate authority's *.crt file; did you want to use the '-createCA' option?", err)
    } else {
      createCertificateAuthorityFiles()
      return
    }
  }

  caCertPEM, _ /*restCaCertBytes*/ := pem.Decode(caCertBytes)
  if caCertPEM == nil || caCertPEM.Type != "CERTIFICATE" {
    if !createCA {
      setupMayBeFatal("could not locate the certificate authority's CERTIFICATE block", err)
    } else {
      createCertificateAuthorityFiles()
      return
    }
  }

  lcaCert, err := x509.ParseCertificate(caCertPEM.Bytes)
  if err != nil {
    if !createCA {
      setupMayBeFatal("could not parse the certificate authority's certificate", err)
    } else {
      createCertificateAuthorityFiles()
      return
    }
  }

  caKeyBytes,  err  := ioutil.ReadFile(caPrivateKeyFileName)
  if err != nil {
    if !createCA {
      setupMayBeFatal("could not load the certificate authority's *.key file", err)
    } else {
      createCertificateAuthorityFiles()
      return
    }
  }

  caKeyPEM, _ /*restCaKeyBytes*/ := pem.Decode(caKeyBytes)
  if caKeyPEM == nil || caKeyPEM.Type != "RSA PRIVATE KEY" {
    if !createCA {
      setupMayBeFatal("could not locate the certificate authority's RSA PRIVATE KEY block", err)
    } else {
      createCertificateAuthorityFiles()
      return
    }
  }

  lcaPrivateKey, err := x509.ParsePKCS1PrivateKey(caKeyPEM.Bytes)
  if err != nil {
    if !createCA {
      setupMayBeFatal("could not parse the certificate authority's private key", err)
    } else {
      createCertificateAuthorityFiles()
      return
    }
  }

  // If we managed to get this far... both the cert and key are OK...
  // so store the local copies in the global variables...
  caCert       = lcaCert
  caPrivateKey = lcaPrivateKey
}
