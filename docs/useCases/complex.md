<header><title>Complex use case</title></header>

# Complex use case

## Objective

Our objective is to typeset/compile/run the JoyLoL implementation 
documentation. 

This requies:

1. Typesetting the document as both pdf and html (using pdf2htmlEx).

2. Extracting and then compiling the source code for joylol.

3. Running joylol snippets from inside the ConTeXt document.

## User's actions

1. User logs into their local tool

2. User selects a particular \*.tex ConTeXt *master* document and requests 
   the html version. 

## System's actions

## Assumed components

1. User's tool

2. Federated messaging system (assumed and working in the background, that 
   is we do not want nor need to focus upon the lower-level messageing). 

3. Task manager (central resouce for a particular task request)

4. ConTeXt typesetter (can be called in parallel)

5. ConText compositor (central resource managing central update of one 
   particular document's \*.aux file data (one per task request)) 

6. GCC compiler (can be called in parallel)

7. GoLang compiler (can ONLY be called sequentially)

8. JoyLoL compiler/interpreter (should be callable in parallel)

9. pdf2htmlEX compiler (can ONLY be called sequentially? or can we stich 
   together runs of smaller pdfs? ) 

## Problems

### Expanding partially specified requests

**Question:** Given the bare minimum of the user's request, which 
component interprets it? 

**Answer:** (In principle) the component "most" responsible for the 
result. 
