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
  "bytes"
  "crypto/rand"
  "crypto/rsa"
  "crypto/x509"
  "crypto/x509/pkix"
  "encoding/pem"
  "fmt"
  "io/ioutil"
  "log"
  "math/big"
  "net"
  "os"
  "path/filepath"
//  "strings"
  "time"
)

/////////////////////////////
// Client Server Certificates

func createNurseryCertificate(nursery *Nursery, nurseryNum int) {
  if nursery.Host == "" {
    log.Printf("cnConfig(WARNING): no host names specified for a Nursery, skipping Nursery[%d]\n", nurseryNum)
    return
  }

  nDir := "servers/"+nursery.Name
  os.MkdirAll(nDir, 0755)
  if nursery.Ca_Cert_Path == "" {
    nursery.Ca_Cert_Path = nDir+"/"+nursery.Name+"-ca-crt.pem"
  }
  caCertFile, caCertErr := os.Open(nursery.Ca_Cert_Path)
  if nursery.Cert_Path == "" {
    nursery.Cert_Path = nDir+"/"+nursery.Name+"-crt.pem"
  }
  certFile, certErr := os.Open(nursery.Ca_Cert_Path)
  if nursery.Key_Path == "" {
    nursery.Key_Path = nDir+"/"+nursery.Name+"-key.pem"
  }
  keyFile, keyErr := os.Open(nursery.Ca_Cert_Path)

  if (caCertErr == nil && certErr == nil && keyErr == nil) {
    fmt.Printf("\n\nCertificate files for the [%s] Nursery already exist\n", nursery.Name)
    fmt.Print( "  not recreating them.\n")
    caCertFile.Close()
    certFile.Close()
    keyFile.Close()
    return
  }

  fmt.Printf("\n\nCreating certificate files for the [%s] Nursery\n", nursery.Name)

  nCert := &x509.Certificate {
    // we need to use DIFFERENT serial numbers for each of CA (1<<32), 
    //  C/S  ((1<<5 + nurseryNum)<<33) and
    //  User ((2<<5 + userNum)<<33)
    SerialNumber: big.NewInt(
      int64(1<<5 + nurseryNum)<<33 |
      int64(config.Certificate_Authority.Serial_Number),
    ),    SignatureAlgorithm: x509.SHA512WithRSA,
    Subject: pkix.Name {
      Organization:  []string{config.Certificate_Authority.Organization},
      Country:       []string{config.Certificate_Authority.Country},
      Province:      []string{config.Certificate_Authority.Province},
      Locality:      []string{config.Certificate_Authority.Locality},
      StreetAddress: []string{config.Certificate_Authority.Street_Address},
      PostalCode:    []string{config.Certificate_Authority.Postal_Code},
      CommonName:    nursery.Name,
    },
    EmailAddresses:  []string{config.Certificate_Authority.Email_Address},
    NotBefore: time.Now(),
    NotAfter:  time.Now().AddDate(int(config.Certificate_Authority.Valid_For.Years),
                                  int(config.Certificate_Authority.Valid_For.Months),
                                  int(config.Certificate_Authority.Valid_For.Days)),
    ExtKeyUsage: []x509.ExtKeyUsage{
      x509.ExtKeyUsageClientAuth,
      x509.ExtKeyUsageServerAuth,
    },
    SubjectKeyId: []byte{1,2,3,4,6},
    KeyUsage:    x509.KeyUsageDigitalSignature |
      x509.KeyUsageKeyEncipherment |
      x509.KeyUsageKeyAgreement |
      x509.KeyUsageDataEncipherment,
  }

  // Add the DNSNames and IPAddresses
  for _, aHost := range nursery.Hosts {
    possibleIPAddress := net.ParseIP(aHost)
    if possibleIPAddress != nil {
      nCert.IPAddresses = append(nCert.IPAddresses, possibleIPAddress)
    } else {
      nCert.DNSNames = append(nCert.DNSNames, aHost)
    }
  }

  nPrivateKey, err := rsa.GenerateKey(rand.Reader, int(config.Key_Size))
  setupMayBeFatal("could not generate rsa key for ["+nursery.Name+"] Nursery", err)

  nBytes, err := x509.CreateCertificate(rand.Reader, nCert, caCert, &nPrivateKey.PublicKey, caPrivateKey)
  setupMayBeFatal("could not create the certificate for ["+nursery.Name+"] Nursery", err)

  nSubject := "Subject: ConTeXt Nursery " + config.Federation_Name + " Server Certificate for ["+nursery.Name+"] Nursery"
  nDate    := "Date:    "+time.Now().String()+"\n"

  caPEM := new(bytes.Buffer)
  caPEM.WriteString("\n")
  caPEM.WriteString(nSubject + " (CA)\n")
  caPEM.WriteString(nDate)
  pem.Encode(caPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: caCert.Raw,
  })
  os.MkdirAll(filepath.Dir(nursery.Ca_Cert_Path), 0755)
  err = ioutil.WriteFile(nursery.Ca_Cert_Path, caPEM.Bytes(), 0644)
  setupMayBeFatal("could not write the ["+nursery.Ca_Cert_Path+"] file", err)

  nPEM := new(bytes.Buffer)
  nPEM.WriteString("\n")
  nPEM.WriteString(nSubject + "\n")
  nPEM.WriteString(nDate)
  pem.Encode(nPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: nBytes,
  })
  //
  // add the CA certificate to the chain..
  //
  nPEM.WriteString("\n")
  nPEM.WriteString(nSubject + " (CA)\n")
  nPEM.WriteString(nDate)
  pem.Encode(nPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: caCert.Raw,
  })
  os.MkdirAll(filepath.Dir(nursery.Cert_Path), 0755)
  err = ioutil.WriteFile(nursery.Cert_Path, nPEM.Bytes(), 0644)
  setupMayBeFatal("could not write the ["+nursery.Cert_Path+"] file", err)

  // NOTE this private key is left UN-ENCRYPTED on the file system!
  // SO you need to ensure it is not readable by anyone other than the
  // user who needs to run the cnNursery!
  //
  nPrivateKeyPEM := new(bytes.Buffer)
  nPrivateKeyPEM.WriteString("\n")
  nPrivateKeyPEM.WriteString(nSubject + "\n")
  nPrivateKeyPEM.WriteString(nDate)
  pem.Encode(nPrivateKeyPEM, &pem.Block {
    Type: "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(nPrivateKey),
  })
  os.MkdirAll(filepath.Dir(nursery.Key_Path), 0755)
  err = ioutil.WriteFile(nursery.Key_Path, nPrivateKeyPEM.Bytes(), 0644)
  setupMayBeFatal("could not write the ["+nursery.Key_Path+"] file", err)
}
