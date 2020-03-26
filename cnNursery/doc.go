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

A given cnNursery interacts with a user's browser, a user's cnTypeSetter 
or other cnNurseries in a federdation of cnNurseries using doubled ended 
https. This means that all browsers, cnTypeSetters, and cnNurseries, on 
both ends of a given https connection, are required to present x509 TLS 
Certificates signed by the Certificate Authority manged by a cnSetup 
command. 

The configuration used to configure a cnNursery is created and managed by 
the cnSetup command. 

Usage of cnSetup:
  
  cnSetup [-c|-config string] [-createCA] [-s|-show]

    -c string
      The configuration file to load (default: "nurseries.yaml")
        
    -config string
      The configuration file to load (default: "nurseries.yaml")
        
    -s	
      Show the loaded configuration (default: configuration not listed)
    
    -show
      Show the loaded configuration (default: configuration not listed)

  
------------------------------------- SECURITY -------------------------------------

Since a ConTeXt Nursery is meant to run (nearly arbitrary) commands on 
each machine running a cnNursery, the SECURITY of the Fedeartion is 
important. 

All interactions between the ConTeXt Nursery components, cnNursery, 
cnTypeSetter as well as user browsers, REQUIRES double ended TLS.

This means that in all cases the "server" authenticates the "client" AND 
the "client" authenticates the "server", in all cases using x509 
certificates signed by the private (self-signed) CA created by the cnSetup 
command. 

HOWEVER, this security is ONLY AS SECURE AS YOUR WEAKEST USER.

The CA, and Server x509 keys ARE NOT PASSWORD encrypted. Your security 
DEPENDS upon the security of your file system. In particular any one with 
access to the CA key can create their own client/server x509 certificates 
and then interact with the ConTeXt Nursery federation. 

It is assumed that the CA as well as the cnNursery certificates and 
configuration are under the control of a ConTeXt Nursery administrator who 
has a reasonable understanding of the security requirements of their 
federation's use. 

Each user's PKCS12 file (which is loaded into each user's web browser) IS 
password encrypted. The passwords for each user are contained in the users 
password file typically loacted in the "user/passwords" file. 


*/
package main 

