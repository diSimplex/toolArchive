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

package CNNurseries

import (
  "github.com/diSimplex/ConTeXtNursery/clientConnection"
  "github.com/diSimplex/ConTeXtNursery/interfaces/discovery"
  "github.com/shirou/gopsutil/cpu"
  "github.com/shirou/gopsutil/load"
  "github.com/shirou/gopsutil/mem"
  "math/rand"
  "time"
)

// Implements the heart beat go routine.
//
// Periodically attempt to contact each cnNursery in the cnInfoMap. If a 
// given cnNursery does not respond, then it is deleted from the 
// cnInfoMap. 
//
// READS config;
// CALLS cnState;
// CALLS cnInfoMap;
// PARAMETER cc (discovery.SendDiscoveryMessage);
//
func SendPeriodicHeartBeats(
  config    *ConfigType,
  cnState   *CNState,
  cnInfoMap *CNInfoMap,
  cc        *clientConnection.CC,
) {
  for {
    time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)
    config.CNLog.Logf("\n\n\nheartBeat state: [%s]\n\n\n", cnState.GetState())
    ni := discovery.NurseryInfo{
      Name: config.Name,
      Port: config.Port,
      Base_Url: config.Base_Url,
      State: cnState.GetState(),
      Processes: 1,
    }

    loads, err := load.Avg()
    if err != nil {
      config.CNLog.MayBeError("Could not read the load average", err)
      loads = &load.AvgStat{ Load1: 1.0, Load5: 1.0, Load15: 1.0, }
    }
    ni.Load.Load1  = loads.Load1
    ni.Load.Load5  = loads.Load5
    ni.Load.Load15 = loads.Load15

    cpuInfo, err := cpu.Info()
    if err != nil {
      config.CNLog.MayBeError("Could not read the cpu information", err)
      cpuInfo = []cpu.InfoStat{ cpu.InfoStat{ Cores: 1, Mhz: 1000 } }
    }
    ni.Cores     = uint(len(cpuInfo))
    ni.Speed_Mhz = cpuInfo[0].Mhz

    virtMem, err := mem.VirtualMemory()
    if err != nil {
      config.CNLog.MayBeError("Could not read the virtual memory information", err)
      virtMem = &mem.VirtualMemoryStat{ Total: 1000, Used: 1000 }
    }
    ni.Memory.Total = virtMem.Total
    ni.Memory.Used  = virtMem.Used

    swapMem, err := mem.SwapMemory()
    if err != nil {
      config.CNLog.MayBeError("Could not read the swap memory information", err)
      swapMem = &mem.SwapMemoryStat{ Total: 1000, Used: 1000 }
    }
    ni.Swap.Total = swapMem.Total
    ni.Swap.Used  = swapMem.Used

    config.CNLog.Json("beat request ", "ni", ni)
    niInfoMap := discovery.SendDiscoveryMessage(config.Primary_Url, ni, cc)
    config.CNLog.Json("beat response ", "niInfoMap", niInfoMap)
    cnInfoMap.ActionUpdateNurseryInfoMap(niInfoMap)
  }

}

