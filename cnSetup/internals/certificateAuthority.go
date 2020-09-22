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
  "github.com/diSimplex/ConTeXtNursery/logger"
  "io/ioutil"
  "math/big"
  "os"
  "time"
)

// CAType contains a Certificate Authority's x509 certificate as well as 
// public/private RSA keys, as well as auxilary fields to control where 
// external PEM files can be found or stored. 
//
// CONSTRAINTS: Once created, the values in this structure SHOULD only be 
// altered by structure methods. 
//
type  CAType struct {
  // Standard x509 fields
  //
  Serial_Number   uint
  Organization    string
  Federation_Name string
  Country         string
  Province        string
  Locality        string
  Street_Address  string
  Postal_Code     string
  Email_Address   string
  Common_Name     string
  //
  Valid_For struct {
    Years  uint `default:"10"`
    Months uint `default:"0"`
    Days   uint `default:"0"`
  }
 
  // Auxilary fields required to write CA files
  //
  Dir            string
  Cert_File_Name string
  Key_File_Name  string
  
  // Auxilary fields required to create, contain and manage access 
  // to the actual certificates and keys. 
  //
  Key_Size       uint
  Cert          *x509.Certificate
  PrivateKey    *rsa.PrivateKey
  CSLog         *logger.LoggerType
}

// Normalize the Certificate Authority auxilary fields.
//
// ALTERS ca;
// NOT THREAD-SAFE;
// CALLED BY: LoadConfiguration ONLY;
//
func (ca *CAType) NormalizeCA(config *ConfigType) {
  // Make sure the Serial_Number is constantly increasing...
  //
  if ca.Serial_Number == 0 { ca.Serial_Number = uint(time.Now().Unix()) }

  if config.Federation_Name != "" {
    ca.Federation_Name = config.Federation_Name
    
    if ca.Dir == "" {
      ca.Dir = "ca/" + config.Federation_Name
    }
    if ca.Cert_File_Name == "" {
      ca.Cert_File_Name = 
        ca.Dir + "/" +
        config.Federation_Name + "-ca-crt.pem"
    }
    if ca.Key_File_Name == "" {
     ca.Key_File_Name = 
       ca.Dir + "/" +
       config.Federation_Name + "-ca-key.pem"
    }
  } else {
    config.CSLog.Logf("You MUST specify a Federation Name")
    os.Exit(-1)
  }
  
  if 1023 < config.Key_Size {
    if ca.Key_Size == 0 { ca.Key_Size = config.Key_Size }
  } else {
    config.CSLog.Logf("You MUST specify a Key_Size of at least 1024")
    os.Exit(-1)
  }
}

// Create the Certificate Authority Structure (only) from the details in 
// the configuraiton. 
//
// CREATES ca;
//
func CreateCA(config *ConfigType) *CAType {
  newCA       := config.Certificate_Authority
  newCA.CSLog  = config.CSLog
  return &newCA
}

// Cerate a new "base" x509 Certificate based upon the CA's configured 
// certificate information.
//
// Various fields specific to a particular certificate use will still need 
// to be filed in by the CA, Nursery, or User certificate code 
// respectively.
//
// It is CRITICAL that we use DIFFERENT serial numbers for each of the: 
//  - Certificate Authority:  1,
//  - Clien/Server:           (1<<5) + nurseryNum, and
//  - User:                   (2<<5) + userNum
// certificates. We do this using the "serialNumModifier" parameter. (This 
// assumes a maximum of 2^5 - 1 = 31 nurseries or 2^6 - 1 = 63 users) 
//
// READS ca;
//
func (ca *CAType) NewBaseCertificate(
  commonName string,
  serialNumModifier uint,
) *x509.Certificate {  
  return &x509.Certificate {
    SerialNumber: big.NewInt(int64(serialNumModifier)<<32 | int64(ca.Serial_Number)),
    SignatureAlgorithm: x509.SHA512WithRSA,
    Subject: pkix.Name {
      Organization:       []string{ca.Organization},
      OrganizationalUnit: []string{ca.Federation_Name},
      Country:            []string{ca.Country},
      Province:           []string{ca.Province},
      Locality:           []string{ca.Locality},
      StreetAddress:      []string{ca.Street_Address},
      PostalCode:         []string{ca.Postal_Code},
      CommonName:         commonName,
    },
    EmailAddresses:       []string{ca.Email_Address},
    NotBefore: time.Now(),
    NotAfter:  time.Now().AddDate(int(ca.Valid_For.Years),
                                  int(ca.Valid_For.Months),
                                  int(ca.Valid_For.Days)),
  }
}

// Create a new RSA Public/Private Key pair.
//
// IGNORES ca;
//
func (ca *CAType) NewRsaKeys(keySize uint) (*rsa.PrivateKey, error) {
  return rsa.GenerateKey(rand.Reader, int(keySize))
}

// Creates a new signed x509 Certificate returned as as "raw" ([]byte) 
// certificate using an x509 Certificate and its associated RSA public 
// key. 
//
// READS ca;
//
func (ca *CAType) SignCertificate(
  certToSign     *x509.Certificate,
  certPublicKey  *rsa.PublicKey,
) ([]byte, error) {
//  key, err := certPrivateKey.(crypto.Signer)
//  if err != nil {
//    return nil, fmt.Errorf("Could not type cast private key to crypto.Signer: %w", err)
//  }
  return x509.CreateCertificate(
    rand.Reader,
    certToSign, ca.Cert,
    certPublicKey,
    ca.PrivateKey,
  )
}

// Create a new Certificate Authority by creating a totally new 
// self-signed x509 certificate and associated public/private RSA keys. 
//
// This code has been inspired by: Shane Utt's excellent article:
//   https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
//
// ALTERS ca;
// NOT THREAD-SAFE;
//
func (ca *CAType) CreateNewCA() error {
  fmt.Printf("\nCreating a new Certificate Authority for [%s]\n", ca.Federation_Name)
  
  ca.Cert = ca.NewBaseCertificate(
    "ConTeXt Nursery "+ca.Common_Name,
    1,
  )
  //
  // Apply Certificate Authority only modifications
  //
  ca.Cert.IsCA = true
//  ca.Cert.ExtKeyUsage = []x509.ExtKeyUsage{
//      x509.ExtKeyUsageClientAuth,
//      x509.ExtKeyUsageServerAuth,
//    }
  ca.Cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign 
  ca.Cert.BasicConstraintsValid =true
  
  // Create a new RSA public/private key pair
  //
  var err error
  ca.PrivateKey, err = ca.NewRsaKeys(ca.Key_Size)
  if err != nil {
    return fmt.Errorf("could not generate rsa key for CA: %w", err)
  }
  
  // create a self-signed certificate using our own ca.Cert and 
  // ca.PrivateKey 
  //
  ca.Cert.Raw, err = ca.SignCertificate(ca.Cert, &ca.PrivateKey.PublicKey)
  if err != nil {
    return fmt.Errorf("could not create the CA certificate: %w", err)
  }

  return nil
}

// Write the Certificate Authority's x509 certificate and RSA keys to files
// on the disk. 
//
// READS ca;
//
func (ca *CAType) WriteCAFiles(config *ConfigType) error {
  caSubject := "Subject: ConTeXt Nursery " + config.Federation_Name + " Certificate Authority\n"
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

// Attempt to load an existing Certificate Authority from PEM files 
// containing x509 certificates and public/private RSA keys. 
//
// ALTERS ca;
// NOT THREAD-SAFE;
//
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
