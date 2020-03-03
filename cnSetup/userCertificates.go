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
  "software.sslmate.com/src/go-pkcs12"
  "strings"
  "time"
)

////////////////////
// User Certificates

func createUserCertificate(theUser string, userNum int) {
  if theUser == "" {
    log.Printf("cnConfig(WARNING): no user name specified for a user, skipping user[%d]\n", userNum)
    return
  }
  fmt.Printf("\n\nCreating configuration for the user [%s]\n", theUser)

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
  setupMayBeFatal("could not generate rsa key for user ["+theUser+"]", err)

  uBytes, err := x509.CreateCertificate(rand.Reader, uCert, caCert, &uPrivateKey.PublicKey, caPrivateKey)
  setupMayBeFatal("could not create the certificate for user ["+theUser+"]", err)

  uSubject := "Subject: ConTeXt Nursery " + config.Federation_Name + " User Certificate for user ["+theUser+"]"
  uDate    := "Date:    "+time.Now().String()+"\n"

  uDir := "users/"+theUser
  os.MkdirAll(uDir, 0755)
  uPath := uDir+"/"+ strings.ReplaceAll(theUser, ".", "-")

  caPEM := new(bytes.Buffer)
  caPEM.WriteString("\n")
  caPEM.WriteString(uSubject + " (CA)\n")
  caPEM.WriteString(uDate)
  pem.Encode(caPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: caCert.Raw,
  })
  uCaCertificateFileName := uPath+"-ca-crt.pem"
  err = ioutil.WriteFile(uCaCertificateFileName, caPEM.Bytes(), 0644)
  setupMayBeFatal("could not write the ["+uCaCertificateFileName+"] file", err)

  uPEM := new(bytes.Buffer)
  uPEM.WriteString("\n")
  uPEM.WriteString(uSubject + "\n")
  uPEM.WriteString(uDate)
  pem.Encode(uPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: uBytes,
  })
  //
  // add the CA certificate to the chain..
  //
  uPEM.WriteString("\n")
  uPEM.WriteString(uSubject + " (CA)\n")
  uPEM.WriteString(uDate)
  pem.Encode(uPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: caCert.Raw,
  })
  uCertificateFileName := uPath+"-crt.pem"
  err = ioutil.WriteFile(uCertificateFileName, uPEM.Bytes(), 0644)
  setupMayBeFatal("could not write the ["+uCertificateFileName+"] file", err)

  uPrivateKeyPEM := new(bytes.Buffer)
  uPrivateKeyPEM.WriteString("\n")
  uPrivateKeyPEM.WriteString(uSubject + "\n")
  uPrivateKeyPEM.WriteString(uDate)
  pem.Encode(uPrivateKeyPEM, &pem.Block {
    Type: "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(uPrivateKey),
  })
  uPrivateKeyFileName := uPath+"-key.pem"
  err = ioutil.WriteFile(uPrivateKeyFileName, uPrivateKeyPEM.Bytes(), 0600)
  setupMayBeFatal("could not write the ["+uPrivateKeyFileName+"] file", err)


//  uCert, err := x509.ParseCertificate(uBytes)
  setupMayBeFatal("could not parse x509 certificate", err)

  pfxBytes, err := pkcs12.Encode(rand.Reader, uPrivateKey, uCert, []*x509.Certificate{caCert}, "test")
  setupMayBeFatal("Could not create the pkcs#12 certificate bundle", err)

  uPKCS12FileName := uPath+".p12"
  err = ioutil.WriteFile(uPKCS12FileName, pfxBytes, 0600)
  setupMayBeFatal("Could not write the pkcs#12 certifcate to a file", err)
}
