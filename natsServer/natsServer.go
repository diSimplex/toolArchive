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

// The connection to the cnNursery nats server.
//
package webserver

import (
  "fmt"
  "strings"
  "sync"

  "github.com/diSimplex/ConTeXtNursery/logger"
  "github.com/nats-io/nats.go"
)

type NATS struct {
  Conn     *nats.Conn
  Subs   []*nats.Subscription
  NatsWG    sync.WaitGroup
  Log      *logger.LoggerType
}

func (ns *NATS) AsyncSubscription(
  natsSubPath  string,
  natsCallback nats.MsgHandler,
) {
  ns.NatsWG.Add(1)
  theSubscription, err := ns.Conn.Subscribe(natsSubPath, natsCallback)
  if err != nil {
    ns.Log.MayBeFatal(
      fmt.Sprintf("Could not register [%s] subscription", natsSubPath),
      err,
    )
  	ns.NatsWG.Done()
  }
  ns.Subs = append(ns.Subs, theSubscription)
}

func (ns *NATS) CloseDown() {
  for _, aSubscription := range ns.Subs {
    _ = aSubscription.Unsubscribe()
    ns.NatsWG.Done()
  }
}

func ConnectServer(
  connections []string,
  cnLog        *logger.LoggerType,
) (*NATS){
  ns := NATS{}
  connectionsStr := strings.Join(connections, ",")
  theConnection, err := nats.Connect(connectionsStr)
  cnLog.MayBeFatal(
    fmt.Sprintf("Could not connect to NATS servers [%s]", connections),
    err,
  )
  ns.Conn = theConnection
  ns.Log  = cnLog
  ns.Subs = make([]*nats.Subscription, 0)
  return &ns
}

func (ns *NATS) RunServer() {
  ns.NatsWG.Wait()
}
