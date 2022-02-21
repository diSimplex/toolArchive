# ConTeXt Nursery
A Nursery for Typesetting large [ConTeXt](https://www.contextgarden.net) 
documents in parallel.

## Problem

Over the next few years, in order to communicate what is in *my* head, I 
will need to write many thousands of pages of Mathematical and 
Computational argumentation. As part of this project, JoyLoL code will be 
embedded directly in the documents. More importantly the act of typesetting 
the documents will invoke a [JoyLoL 
system](https://github.com/diSimplex/JoyLoLComputeFarm) to interpret, 
compile, and check the correctness of the code and examples.

TeX, and hence LaTeX and ConTeXt, as well as many compilers are essentially 
single threaded. With thousands of pages and lots of code, this means a lot 
of waiting to see the effect of small changes as I write.

While compilation of a *single* code file is essentially single threaded, 
*building* a large project is an "embarassingly" parallel problem. Since 
the *compilation* and the *linking* are kept as separate tasks, the bulk of 
the compilation can be done in parallel, using as many CPU threads as are 
available.

More importantly since the results of the compilation of a given code file 
will be linked in a later step, if neither the code file nor any of its 
dependencies change, then the compilation can simply be skipped, using no 
CPU threads what so ever.

At the moment, TeX, LaTeX, and ConTeXt, can not enjoy the same benefits of 
multiple CPU threads. However, *if* we could break a large document into 
its individual chapters, which are known to start on new page boundaries, 
and then find some way to combine them back together, we *could* typeset 
each chapter in parallel, and hence potentially make use of all available 
CPU threads.

Similarly if we can be assured that a given ConTeXt chapter and none of its 
dependencies have changed, then with a typeset and combine model, such 
chapters do not need to be re-typeset, again using no CPU threads what so 
ever.

The ConTeXt Nursery project aims to solve exaclty this problem using the 
[ConTeXt typesetting system](https://www.contextgarden.net) and the
associated [LuaTeX](http://luatex.org/).

## Details

More details can be found in the [docs](docs) directory.

## License

Unless explicitly stated otherwise, all content contained in this repository is

```
          Copyright (C) 2020 the contributors to the diSimplex project
               as currently listed by https://github.com/diSimplex
```

and is licensed under the Apache License, Version 2.0 (the "License"); you 
may not use this code except in compliance with the License. You may obtain 
a copy of the License at

```
                    http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software 
distributed under the Apache 4.0 License is distributed on an **‘AS IS’
BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND**, either express or 
implied. See the Apache 2.0 License whose URL is listed above for the 
specific language governing permissions and limitations under the License.

A local copy of the Apache License version 2.0 can be found in the 
[LICENSE.txt](LICENSE.txt) file.
