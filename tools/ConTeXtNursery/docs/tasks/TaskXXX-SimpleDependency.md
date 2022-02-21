<header><title>Task 005 - Simple Dependency</title></header>

# Task 005 - Simple Dependency 

1.

2.

## Architecture

- **`projectManager`** (*interface*): 

  - Automates the 'build' of a given artifact by performing a [depth first 
  topological 
  sort](https://en.wikipedia.org/wiki/Topological_sorting#Depth-first_search) 
  of the distributed dependency graph contained in the system's collection 
  of `librarian` databases. 

- **`cnProjectManager`** (*microservice*) implements the `projectManager` 
interface.

- **`builder`** (*interface*):

  - Builds a given artifact using particular (gcc, context, joylol, 
  pdf2htmlEX, etc) build tool. 

  - Provides a NATS message interface using the following subjects:

    - Listens for `artifact.build` messages, responding with an offer to 
    build (which contains some measure of the local cost of the build -- 
    such as how many of the dependencies it already has). 

    - Listens for `artifact.dependencies` messages, responding with any 
    dependency information computed for a given artifact. 

- There will be a range of distinct *microservers* which each implement 
the `builder` interface in different ways. 

## Reflections
