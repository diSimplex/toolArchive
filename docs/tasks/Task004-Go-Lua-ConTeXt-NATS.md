<header><title>Task 004 - Go, Lua, ConTeXt, NATS</title></header>

# Task 004 - Showing that GoLang, Lua, ConTeXt and NATS can interoperate

1. **completed** Run Lua scripts inside GoLang
   1. **completed** using simple Lua 5.4 ( https://www.lua.org/download.html )
   2. **completed** using luaFileSystem ( https://github.com/keplerproject/luafilesystem )
   3. **completed** using luaSockets ( https://github.com/diegonehab/luasocket )
   4. **completed** JSON: 
      1. **completed** using lua-cjson ( https://github.com/mpx/lua-cjson )
      2. **completed** using lunajson ( https://github.com/grafi-tt/lunajson )
      3. **completed** using dkjson ( http://dkolf.de/src/dkjson-lua.fsl/home )
   5. **completed** using uuid (lua) ( https://github.com/Tieske/uuid )
   6. **completed** using lua-nats ( https://github.com/DawnAngel/lua-nats ) using all of lua-cjson, lunajson and dkjson
2. **completed** Run Lua scripts inside ConTeXt
   1. **completed** using simple Lua 5.4 ( https://www.lua.org/download.html )
   2. **completed** using luaFileSystem ( https://github.com/keplerproject/luafilesystem )
   3. **completed** using luaSockets ( https://github.com/diegonehab/luasocket )
   4. **completed** JSON: 
      1. **FAILED** using lua-cjson ( https://github.com/mpx/lua-cjson ) missing `lua_checkstack`.
      2. **completed** using lunajson ( https://github.com/grafi-tt/lunajson )
      3. **completed** using dkjson ( http://dkolf.de/src/dkjson-lua.fsl/home )
   5. **completed** using uuid (lua) ( https://github.com/Tieske/uuid )
   6. **completed** using lua-nats ( https://github.com/DawnAngel/lua-nats ) using both lunajson and dkjson
3. **completed** (re)setup OCI-repository on pi01
4. (re)Run NATS inside a podman pod
5. (re)Ensure NATS is federated across muplitiple machines/pods.
6. Issue NATS messages from Lua scripts running inside ConTeXt
7. Issue NATS messages from Lua scripts running inside GoLang
8. Issue NATS messages directly from GoLang
 
## Architecture

**Run Lua scripts inside GoLang**

We will use https://github.com/xiexiao/golua since it already adds the 
-ldl to allow Lua's use of shared libraries on linux machines.

**(re)Run NATS inside a podman pod**


## Problems

1. The current lua-nats interface is very old *but* still works with NATS 
   v2.0 for subscription and publishing on simple channels with simple 
   messages. 

   As an alternative, for GoLang we can and should simply use the GoLang 
   NATS client. 

   As an alterntive, for ConTeXt, if we have through-put problems with the 
   pure lua JSON implementations, we might use the C NATS client 
   https://github.com/nats-io/nats.c which is up to date. 

2. ConTeXt requires the `--permitloadlib` option to allow external Lua 
   shared modules to be loaded. However some `lua_` symbols seem to be 
   missing. 

3. ConTeXt (even with `--permiteloadlib`) seems to be missing the 
   `lua_checkstack` symbol (among possibly others). This means that the 
   Lua require of `cjson` fails. We can use either of the pure lua JSON 
   tools: http://dkolf.de/src/dkjson-lua.fsl/home (which uses lpeg for 
   speedup) or https://github.com/grafi-tt/lunajson (which uses its own 
   internal optimization). 

4. We need a simplified version of docker-compose probably written in 
   Lua/YAML (since ConTeXt provides Lua for scritps). Use 
   https://github.com/exosite/lua-yaml for the Lua-YAML loading. (We can 
   not use wrapper of LibYAML since that requires `--permitloadlib` which 
   does not yet work soon.... but not yet) 

## Reflections

1. golua changes `pcall` to `unsafe_pcall` and `xpcall` to `unsafe_xpcall`.
   These unsafe versions will FAIL spectacularly if they call back into go.
   They are OK if they only interact with Lua.

   We will wrap any use of `L.Dofile` in the `L.GetGlobal("unsafe_pcall") 
   ; L.SetGlobal("pcall")` and `L.PushNil() ; L.SetGlobal("pcall")` pair 
   in GoLang in order to allow the loaded files to use `pcall` as normal. 

2. For future: see ntg-context email archive for the subjects:
   - "Support for musl"
   - "Three problems with ConTeXt standalone for armhf"
   - "Best way to create a large number of documents from  database"
   for a discussion on how to compile inside Alipine linux.

3. **2020-11-27:** Searched through context email archive for any 
   discussion on parallelizing large documents. None could be found. 
   (There was discussion on parallelizing lots of small documents driven 
   from databases). 

   Hans seemed to suggest that the greatest time (in LuaMetaTex) is now 
   taken up by the Lua code rather than the old TeX code. I think he also 
   seemed to suggest that a lot of this time is taken loading/manipulating 
   fonts/images (so in the pdf generation?). 
   
