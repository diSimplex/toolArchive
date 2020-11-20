<header><title>User tool use case</title></header>

# User tool use case

## Objective

Our objective is to allow the user to select a ConTeXt document and a 
given "deliverable" for typesetting. 

## User's actions

1. User logs into their local User's Tool (acting as an ArtifactManager). 

2. User selects a particular \*.tex *master* document and requests the 
   html version from a list of "deliverables".

## System's actions

0. The User's Tool loads a (YAML?/LUA?) configuration file and scans its 
   local working directories looking for 

   1. ConTeXt archives / code archives / git repositories
   2. Dependency descriptions of particular deliverables.

1. The User's Tool advertises for a TaskManager.

2. Any free TaskManagers respond with offers to manage a task.

3. The User's Tool chooses a TaskManager and using the Lua dependency DSL 
   associated with any user chosen "deliverables", informs the Task 
   Manager of the inital dependency graph. 

4. The chosen TaskManager begins managing the dependency graph to 
   ultimately build the requested goals. 

## Configurtion

The User's tool is configured by specifying a (collection of) Lua script(s).

These Lua scripts define a number of "deliverables" from which the use can 
choose one or more items.

Corresponding to each "deliverable" is an associated Lua dependency DSL 
which the User's Tool uses to inform the choosen TaskManager of the 
initial dependency graph together with required goals. 

## Problems

### Interpreting Lua

See: https://github.com/aarzilli/golua or? https://github.com/xiexiao/golua

https://github.com/stevedonovan/luar/

https://github.com/fiatjaf/lunatico

### Interacting with the NATS messageing service

See: https://github.com/nats-io/nats.go

See: https://github.com/DawnAngel/lua-nats

### Serving files

There will be two types of configuration:

1. Global configuration of repositories.

2. Local configuration inside each repository.

The Global configuration will consisit of a collection of YAML files listing:

1. Repositories (either directories, tar/zip files, or git repositories).

2. Either individual or directories of other global configuration files.

The Local configuration files will describe the top-level files and their 
associated deliverables.

Will use Lua scripts for local configuration and YAML or Lua for global 
configuration. 

See: https://blog.gopheracademy.com/advent-2016/go-syntax-for-dsls/
  http://lua-users.org/lists/lua-l/2006-03/msg00259.html
  http://lua-users.org/lists/lua-l/2006-03/msg00265.html

  https://martinfowler.com/articles/rake.html
  https://en.wikipedia.org/wiki/SCons
  https://stackoverflow.com/a/19182835
  https://dnaeon.github.io/choosing-lua-as-the-ddl-and-config-language/

YAML replacement?!? https://github.com/hashicorp/hcl

### Dependency analysis

**ANSI-C:** see `-M` (and friends) in:
https://gcc.gnu.org/onlinedocs/gcc-10.2.0/gcc/Preprocessor-Options.html#Preprocessor-Options 

**ConTeXt:** see ???

We will use the *new* (and improved?) LMTX version of ConTeXt (which uses 
even more Lua internally). 

Our own ConTeXt modules which have a dependency mode to only output a 
dependency file as a Lua script. With the `context --once` command line. 
See page 9 of http://www.pragma-ade.com/general/manuals/tools-mkiv.pdf . 

See also: http://www.pragma-ade.com/general/manuals/luametatex.pdf


