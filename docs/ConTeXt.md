<header><title>Parallelizing the typesetting of ConTeXt documents</title></header>

# Parallelizing the typesetting of ConTeXt documents

## Problem

ConTeXt, being based upon 
[TeX](http://www.tug.org/)/[LuaTeX](http://luatex.org/), is *single 
threaded*. Our goal is to be able to mimic the 
[LaTeX](https://www.latex-project.org/) [packages](https://www.ctan.org/) 
[subdocs](https://ctan.org/pkg/subdocs) 
([github](https://github.com/jbezos/subdocs)), 
[subfiles](https://ctan.org/pkg/subfiles) 
([github](https://github.com/gsalzer/subfiles)), and/or 
[standalone](https://www.ctan.org/pkg/standalone) 
([bitbucket](https://bitbucket.org/martin_scharrer/standalone/src/default/)) 
in ConTeXt (see also [LaTeX/Modular 
Documents](https://en.wikibooks.org/wiki/LaTeX/Modular_Documents)).

## Assumptions

1. We assume all documents will be structured as a tree of sub-documents 
   where **only** the leaf sub-documents of this tree contain any text to 
   be typeset. 


## Questions

1. Should all sub-documents include *explicit* pre/post-ambles?

   If we do not use *explicit* pre/post-ambles, then we expect the 
   typesetter to automatically add pre/post-ambles. This in turn requires 
   someway of *registering* pre/post-ambles for the whole document. 

   What about sub-parts with differing pre/post-ambles?

   At the moment, this would require either 

     - explicit generation of sub-document wrappers
     - explicit ConTeXt module which ignores "un-used" "includes"

   Using explicit generation would allow the system to additionaly 
   update/control page numbers... but is rather "clunky".

2. How do we control the starting page numbers of sub-documents?

3. How do we determine sub-document dependencies?

   "Usually" the "only" sub-document dependency is that required for 
   sequential page numbering. Since ConTeXt requires a multi-pass approach 
   to get this correct anyway, individual sub-documents only depend upon 
   they "previous" sub-document from *previous* ConTeXt passes.

## Solutions




## Issues

There are a number of potential issues:

1. Management of the a sub-document's preamble. Depending upon the 
   complexity of the over all document, this "preamble" might be scattered 
   across a number of files.

   This is what subfiles and standalone attempt to solve.

2. Sharing of the `aux` (TeX) or `tuc` (LuaTeX/ConTeXt) information which 
   includes internal cross references *between* parallelized sub-documents. 

   This is what subdocs attempts to solve.

3. Recombination of the typeset result of ConTeXt run on each of the 
   parallelized sub-documents separately.

   Both [`pdftk`](https://www.pdflabs.com/tools/pdftk-the-pdf-toolkit/) 
   (via [`iText`](https://itextpdf.com/en) 
   ([github](https://github.com/itext))) and `gs` 
   ([ghostscript](https://www.ghostscript.com/)) are able to do this in at 
   least draft form.

4. Recompute the page numbers of each parallelized (sub)document.

   This is part of the problem subdocs attempts to solve via managment of 
   the `aux` or `tuc` files.

5. Complex document structures with intermingled text and 
   sub-documents/chapters.

   Our solution (above) depends upon gluing together pdfs produced in 
   parallel which all start/end of page boundaries. Intermingling text and 
   sub-documents makes it very hard for the parallelization to keep to 
   page boundaries, as well as identify "whole" collections of pages.

   This suggests the use of a pre-processor to break the original document 
   up into a tree of sub-documents where **only** the leaf nodes of the 
   tree actually have text. Sometime in the future ConTeXt itself may be 
   able to do this..

   At the moment we *expect* the user to **only** create documents with 
   this "nice" tree structure. 
