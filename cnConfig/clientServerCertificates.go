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
  "fmt"
  "log"
  "strings"
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
}
