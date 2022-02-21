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

ConTeXt Setup command internals.
  
The cnSetup command consists of methods to manage:
  
    1. Configuraiton (ConfigType)
    2. CertificateAuthority (CAType)
    3. Nursery Certificates and Configuration (NurseryType)
    4. User Certificates and Configuration (UserType)
  
This CNSetup package is used by the cnSetup command to orchestrate the 
creation of a Certificate Authority, as well as Certificates and 
Configuration for each Nursery and User. 

*/
package CNSetup