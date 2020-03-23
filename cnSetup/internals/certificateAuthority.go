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
  "math/big"
  "os"
  "sync"
  "time"
)

////////////////////////
// Certificate Authority
//

type  CAType struct {
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
  
  Dir            string
  Cert_File_Name string
  Key_File_Name  string
  
  Cert          *x509.Certificate
  PrivateKey    *rsa.PrivateKey
  
  federationName string
  keySize        uint
  
  mutex          *sync.RWMutex
}

func (config *ConfigType) NormalizeCA() {
  if config.Federation_Name != "" {
     config.Certificate_Authority.federationName = config.Federation_Name
     
    if config.Certificate_Authority.Dir == "" {
      config.Certificate_Authority.Dir = "ca/" + config.Federation_Name
    }
    if config.Certificate_Authority.Cert_File_Name == "" {
      config.Certificate_Authority.Cert_File_Name = 
        config.Certificate_Authority.Dir + "/" +
        config.Federation_Name + "-ca-crt.pem"
    }
    if config.Certificate_Authority.Key_File_Name == "" {
     config.Certificate_Authority.Key_File_Name = 
       config.Certificate_Authority.Dir + "/" +
       config.Federation_Name + "-ca-key.pem"
    }
  } else {
    config.csLog.Logf("You MUST specify a Federation Name")
    os.Exit(-1)
  }
  
  if 1023 < config.Key_Size {
    config.Certificate_Authority.keySize = config.Key_Size
  } else {
    config.csLog.Logf("You MUST specify a Key_Size of at least 1024")
    os.Exit(-1)
  }
}

func CreateCA(config *ConfigType) *CAType {
  newCA := config.Certificate_Authority
  return &newCA
}

func (ca *CAType) StartUsing() {
  ca.mutex.Lock()
}

func (ca *CAType) StopUsing() {
  ca.mutex.Unlock()
}

func (ca *CAType) CreateNewCA() error {
  fmt.Print("\nCreating a new Certificate Authority for [%s]\n", ca.federationName)

  ca.Cert = &x509.Certificate {
    // we need to use DIFFERENT serial numbers for each of CA (1<<32), 
    //  C/S (1<<33) and User (1<<34)
    SerialNumber: big.NewInt(int64(1<<32) | int64(ca.Serial_Number)),
    SignatureAlgorithm: x509.SHA512WithRSA,
    Subject: pkix.Name {
      Organization:  []string{ca.Organization},
      Country:       []string{ca.Country},
      Province:      []string{ca.Province},
      Locality:      []string{ca.Locality},
      StreetAddress: []string{ca.Street_Address},
      PostalCode:    []string{ca.Postal_Code},
      CommonName:    "ConTeXt Nursery "+ca.Common_Name,
    },
    EmailAddresses:  []string{ca.Email_Address},
    NotBefore: time.Now(),
    NotAfter:  time.Now().AddDate(int(ca.Valid_For.Years),
                                  int(ca.Valid_For.Months),
                                  int(ca.Valid_For.Days)),
    IsCA:        true,
    ExtKeyUsage: []x509.ExtKeyUsage{
      x509.ExtKeyUsageClientAuth,
      x509.ExtKeyUsageServerAuth,
    },
    KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign ,
    BasicConstraintsValid: true,
  }

  var err error
  ca.PrivateKey, err = rsa.GenerateKey(rand.Reader, int(ca.keySize))
  if err != nil {
    return fmt.Errorf("could not generate rsa key for CA: %w", err)
  }

  ca.Cert.Raw, err = x509.CreateCertificate(
    rand.Reader,
    ca.Cert, ca.Cert,
    &ca.PrivateKey.PublicKey,
    ca.PrivateKey,
  )
  if err != nil {
    return fmt.Errorf("could not create the CA certificate: %w", err)
  }

  return nil
}

func (ca *CAType) WriteCAFiles() error {
  caSubject := "Subject: ConTeXt Nursery " + ca.federationName + " Certificate Authority\n"
  caDate    := "Date:    "+time.Now().String()+"\n"

  os.MkdirAll(ca.Dir, 0755)

  caPEM := new(bytes.Buffer)
  caPEM.WriteString("\n")
  caPEM.WriteString(caSubject)
  caPEM.WriteString(caDate)
  pem.Encode(caPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: ca.Cert.Raw,
  })
  err := ioutil.WriteFile(ca.Cert_File_Name, caPEM.Bytes(), 0644)
  if err != nil {
    return fmt.Errorf("could not write the certificateAuthority.crt file: %w", err)
  }
  
  // NOTE this private key is left UN-ENCRYPTED on the file system!
  // SO you need to ensure it is not readable by anyone other than the
  // user who needs to run the cnSetup!
  //
  caPrivateKeyPEM := new(bytes.Buffer)
  caPrivateKeyPEM.WriteString("\n")
  caPrivateKeyPEM.WriteString(caSubject)
  caPrivateKeyPEM.WriteString(caDate)
  pem.Encode(caPrivateKeyPEM, &pem.Block {
    Type: "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(ca.PrivateKey),
  })
  err = ioutil.WriteFile(ca.Key_File_Name, caPrivateKeyPEM.Bytes(), 0600)
  if err != nil {
    return fmt.Errorf("could not write the certificateAuthority.key file: %w", err)
  }
  
  return nil
}

func (ca *CAType) LoadCAFromFiles() error {
  caCertBytes, err := ioutil.ReadFile(ca.Cert_File_Name)
  if err != nil {
    return fmt.Errorf("could not load the certificate authority's *.crt file: %w", err)
  }

  caCertPEM, _ /*restCaCertBytes*/ := pem.Decode(caCertBytes)
  if caCertPEM == nil || caCertPEM.Type != "CERTIFICATE" {
    return fmt.Errorf("could not locate the certificate authority's CERTIFICATE block: %w", err)
  }

  lcaCert, err := x509.ParseCertificate(caCertPEM.Bytes)
  if err != nil {
    return fmt.Errorf("could not parse the certificate authority's certificate: %w", err)
  }

  caKeyBytes,  err  := ioutil.ReadFile(ca.Key_File_Name)
  if err != nil {
    return fmt.Errorf("could not load the certificate authority's *.key file: %w", err)
  }

  caKeyPEM, _ /*restCaKeyBytes*/ := pem.Decode(caKeyBytes)
  if caKeyPEM == nil || caKeyPEM.Type != "RSA PRIVATE KEY" {
    return fmt.Errorf("could not locate the certificate authority's RSA PRIVATE KEY block: %w", err)
  }

  lcaPrivateKey, err := x509.ParsePKCS1PrivateKey(caKeyPEM.Bytes)
  if err != nil {
    return fmt.Errorf("could not parse the certificate authority's private key: 5w", err)
  }

  // If we managed to get this far... both the cert and key are OK...
  // so store the local copies in the global variables...
  ca.Cert       = lcaCert
  ca.PrivateKey = lcaPrivateKey
  
  return nil
}
