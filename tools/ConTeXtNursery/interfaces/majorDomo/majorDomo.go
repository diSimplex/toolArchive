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

// A pair of NATS and RESTful HTTP interfaces which are responsible for 
// managing the location, dependency, and topological sort of one or more
// artifacts.
//
package majorDomo

import (
  "github.com/nats-io/nats.go"
)

//////////////////////////////////////////////////////////////////////////
// majorDomo interface types

type Artifact_Identity struct {
  Name           string
  Type           string
  Location       string
  Time_Stamp     uint32
  File_Hash      string
}

type Artifact struct {
  Artifact       Artifact_Identity
  Dependencies []string
  Copies       []*Artifact_Identity
}

// Listen to all `artifact.have` messages.
//
// Populates the dependency graph with all new or updated meta-data. 
//
func ArtifactHave(msg *nats.Msg) {
  // do nothing
}

// Listen to all `artifact.wants` messages.
//
// Populates the dependency graph with all new or updated meta-data. 
//
func ArtifactWants( msg *nats.Msg) {
  // do nothing
}

// Listen to all `artifact.delet` messages.
//
// Removes the corresponding artifact from the dependency graph.
//
func ArtifactDelete(msg *nats.Msg) {
  // do nothing
}

func AddMajorDomoInterface(natsConn *nats.Conn) {

  nastConn

  artifactWantsSub, err := natsConn.Subscribe("artifact.wants.>", ArtifactWants)
  if err != nil {
  	//??
  }

  artifactDeleteSub, err := natsConn.Subscribe("artifact.delete.>", ArtifactDelete)
  if err != nil {
  	//??
  }
}
