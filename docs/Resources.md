# Resources

### Development processes

- [zeroMQ's C4](https://rfc.zeromq.org/spec/22/)

### LaTeX discussions of how and when to break up texts

- [Splitting a large document into several 
files](https://tex.stackexchange.com/questions/29577/splitting-a-large-document-into-several-files)

- [Writing and Managing Thesis in 
LaTeX](https://tex.stackexchange.com/a/29534)

- [On managing large 
documents](https://texblog.org/2016/06/20/on-managing-large-documents/)

- [Keeping things organized in large 
documents](https://texblog.org/2012/12/04/keeping-things-organized-in-large-documents/) 
discusses how to use 
[catchfilebetweentags](https://www.ctan.org/pkg/catchfilebetweentags?lang=en) 
package

- [Techniques and packages to keep up with good 
practices](https://tex.stackexchange.com/questions/19264/techniques-and-packages-to-keep-up-with-good-practices/20191#20191)

- [Everyday LaTeX and workflow?](https://tex.stackexchange.com/a/22433)

- [overleaf: Management in a large 
project](https://www.overleaf.com/learn/latex/Management_in_a_large_project) 
see the section on the [import package](https://www.ctan.org/pkg/import)...

- [overleaf: Multi-file LaTeX 
projects](https://www.overleaf.com/learn/latex/Multi-file_LaTeX_projects)

- [wikiBooks: LaTeX/Modular 
Documents](https://en.wikibooks.org/wiki/LaTeX/Modular_Documents) see: 
[Separate compilation of child 
documents](https://en.wikibooks.org/wiki/LaTeX/Modular_Documents#Separate_compilation_of_child_documents)

- [LaTeX Notes: Structuring Large Documents (include vs 
input)](https://web.science.mq.edu.au/~rdale/resources/writingnotes/latexstruct.html)

- [Splitting a Large Document into Several Files (include vs 
input)](https://www.dickimaw-books.com/latex/thesis/html/include.html)

### Combining PDFs with embedded hyper-links

- [Merge PDF's with PDFTK with 
Bookmarks?](https://stackoverflow.com/questions/2969479/merge-pdfs-with-pdftk-with-bookmarks/3139897) 
- [extractpdfmark](https://github.com/trueroad/extractpdfmark) (in Ubuntu)

- [python script to create PDFMarks](https://stackoverflow.com/a/30524828)

- [Scripts to merge pdfs and add bookmarks with pdftk 
](https://autohotkey.com/board/topic/98985-scripts-to-merge-pdfs-and-add-bookmarks-with-pdftk/) 

- [`pdftk`](https://www.pdflabs.com/tools/pdftk-the-pdf-toolkit/) (pdftk 
  uses iText under the covers)

- [`iText`](https://itextpdf.com/en) ([github](https://github.com/itext)))

- `gs` ([ghostscript](https://www.ghostscript.com/))

- there are some comments that the combined PDF might have a lower quality 

- consider using iText (java)

- [Need to merge multiple pdf's into a single PDF with Table Of Contents 
sections](https://stackoverflow.com/questions/2418871/need-to-merge-multiple-pdfs-into-a-single-pdf-with-table-of-contents-sections/40222656#40222656)

- [Add and edit bookmarks to 
pdf](https://unix.stackexchange.com/questions/17065/add-and-edit-bookmarks-to-pdf/31070)

- [using pure ghostscript](https://stackoverflow.com/a/16027780)

- Using pdftk-java: `pdftk in1.pdf in2.pdf cat output out1.pdf`

### Overleaf (online LaTeX)

- [overleaf](https://www.overleaf.com/)

- [overleaf github](https://github.com/overleaf/overleaf)

- [Overleaf CLSI](https://github.com/overleaf/clsi) also has an API to 
  collect and/or distribute artifacts. It works on a "PULL" model, in that 
  artifacts are pulled from some (web)server as and when needed.

### Distributed (C) compilers

- [distcc](https://github.com/distcc/distcc) Both distcc (and its 
  successor, icecream) assume you have already broken the problem into 
  embarassingly parallel components (of individual C-files) and "simply" 
  want to distribute the compilation across multiple machines (instead of 
  make -j's multiple cpus on a single machine).

- [icecream](https://github.com/icecc/icecream)

- [recc](https://gitlab.com/bloomberg/recc) recc is the Remote Execution 
  Caching Compiler. It is a cross between ccache and distcc using the 
  Remote Execution APIs.

### Code artifact caching

- [ccache](https://ccache.dev/) Ccache is a (distributed?) compilation 
  caching system designed specifically for C-code. Ccache acts as a "proxy" 
  for a given C-compiler (such as gcc) and analysing the input files and 
  command line options can determine if the result of this compilation has 
  already been cached...

- [recc](https://gitlab.com/bloomberg/recc) recc is the Remote Execution 
  Caching Compiler. It is a cross between ccache and distcc using the 
  Remote Execution APIs.

- [Overleaf](https://github.com/overleaf/overleaf#other-repositories) has a 
  collection of artifact caching and/or storing services.

- [Overleaf CLSI](https://github.com/overleaf/clsi) also has an API to 
  collect and/or distribute artifacts. It works on a "PULL" model, in that 
  artifacts are pulled from some (web)server as and when needed.

### Distributed build systems

Both bazel and gradle require a pre-known dependency analysis in the form 
of a "build (dependency) script". Bazel seems to use a "higher-level" 
description of the "compilation" problem.

They both assume your have previously split your problem into individual 
"parts" and have a tool to combine the results into the "final" "whole".

- [bazel](https://bazel.build/) 
  ([Starlark](https://github.com/bazelbuild/starlark/)) (extensions based 
  on Starlark (and goLang?))

- [gradle](https://gradle.org/) (extensions based on Java)

- [Magefile](https://magefile.org/) a simple make/rake in goLang

- [justfile](https://github.com/casey/just) a simple make/rake in rust

- [Modern Make](https://github.com/tj/mmake) a make wrapper

### Distributed build farms

**All of the following distributed build farms, implement or use the next two 
google apis:**

- Google's [Remote Execution 
  API](https://github.com/bazelbuild/remote-apis) (based on goLang)

- Google's [Remote Workers 
  API](https://docs.google.com/document/d/1s_AzRRD2mdyktKUj2HWBn99rMg_3tcPvdjx3MPbFidU/edit#heading=h.1u2taqr2h940).

Both of these APIs require Googles's 
[ProtoBuf](https://developers.google.com/protocol-buffers/) tool. There is 
a [Lua ProtoBuf](https://github.com/starwing/lua-protobuf).

**Each of these build farms uses the above two google APIs:**

- [buildbarn (original)](https://github.com/EdSchouten/bazel-buildbarn) see 
  the very useful process interaction diagrams. (based on goLang)

- [BuildBarn (cluster of github projects)](https://github.com/buildbarn) 
  see [Example deployments of 
  Buildbarn](https://github.com/buildbarn/bb-deployments) for useful 
  process interaction diagrams. (based on goLang)

- [Scoot](https://github.com/twitter/scoot)

- [BuildStream](https://www.buildstream.build/)

- [BuildGrid](https://gitlab.com/BuildGrid/buildgrid) BuildGrid is a Python 
  remote execution service which implements Google's [Remote Execution 
  API](https://github.com/bazelbuild/remote-apis) and the [Remote Workers 
  API](https://docs.google.com/document/d/1s_AzRRD2mdyktKUj2HWBn99rMg_3tcPvdjx3MPbFidU/edit#heading=h.1u2taqr2h940).

- [Buildfarm](https://github.com/bazelbuild/bazel-buildfarm)

### Asynchronous wire protocols

**zeroMQ**

- [zeroMQ documentation](https://zeromq.org/get-started/)

- [zGuide](http://zguide.zeromq.org/page:all)

- [zeroMQ MajorDomo project](https://github.com/zeromq/majordomo)

- [zeroMQ specifications](https://rfc.zeromq.org/)

- [zeroMQ gitHub](https://github.com/zeromq/)

**NanoMsg and Scalability Protocols**

- [A Look at Nanomsg and Scalability Protocols (Why ZeroMQ Shouldn’t Be 
  Your First 
  Choice)](https://bravenewgeek.com/a-look-at-nanomsg-and-scalability-protocols/)

- [nanomsg documentation](https://nanomsg.org/documentation.html)

- [nanomsg github](https://github.com/nanomsg/nanomsg)

**nng**

- [Rationale: Or why am I bothering to rewrite 
  nanomsg?](https://nng.nanomsg.org/RATIONALE.html)

- [nng](https://nng.nanomsg.org)

- [nng github](https://github.com/nanomsg/nng)

**RESTful protocols** (used by overleaf's CLSI)

- [Representational state transfer 
  (REST)](https://en.wikipedia.org/wiki/Representational_state_transfer)

- [REST cookbook](http://restcookbook.com/)

- [RESTful 
  TutorialsPoint](https://www.tutorialspoint.com/restful/index.htm)

- [What is the difference between REST and HTTP 
  protocols?](https://stackoverflow.com/questions/5449034/what-is-the-difference-between-rest-and-http-protocols)

- [REST API Tutorial](https://restfulapi.net/)

- [Representational State Transfer (REST) (Roy Fielding's PhD chapter 
  5)](https://www.ics.uci.edu/~fielding/pubs/dissertation/rest_arch_style.htm)

- [Roy Fielding's PhD 
  thesis](https://www.ics.uci.edu/~fielding/pubs/dissertation/top.htm)

**zeroMQ / NanoMsg / nng VS Rest**

- [What are the benefits of using a ZeroMQ-like messaging library over a 
  REST API based 
  architecture?](https://www.quora.com/What-are-the-benefits-of-using-a-ZeroMQ-like-messaging-library-over-a-REST-API-based-architecture)

- [A Protocol for REST over ZeroMQ](http://hintjens.com/blog:86)

- [40/XRAP Extensible Resource Access 
  Protocol](https://rfc.zeromq.org/spec/40/)
