<header><title>Task Manager</title></header>

# Task Manager

## Problem

We need a "Task Manager" to manage the parallelization of inter-dependent 
tasks. 

"Classically" this is done assuming the inter-dependent tasks are an 
acyclic-tree of dependencies. With this assumption the task manager "just" 
has to do a ["topological 
sort"](https://en.wikipedia.org/wiki/Topological_sorting) (often by using 
a "reverse" recursive descent). 

With the use of ConTeXt (and possibly JoyLoL), this assumption of an 
acyclic-dependency-tree is broken. This means the task manager can no 
longer do a "simple" topological sort.

Since the inter-dependencies contain explicit dependency cycles, we must

1. identify dependency cycles

2. selectively un-wind these cycles back into a (truncated) acyclic-tree 
   of dependencies. 



## Questions

## Solution

1. Allow dependencies to explicitly depend upon a task at a particular 
   un-winding depth ("time"). Such dependencies will only be expressed 
   *relatively*. (Note that *absolute* references if used make the overall 
   task search structure more "fragile"). 


## Issues

1. Most 

2. Un-winding dependency cycles produces an infinite acyclic-tree of 
   dependencies. The fundamental problem with this is that the "top" of 
   topological sort for the dependency tree is located at the "infinity". 

   This means we *must* work with a truncated un-wound dependency tree, 
   where the un-winding depth is *pre-determined*. 

   How do we detect that the current un-winding depth is insufficient?

3. The dependencies of un-wound tasks must also exist at multiple 
   un-winding levels. 
