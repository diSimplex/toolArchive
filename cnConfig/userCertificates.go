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
  "time"
)

////////////////////
// User Certificates

func createUserCertificate(theUser string, userNum int) {
  if theUser == "" {
    log.Printf("cnConfig(WARNING): no user name specified for a user, skipping user[%d]\n", userNum)
    return
  }
  fmt.Printf("\nCreating configuration for the user [%s]\n", theUser)

  uCert := &x509.Certificate {
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
    },
    SubjectKeyId: []byte{1,2,3,4,6},
    KeyUsage:    x509.KeyUsageDigitalSignature,
  }

  uPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
  configMayBeFatal("could not generate rsa key for user ["+theUser+"]", err)

  uBytes, err := x509.CreateCertificate(rand.Reader, uCert, caCert, &uPrivateKey.PublicKey, caPrivateKey)
  configMayBeFatal("could not create the certificate for user ["+theUser+"]", err)

  uSubject := "ConTeXt Nursery " + config.Federation_Name + " User Certificate for user ["+theUser+"]"
  uDate    := time.Now().String()

  uDir := "users/"+theUser
  os.MkdirAll(uDir, 0755)

  uPEM := new(bytes.Buffer)
  pem.Encode(uPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Headers: map[string]string {
      "Subject": uSubject,
      "Date":    uDate,
    },
    Bytes: uBytes,
  })
  uCertificateFileName := uDir+"/"+theUser+".crt"
  err = ioutil.WriteFile(uCertificateFileName, uPEM.Bytes(), 0644)
  configMayBeFatal("could not write the ["+uCertificateFileName+"] file", err)

  uPrivateKeyPEM := new(bytes.Buffer)
  pem.Encode(uPrivateKeyPEM, &pem.Block {
    Type: "RSA PRIVATE KEY",
    Headers: map[string]string {
      "Subject": uSubject,
      "Date":    uDate,
    },
    Bytes: x509.MarshalPKCS1PrivateKey(uPrivateKey),
  })
  uPrivateKeyFileName := uDir+"/"+theUser+".key"
  err = ioutil.WriteFile(uPrivateKeyFileName, uPrivateKeyPEM.Bytes(), 0644)
  configMayBeFatal("could not write the ["+uPrivateKeyFileName+"] file", err)
}
