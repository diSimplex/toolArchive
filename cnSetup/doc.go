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

ConTeXt Nurseries Setup command (cnSetup) sets up the configuration 
required by the cnNursery and cnTypeSetting commands as well as each 
user's web browser. 

In particular the cnSetup command:

  1. Creates a private (self-signed) certificate authority (CA).
  
  2. Creates x509 server certificates for the cnNursery server signed by
     the CA. 

  3. Creates x509 client certificates for each user signed by the CA.
  
  4. Creates the (YAML) configuration files used by the cnNursery and 
     cnTypeSetting commands. 

The x509 client certificates are meant to be loaded by each user into 
their web broswer to enable the user to browse the HTML ConTeXt Nursery 
interfaces. 

Usage of cnSetup:
  
  cnSetup [-c|-config string] [-createCA] [-s|-show]

    -c string
      The configuration file to load (default: "nurseries.yaml")
        
    -config string
      The configuration file to load (default: "nurseries.yaml")
        
    -createCA
      IF the Certificate Authority (CA) crt and key files can't be loaded,
      THEN the CA will only be re-created IF the "-createCA" switch is present.
      (default: no CA will be re-created)
        
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

