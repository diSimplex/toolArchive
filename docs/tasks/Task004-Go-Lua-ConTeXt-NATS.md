<header><title>Task 004 - Go, Lua, ConTeXt, NATS</title></header>

# Task 004 - Showing that GoLang, Lua, ConTeXt and NATS can interoperate

1. Run Lua scripts inside GoLang
   1. **completed** using simple Lua 5.4 ( https://www.lua.org/download.html )
   2. **completed** using luaFileSystem ( https://github.com/keplerproject/luafilesystem )
   3. **completed** using luaSockets ( https://github.com/diegonehab/luasocket )
   4. **completed** using lua-cjson ( https://github.com/mpx/lua-cjson )
   5. **completed** using uuid (lua) ( https://github.com/Tieske/uuid )
   6. using lua-nats ( https://github.com/DawnAngel/lua-nats )
2. Run Lua scripts inside ConTeXt
   1. using simple Lua 5.4 ( https://www.lua.org/download.html )
   2. using luaFileSystem ( https://github.com/keplerproject/luafilesystem )
   3. using luaSockets ( https://github.com/diegonehab/luasocket )
   4. using lua-cjson ( https://github.com/mpx/lua-cjson )
   5. using uuid (lua) ( https://github.com/Tieske/uuid )
   6. using lua-nats ( https://github.com/DawnAngel/lua-nats )
3. (re)Run NATS inside a podman pod
4. (re)Ensure NATS is federated across muplitiple machines/pods.
5. Issue NATS messages from Lua scripts running inside ConTeXt
6. Issue NATS messages from Lua scripts running inside GoLang
7. Issue NATS messages directly GoLang
 
## Architecture

**Run Lua scripts inside GoLang**

We will use either https://github.com/aarzilli/golua or 
https://github.com/xiexiao/golua 

**(re)Run NATS inside a podman pod**


## Problems

The current lua-nats interface is very old and may not work with NATS v2.0.

For GoLang we can and should simply use the GoLang NATS client.

For ConTeXt we might use the C NATS client 
https://github.com/nats-io/nats.c which is up to date. 

## Reflections

1. golua changes `pcall` to `unsafe_pcall` and `xpcall` to `unsafe_xpcall`.
   These unsafe versions will FAIL spectacularly if they call back into go.
   They are OK if they only interact with Lua.
