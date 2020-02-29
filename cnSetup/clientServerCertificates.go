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
  "log"
  "os"
  "strings"
  "time"
)

/////////////////////////////
// Client Server Certificates

func createNurseryCertificate(theNursery Nursery, nurseryNum int) {
  if theNursery.Host == "" {
    log.Printf("cnConfig(WARNING): no host names specified for a Nursery, skipping Nursery[%d]\n", nurseryNum)
    return
  }
  hosts := strings.Split(theNursery.Host, ",")
  for i, aString := range hosts {
    hosts[i] = strings.TrimSpace(aString)
  }
  fmt.Printf("\nCreating configuration for the [%s] Nursery\n", hosts[0])

  nCert := &x509.Certificate {
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
    ExtKeyUsage: []x509.ExtKeyUsage{
      x509.ExtKeyUsageClientAuth,
      x509.ExtKeyUsageServerAuth,
    },
    SubjectKeyId: []byte{1,2,3,4,6},
    KeyUsage:    x509.KeyUsageDigitalSignature,
  }

  nPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
  configMayBeFatal("could not generate rsa key for ["+hosts[0]+"] Nursery", err)

  nBytes, err := x509.CreateCertificate(rand.Reader, nCert, caCert, &nPrivateKey.PublicKey, caPrivateKey)
  configMayBeFatal("could not create the certificate for ["+hosts[0]+"] Nursery", err)

  nSubject := "ConTeXt Nursery " + config.Federation_Name + " Server Certificate for ["+hosts[0]+"] Nursery"
  nDate    := time.Now().String()

  nDir := "servers/"+hosts[0]
  os.MkdirAll(nDir, 0755)

  nPEM := new(bytes.Buffer)
  pem.Encode(nPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Headers: map[string]string {
      "Subject": nSubject,
      "Date":    nDate,
    },
    Bytes: nBytes,
  })
  nCertificateFileName := nDir+"/"+hosts[0]+".crt"
  err = ioutil.WriteFile(nCertificateFileName, nPEM.Bytes(), 0644)
  configMayBeFatal("could not write the ["+nCertificateFileName+"] file", err)

  nPrivateKeyPEM := new(bytes.Buffer)
  pem.Encode(nPrivateKeyPEM, &pem.Block {
    Type: "RSA PRIVATE KEY",
    Headers: map[string]string {
      "Subject": nSubject,
      "Date":    nDate,
    },
    Bytes: x509.MarshalPKCS1PrivateKey(nPrivateKey),
  })
  nPrivateKeyFileName := nDir+"/"+hosts[0]+".key"
  err = ioutil.WriteFile(nPrivateKeyFileName, nPrivateKeyPEM.Bytes(), 0644)
  configMayBeFatal("could not write the ["+nPrivateKeyFileName+"] file", err)
}
