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

/*

ConTeXt Nurseries command internals
  
The ConTeXt Nurseries Nursery command (cnNursery) runs a ConTeXt Nursery 
on a given machine. 

In particular the cnNursery command:

  1. Runs a combined HTML and RESTfull HTTP/JSON interface which can be used 
     to manage a running cnNursery. 
  
  2. Manages local collections of working files ("workspaces") used by 
     commands being run by the cnNursery on a user's behalf. 
  
  3. Manages a collection of command output to allow users to understand how
     a command is progressing. 

  4. Manages a collection of runable commands ("actions").
  
  5. Allows a registered user to configure and run one or more command 
     actions in a specific workspace. 

This CNNurseries package is used by the cnNursery command to orchestrate 
the creation of local workspaces, and running of local actions, as well as 
(potentially) forwarding "work" to other less heavily loaded cnNurseries 
in a federation of cnNurseries. 

*/
package CNNurseries