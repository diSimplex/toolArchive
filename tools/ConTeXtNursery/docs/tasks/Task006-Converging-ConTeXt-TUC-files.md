<header><title>Task 006 - Converging ConTeXt TUC files</title></header>

# Task 006 - Converging ConTeXt TUC files

0. **not needed** Patch mtx-context.lua script to ensure correct `--once` 
   behaviour with start page control. Use `--runs=1` instead. 

1. **completed** Identify page numbers inside the TUC file.

   - run a simple tex document with different starting page numbers and 
     compare the differences. 

2. Start writting a simple Ruby script which simulates splitting the 
   luametatex document into type set sub-documents. 

3. Start writing a GoLang/Lua based TUC coordinator

4. Start writing a GoLang/ConTeXt based worker

5. Transfer TUC files to/from worker and coordinator

## Architecture



## Problems

1. Controlling page numbers of individual sub-documents.

2. Controlling the chapter/section number of individual sub-documents.

## Questions

1. How will the `*.tuc` files to transfered between the coordinator and 
   the workers? 

2. What is the structure of the TUC file?

2. How will the coordinator choose which values to update?

2. How will the coordinator keep track of which TUC values depend upon 
   which sub-documents? 

2. Where should new ConTeXt modules be put? In one of the following 
   places: 

   - `$HOME/texmf`
   - `$CONTEXT/tex/texmf-local`
   - `$CONTEXT/tex/texmf-projects`

   *Answer*: any of the above.
   
## Reflections

Reading the mtx-contex.lua script only the `*.tuc` is used/monitored for 
multi-pass information. 

A deeper exploration of the ConTeXt code (both TeX and Lua) shows that 
while it *should* be possible to reproduce the *effective* behaviours of 
sub-documents, it will take the knowledge of a ConTeXt wizard, of which 
there are probably only a handful in the world. While I *could* (probably) 
get to that standard, I have other things to do. Equally importantly 
getting to a wizard standard would almost certainly require access to the 
`LuaMetaTex` source code (which is not yet generally available). 

At the moment I *can* fake reasonable facimiles of the sub-document PDFs, 
alas this lacks cross-references in any meaningful way. It *might* be 
possible to even fake these cross-references if I find a way to create the 
missing pages as blanks. The real problem here is that it is the 
(relatively invisible) `realpages` not the (visible) `userpages` which are 
the underlying unit for all of the TUC cross-references. 

## Resources

- https://wiki.contextgarden.net/Command/env
- https://wiki.contextgarden.net/Commands_with_KeyVal_arguments

