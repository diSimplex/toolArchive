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
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "github.com/sethvargo/go-password/password"
  "io/ioutil"
  "os"
  "os/exec"
  "time"
)

////////////////////
// User Certificates

// Create a user's X509 certificates and public/private keys.
//
// We provide the name of the user (usually one of their email addresses), 
// a user unique int (used to ensure the user's serial number is unique). 
//
// This code has been inspired by: Shane Utt's excellent article:
//   https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
//
// READS ca;
// READS user;
//
func (user *UserType) CreateUserCertificate(
  userNum        int,
  ca            *CAType,
  federationName string,
) error {

  os.MkdirAll(user.Cert_Dir, 0755)
  caCertFile, caCertErr := os.Open(user.Ca_Cert_Path)
  certFile, certErr := os.Open(user.Cert_Path)
  keyFile, keyErr := os.Open(user.Key_Path)
  pkcsFile, pkcsErr := os.Open(user.Pkcs12_Path)

  if (caCertErr == nil && certErr == nil && keyErr == nil && pkcsErr == nil) {
    fmt.Printf("\n\nCertificate files for the user [%s] already exist\n", user.Name)
    fmt.Print( "  not recreating them.\n")
    caCertFile.Close()
    certFile.Close()
    keyFile.Close()
    pkcsFile.Close()
    return nil
  }

  fmt.Printf("\n\nCreating certificate files for the user [%s]\n", user.Name)
  
  uCert := ca.NewBaseCertificate(
    user.Name + " ( ConTeXt Nursery " + federationName + " )",
    user.Serial_Number,
  )
  uCert.ExtKeyUsage  = []x509.ExtKeyUsage{ x509.ExtKeyUsageClientAuth }
  uCert.SubjectKeyId = []byte{1,2,3,4,6}
  uCert.KeyUsage     = x509.KeyUsageDigitalSignature |
      x509.KeyUsageKeyEncipherment |
      x509.KeyUsageKeyAgreement |
      x509.KeyUsageDataEncipherment
  
  uPrivateKey, err := ca.NewRsaKeys(user.Key_Size)
  if err != nil {
    return fmt.Errorf("could not generate rsa key for user [%s]: %w", user.Name, err)
  }

  uBytes, err := ca.SignCertificate(uCert, &uPrivateKey.PublicKey)
  if err != nil {
    return fmt.Errorf("could not create the certificate for user [%s]: %w", user.Name, err)
  }

  uSubject := "Subject: ConTeXt Nursery " + federationName + " User Certificate for user ["+user.Name+"]"
  uDate    := "Date:    "+time.Now().String()+"\n"

  caPEM := new(bytes.Buffer)
  caPEM.WriteString("\n")
  caPEM.WriteString(uSubject + " (CA)\n")
  caPEM.WriteString(uDate)
  pem.Encode(caPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: ca.Cert.Raw,
  })
  err = ioutil.WriteFile(user.Ca_Cert_Path, caPEM.Bytes(), 0644)
  if err != nil {
    return fmt.Errorf("could not write the [%s] file: %w",  user.Ca_Cert_Path, err)
  }

  uPEM := new(bytes.Buffer)
  uPEM.WriteString("\n")
  uPEM.WriteString(uSubject + "\n")
  uPEM.WriteString(uDate)
  pem.Encode(uPEM, &pem.Block {
    Type:  "CERTIFICATE",
    Bytes: uBytes,
  })

  err = ioutil.WriteFile(user.Cert_Path, uPEM.Bytes(), 0644)
  if err != nil {
    return fmt.Errorf("could not write the [%s] file: %w", user.Cert_Path, err)
  }

  uPrivateKeyPEM := new(bytes.Buffer)
  uPrivateKeyPEM.WriteString("\n")
  uPrivateKeyPEM.WriteString(uSubject + "\n")
  uPrivateKeyPEM.WriteString(uDate)
  pem.Encode(uPrivateKeyPEM, &pem.Block {
    Type: "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(uPrivateKey),
  })
  err = ioutil.WriteFile(user.Key_Path, uPrivateKeyPEM.Bytes(), 0600)
  if err != nil {
    return fmt.Errorf("could not write the [%s] file: %w", user.Key_Path, err)
  }

//  uCert, err := x509.ParseCertificate(uBytes)
//  if err != nil {
//    return fmt.Errorf("could not parse x509 certificate: %w", err)
//  }

//  pfxBytes, err := pkcs12.Encode(rand.Reader, uPrivateKey, uCert, []*x509.Certificate{caCert}, "test")
//  if err != nil {
//    return fmt.Errorf("Could not create the pkcs#12 certificate bundle: %w", err)
//  }

//  err = ioutil.WriteFile(user.Pkcs12_Paht, pfxBytes, 0600)
//  if err != nil {
//    return fmt.Errorf("Could not write the pkcs#12 certifcate to a file: %w", err)
//  }

//  openssl pkcs12 -export
//    -out stephen\@perceptisys-co-uk.p12
//    -inkey stephen\@perceptisys-co-uk-key.pem
//    -in stephen\@perceptisys-co-uk-crt.pem
//    -certfile stephen\@perceptisys-co-uk-ca-crt.pem

  thePassword, err := password.Generate(8, 2, 0, false, false)
  if err != nil {
    return fmt.Errorf("Could not generate a password: %w", err)
  }
  user.Password = thePassword

  err = os.Setenv("OPENSSL_PASSWORD", thePassword)
  if err != nil {
    return fmt.Errorf("Could not set the OPENSSL_PASSWORD environment variable: %w", err)
  }

  cmd := exec.Command("openssl", "pkcs12", "-export",
    "-out", user.Pkcs12_Path,
    "-inkey", user.Key_Path,
    "-in", user.Cert_Path,
    "-certfile", user.Ca_Cert_Path,
    "-passout", "env:OPENSSL_PASSWORD",
  )
  outErr, err := cmd.CombinedOutput()
  if err != nil {
    fmt.Printf("ERROR:\n--------------\n%s\n--------------\n", outErr)
    return fmt.Errorf("Could not create the pkcs#12 file: %w", err)
  }
  return nil
}
