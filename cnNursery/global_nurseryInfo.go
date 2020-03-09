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

// This file collects all of the globals required for the cnNursery process.
//
// Since cnNursery makes essential use of multi-threading, we need to 
// ensure all globals are thread safe. To do this we make liberal use
// of the sync.RWMutexes, one for each global.
//

package main

import (
  "encoding/json"
  "sync"
)

//////////////////////
// Nursery Information

type MemoryTU struct {
  Total uint64
  Used  uint64
}

type NurseryInfo struct {
  Name      string
  Port      uint
  State     string
  Processes uint
  Cores     uint
  Speed_Mhz float64
  Memory    MemoryTU
  Swap      MemoryTU
  Load      struct {
    Load1   float64
    Load5   float64
    Load15  float64
  }
}

type NurseryInfoMap map[string]NurseryInfo

// The heartBeatInfo is the federation information as managed by the heart
// beat *collector* (POST at "/heartbeat") at the primary cnNursery

var heartBeatInfoSync sync.RWMutex
var heartBeatInfoPriv NurseryInfoMap

func initHeartBeatInfo() {
 heartBeatInfoSync.Lock()
 defer heartBeatInfoSync.Unlock()

 heartBeatInfoPriv = make(map[string]NurseryInfo)
}

func getHeartBeatInfoMap() NurseryInfoMap {
  heartBeatInfoSync.RLock()
  defer heartBeatInfoSync.RUnlock()

  return heartBeatInfoPriv
}

func setHeartBeatInfoNurseryFromJsonBytes(jsonBytes[]byte) error {
  heartBeatInfoSync.Lock()
  defer heartBeatInfoSync.Unlock()

  var ni NurseryInfo
  err := json.Unmarshal(jsonBytes, &ni)
  if err != nil { return err }

// debug
  cnNurseryJson("heartBeat ni:", "NurseryInfo", ni)

  if ni.Name != "" {
    heartBeatInfoPriv[ni.Name] = ni
  }

  return nil
}

func jsonMarshalHeartBeatInfo() ([]byte, error) {
  heartBeatInfoSync.RLock()
  defer heartBeatInfoSync.RUnlock()

  return json.Marshal(heartBeatInfoPriv)
}

// The federationInfo is the federation information as managed by the heart
// beat sender in each "seconday" cnNursery (including the primary
// cnNursery's heart beat).

var federationInfoSync sync.RWMutex
var federationInfoPriv NurseryInfoMap

func getFederationInfoMap() NurseryInfoMap {
  federationInfoSync.RLock()
  defer federationInfoSync.RUnlock()

  return federationInfoPriv
}

func jsonUnmarshalFederationInfo(jsonBytes []byte) {
  federationInfoSync.Lock()
  defer federationInfoSync.Unlock()

  json.Unmarshal(jsonBytes, &federationInfoPriv)
}
