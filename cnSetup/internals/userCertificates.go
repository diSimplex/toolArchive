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
  "github.com/sethvargo/go-password/password"
  "io/ioutil"
  "math/big"
  "log"
  "os"
  "os/exec"
//  "software.sslmate.com/src/go-pkcs12"
  "strings"
  "sync"
  "time"
)

////////////////////
// User Certificates

// Create a user's X509 certificates and public/private keys.
//
// We provide the name of the user (usually one of their email addresses), 
// a user unique int (used to ensure the user's serial number is unique). 
//
// We also provide an optional WaitGroup which, if not nil, is used to 
// allow this function to be called asynchronously as a go routine. 
//
func (user *User) CreateUserCertificate(
  ca *CertificateAuthority,
  wg *sync.WaitGroup,
) {
  if wg != nil {
    wg.Add(1)
    defer wg.Done()
  }

  if user.Name == "" {
    log.Printf("cnConfig(WARNING): no user name specified for a user, skipping user[%d]\n", userNum)
    return
  }

  // TODO sort this out with user.NormalizeConfiguration
  uDir := "users/"+user.Name
  os.MkdirAll(uDir, 0755)
  uPath := uDir+"/"+ strings.ReplaceAll(user.Name, ".", "-")

  uCaCertificateFileName := uPath+"-ca-crt.pem"
  caCertFile, caCertErr := os.Open(uCaCertificateFileName)

  uCertificateFileName := uPath+"-crt.pem"
  certFile, certErr := os.Open(uCertificateFileName)

  uPrivateKeyFileName := uPath+"-key.pem"
  keyFile, keyErr := os.Open(uPrivateKeyFileName)

  uPKCS12FileName := uPath+"-pkcs12.p12"
  pkcsFile, pkcsErr := os.Open(uPKCS12FileName)

  if (caCertErr == nil && certErr == nil && keyErr == nil && pkcsErr == nil) {
    fmt.Printf("\n\nCertificate files for the user [%s] already exist\n", user.Name)
    fmt.Print( "  not recreating them.\n")
    caCertFile.Close()
    certFile.Close()
    keyFile.Close()
    pkcsFile.Close()
    return
  }

  fmt.Printf("\n\nCreating certificate files for the user [%s]\n", user.Name)

  ca.StartUsing()
  defer ca.StopUsing()
  
  uCert := &x509.Certificate {
    // we need to use DIFFERENT serial numbers for each of CA (1<<32),
    //  C/S  ((1<<5 + nurseryNum)<<33) and
    //  User ((2<<5 + userNum)<<33)
    SerialNumber: big.NewInt(
      (user.Serial_Number)<<33 |
      int64(ca.Serial_Number),
    ),
    SignatureAlgorithm: x509.SHA512WithRSA,
    Subject: pkix.Name {
      Organization:  []string{ca.Organization},
      Country:       []string{ca.Country},
      Province:      []string{ca.Province},
      Locality:      []string{ca.Locality},
      StreetAddress: []string{ca.Street_Address},
      PostalCode:    []string{ca.Postal_Code},
      CommonName:    user.Name + " ( ConTeXt Nursery " + config.Federation_Name + " )",
    },
    EmailAddresses:  []string{ca.Email_Address},
    NotBefore: time.Now(),
    NotAfter:  time.Now().AddDate(int(ca.Valid_For.Years),
                                  int(ca.Valid_For.Months),
                                  int(ca.Valid_For.Days)),
    ExtKeyUsage: []x509.ExtKeyUsage{
      x509.ExtKeyUsageClientAuth,
    },
    SubjectKeyId: []byte{1,2,3,4,6},
    KeyUsage:    x509.KeyUsageDigitalSignature |
      x509.KeyUsageKeyEncipherment |
      x509.KeyUsageKeyAgreement |
      x509.KeyUsageDataEncipherment,
  }
  ca.StopUsing()
  
  uPrivateKey, err := rsa.GenerateKey(rand.Reader, int(config.Key_Size))
  setupMayBeFatal("could not generate rsa key for user ["+user.Name+"]", err)

  uBytes, err := x509.CreateCertificate(rand.Reader, uCert, caCert, &uPrivateKey.PublicKey, caPrivateKey)
  setupMayBeFatal("could not create the certificate for user ["+user.Name+"]", err)

  uSubject := "Subject: ConTeXt Nursery " + config.Federation_Name + " User Certificate for user ["+user.Name+"]"
  uDate    := "Date:    "+time.Now().String()+"\n"

  caPEM := new(bytes.Buffer)
  caPEM.WriteString("\n")
  caPEM.WriteString(uSubject + " (CA)\n")
  caPEM.WriteString(uDate)
  pem.Encode(caPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: caCert.Raw,
  })
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
  err = ioutil.WriteFile(uPrivateKeyFileName, uPrivateKeyPEM.Bytes(), 0600)
  setupMayBeFatal("could not write the ["+uPrivateKeyFileName+"] file", err)

//  uCert, err := x509.ParseCertificate(uBytes)
//  setupMayBeFatal("could not parse x509 certificate", err)

//  pfxBytes, err := pkcs12.Encode(rand.Reader, uPrivateKey, uCert, []*x509.Certificate{caCert}, "test")
//  setupMayBeFatal("Could not create the pkcs#12 certificate bundle", err)

//  uPKCS12FileName := uPath+".p12"
//  err = ioutil.WriteFile(uPKCS12FileName, pfxBytes, 0600)
//  setupMayBeFatal("Could not write the pkcs#12 certifcate to a file", err)

//  openssl pkcs12 -export
//    -out stephen\@perceptisys-co-uk.p12
//    -inkey stephen\@perceptisys-co-uk-key.pem
//    -in stephen\@perceptisys-co-uk-crt.pem
//    -certfile stephen\@perceptisys-co-uk-ca-crt.pem

  thePassword, err := password.Generate(8, 2, 0, false, false)
  setupMayBeFatal("Could not generate a password", err)
  userPasswords[user.Name] = thePassword

  err = os.Setenv("OPENSSL_PASSWORD", thePassword)
  setupMayBeFatal("Could not set the OPENSSL_PASSWORD environment variable", err)

  cmd := exec.Command("openssl", "pkcs12", "-export",
    "-out", uPKCS12FileName,
    "-inkey", uPrivateKeyFileName,
    "-in", uCertificateFileName,
    "-certfile", uCaCertificateFileName,
    "-passout", "env:OPENSSL_PASSWORD",
  )
  outErr, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Printf("ERROR:\n--------------\n%s\n--------------\n", outErr)
    setupMayBeFatal("Could not create the pkcs#12 file", err)
  }
}
