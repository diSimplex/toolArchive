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

// A clientConnection which sends RESTful HTTP messages to the 
// webserver.
//
package clientConnection

import (
  "bytes"
  "crypto/tls"
  "crypto/x509"
  "github.com/diSimplex/ConTeXtNursery/logger"
  "io/ioutil"
  "net/http"
  "time"
)

type CC struct {
  Client     *http.Client
  Log        *logger.LoggerType
}

func CreateClientConnection(
  caCertPath, certPath, keyPath string,
  cnLog      *logger.LoggerType,
) *CC {

  // Load the Server x509 Certificates and keys for this client connection
  //
  serverCert, err := tls.LoadX509KeyPair( certPath, keyPath )
  cnLog.MayBeFatal("Could not load cert/key pair", err)
  //
  caCert, err := ioutil.ReadFile(caCertPath)
  cnLog.MayBeFatal("Could not load the CA certificate", err)
  //
  caCertPool := x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caCert)
  //
  // Setup HTTPS server configuration
  //
  tlsConfig := &tls.Config{
    ClientAuth:     tls.RequireAndVerifyClientCert,
    Certificates: []tls.Certificate{serverCert},
    RootCAs:        caCertPool,
    ClientCAs:      caCertPool,
  }

  cc := CC{}

  transport := &http.Transport{
    TLSClientConfig:    tlsConfig,
    ForceAttemptHTTP2:  true,
    MaxIdleConns:       10,
    IdleConnTimeout:    30 * time.Second,
    DisableCompression: true,
  }

  cc.Client = &http.Client{
    Transport: transport,
  }

  cc.Log = cnLog

  return &cc
}


func (cc CC) GetMessage(baseUrl, url string) []byte {

  replyBytes := []byte{}

  resp, err := cc.Client.Get(baseUrl + url)
  if err != nil {
    cc.Log.MayBeError(
      "Could not get client connection request to the Nursery: "+baseUrl,
      err,
    )
    if resp != nil {
      cc.Log.Logf("Response code: %s / %s", resp.Status, resp.Proto)
    }
    return replyBytes
  }
  defer resp.Body.Close()

  respBody, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    cc.Log.MayBeFatal(
      "Could not read the body of the client connection response",
       err,
    )
    return replyBytes
  }

  return respBody
}

func (cc CC) SendJsonMessage(baseUrl, url, method string, jsonBytes []byte) []byte {

  replyBytes := []byte{}

  ccReq, err := http.NewRequest(
    method,
    baseUrl + url,
    bytes.NewReader(jsonBytes),
  )
  if err != nil {
    cc.Log.MayBeError("Could not create client connection request", err)
    return replyBytes
  }

  ccReq.Header.Add("Accept", "application/json")

  resp, err := cc.Client.Do(ccReq)
  if err != nil {
    cc.Log.MayBeError(
      "Could not send client connection request to the Nursery: "+baseUrl,
      err,
    )
    if resp != nil {
      cc.Log.Logf("Response code: %s / %s", resp.Status, resp.Proto)
    }
    return replyBytes
  }
  defer resp.Body.Close()

  respBody, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    cc.Log.MayBeFatal(
      "Could not read the body of the client connection response",
       err,
    )
    return replyBytes
  }

  return respBody
}

