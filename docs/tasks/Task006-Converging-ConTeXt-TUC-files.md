<header><title>Task 006 - Converging ConTeXt TUC files</title></header>

# Task 006 - Converging ConTeXt TUC files

0. Patch mtx-context.lua script to ensure correct `--once` behaviour with 
   start page control. 

1. Identify page numbers inside the TUC file.

   - run a simple tex document with different starting page numbers and 
     compare the differences. 

2. Start writing a GoLang/Lua based TUC coordinator

3. Start writing a GoLang/ConTeXt based worker

4. Transfer TUC files to/from worker and coordinator

## Architecture



## Problems

1. Controlling page numbers of individual sub-documents

## Questions

1. How will the `*.tuc` files to transfered between the coordinator and 
   the workers? 

2. What is the structure of the TUC file?

2. How will the coordinator choose which values to update?

2. How will the coordinator keep track of which TUC values depend upon 
   which sub-documents? 

## Reflections

Reading the mtx-contex.lua script only the `*.tuc` is used/monitored for 
multi-pass information. 
