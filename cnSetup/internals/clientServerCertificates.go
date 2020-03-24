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
  "math/big"
  "net"
  "os"
  "path/filepath"
  "time"
)

// Create a Server's' x509 Certificate and associated public/private RSA 
// keys. The Server certificate and keys are written to disk in the PEM 
// format. 
//
// This code has been inspired by: Shane Utt's excellent article:
//   https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
//
func (nursery *NurseryType) CreateNurseryCertificate(
  nurseryNum int,
  ca        *CAType,
  config    *ConfigType,
) error {

  os.MkdirAll(nursery.Cert_Dir, 0755)
  caCertFile, caCertErr := os.Open(nursery.Ca_Cert_Path)
  certFile,   certErr   := os.Open(nursery.Cert_Path)
  keyFile,    keyErr    := os.Open(nursery.Key_Path)

  if (caCertErr == nil && certErr == nil && keyErr == nil) {
    fmt.Printf("\n\nCertificate files for the [%s] Nursery already exist\n", nursery.Name)
    fmt.Print( "  not recreating them.\n")
    caCertFile.Close()
    certFile.Close()
    keyFile.Close()
    return nil
  }

  fmt.Printf("\n\nCreating certificate files for the [%s] Nursery\n", nursery.Name)

  ca.StartReading()
  defer ca.StopReading()
  
  nCert := &x509.Certificate {
    // we need to use DIFFERENT serial numbers for each of CA (1<<32), 
    //  C/S  ((1<<5 + nurseryNum)<<33) and
    //  User ((2<<5 + userNum)<<33)
    SerialNumber: big.NewInt(
      (nursery.Serial_Number)<<33 |
      int64(ca.Serial_Number),
    ),    SignatureAlgorithm: x509.SHA512WithRSA,
    Subject: pkix.Name {
      Organization:  []string{ca.Organization},
      Country:       []string{ca.Country},
      Province:      []string{ca.Province},
      Locality:      []string{ca.Locality},
      StreetAddress: []string{ca.Street_Address},
      PostalCode:    []string{ca.Postal_Code},
      CommonName:    nursery.Name,
    },
    EmailAddresses:  []string{ca.Email_Address},
    NotBefore: time.Now(),
    NotAfter:  time.Now().AddDate(int(ca.Valid_For.Years),
                                  int(ca.Valid_For.Months),
                                  int(ca.Valid_For.Days)),
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
  ca.StopReading()

  // Add the DNSNames and IPAddresses
  for _, aHost := range nursery.Hosts {
    possibleIPAddress := net.ParseIP(aHost)
    if possibleIPAddress != nil {
      nCert.IPAddresses = append(nCert.IPAddresses, possibleIPAddress)
    } else {
      nCert.DNSNames = append(nCert.DNSNames, aHost)
    }
  }

  nPrivateKey, err := rsa.GenerateKey(rand.Reader, int(nursery.Key_Size))
  if err != nil {
    return fmt.Errorf("could not generate rsa key for ["+nursery.Name+"] Nursery: %w", err)
  }

  nBytes, err := x509.CreateCertificate(
    rand.Reader,
    nCert, ca.Cert,
    &nPrivateKey.PublicKey,
    ca.PrivateKey,
  )
  if err != nil {
    return fmt.Errorf("could not create the certificate for ["+nursery.Name+"] Nursery: %w", err)
  }

  nSubject := "Subject: ConTeXt Nursery " + config.Federation_Name + " Server Certificate for ["+nursery.Name+"] Nursery"
  nDate    := "Date:    "+time.Now().String()+"\n"

  caPEM := new(bytes.Buffer)
  caPEM.WriteString("\n")
  caPEM.WriteString(nSubject + " (CA)\n")
  caPEM.WriteString(nDate)
  pem.Encode(caPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: ca.Cert.Raw,
  })
  os.MkdirAll(filepath.Dir(nursery.Ca_Cert_Path), 0755)
  err = ioutil.WriteFile(nursery.Ca_Cert_Path, caPEM.Bytes(), 0644)
  if err != nil {
    return fmt.Errorf("could not write the ["+nursery.Ca_Cert_Path+"] file: %w", err)
  }

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
    Bytes: ca.Cert.Raw,
  })
  os.MkdirAll(filepath.Dir(nursery.Cert_Path), 0755)
  err = ioutil.WriteFile(nursery.Cert_Path, nPEM.Bytes(), 0644)
  if err != nil {
    return fmt.Errorf("could not write the ["+nursery.Cert_Path+"] file: %w", err)
  }

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
  if err != nil {
    return fmt.Errorf("could not write the ["+nursery.Key_Path+"] file: %w", err)
  }
  return nil
}
