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

package main

import (
  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/interfaces/discovery"
  "github.com/shirou/gopsutil/cpu"
  "github.com/shirou/gopsutil/load"
  "github.com/shirou/gopsutil/mem"
  "math/rand"
  "time"
)

///////////////////////
// Heart Beat interface

func sendPeriodicHeartBeats(cc *clientConnection.CC) {
  lConfig := getConfig()

  for {
    time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
    ni := discovery.NurseryInfo{
      Name: lConfig.Name,
      Port: lConfig.Port,
      Base_Url: lConfig.Base_Url,
      State: "up",
      Processes: 1,
    }

    loads, err := load.Avg()
    if err != nil {
      cnLog.MayBeError("Could not read the load average", err)
      loads = &load.AvgStat{ Load1: 1.0, Load5: 1.0, Load15: 1.0, }
    }
    ni.Load.Load1  = loads.Load1
    ni.Load.Load5  = loads.Load5
    ni.Load.Load15 = loads.Load15

    cpuInfo, err := cpu.Info()
    if err != nil {
      cnLog.MayBeError("Could not read the cpu information", err)
      cpuInfo = []cpu.InfoStat{ cpu.InfoStat{ Cores: 1, Mhz: 1000 } }
    }
    ni.Cores     = uint(len(cpuInfo))
    ni.Speed_Mhz = cpuInfo[0].Mhz

    virtMem, err := mem.VirtualMemory()
    if err != nil {
      cnLog.MayBeError("Could not read the virtual memory information", err)
      virtMem = &mem.VirtualMemoryStat{ Total: 1000, Used: 1000 }
    }
    ni.Memory.Total = virtMem.Total
    ni.Memory.Used  = virtMem.Used

    swapMem, err := mem.SwapMemory()
    if err != nil {
      cnLog.MayBeError("Could not read the swap memory information", err)
      swapMem = &mem.SwapMemoryStat{ Total: 1000, Used: 1000 }
    }
    ni.Swap.Total = swapMem.Total
    ni.Swap.Used  = swapMem.Used

    cnLog.Json("beat request ", "ni", ni)
    niInfoMap := discovery.SendDiscoveryMessage(lConfig.Primary_Url, ni, cc)
    cnLog.Json("beat response ", "niInfoMap", niInfoMap)
    cnInfoMap.ActionUpdateNurseryInfoMap(niInfoMap)
  }

}

