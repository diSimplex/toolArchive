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

package logger

import (
  "encoding/json"
  "log"
  "runtime/debug"
)

/////////////////////////////
// Logging and Error handling
//

type LoggerType struct {
  AppName    string
  PrintStack bool
}

func CreateLogger(appName string) *LoggerType {
  return &LoggerType{ AppName: appName }
}

func (l *LoggerType) SetPrintStack(printStack bool) {
  l.PrintStack = printStack
}

func (l *LoggerType) MayBeFatal(logMessage string, err error) {
  if err != nil {
    if l.PrintStack { debug.PrintStack() }
    log.Fatalf("%s(FATAL): %s ERROR: %s", l.AppName, logMessage, err)
  }
}

func (l *LoggerType) MayBeError(logMessage string, err error) {
  if err != nil {
    if l.PrintStack { debug.PrintStack() }
    log.Printf("%s(error): %s error: %s", l.AppName, logMessage, err)
  }
}

func (l *LoggerType) Log(logMesg string) {
  log.Printf("%s(info): %s", l.AppName, logMesg)
}

func (l *LoggerType) Logf(logFormat string, v ...interface{}) {
  log.Printf(l.AppName+"(info): "+logFormat, v...)
}

func (l *LoggerType) Json(logMesg string, valName string, aValue interface{}) {
  jsonBytes, err := json.MarshalIndent(aValue, "", "  ")
  if err != nil {
    l.MayBeError("Could not marshal "+valName+" into json", err)
    jsonBytes = make([]byte, 0)
  }
  log.Printf("%s(json): %s %s", l.AppName, logMesg, string(jsonBytes))
}

