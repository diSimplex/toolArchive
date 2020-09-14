<header><title>Task 004 - Simple Librarian</title></header>

# Task 004 - Simple Librarian

1. Start implementation of the `librarian` message and http interfaces.

2. Start implementation of the `cnLibrarian` microservice.

3. Start implementation of the `cnClient` microservice.

4. A user should be able to publish the existence of a collection of 
artifacts using a `cnClient` microservice.

5. A user should be able to 'typeset' a given artifact using the 
`cnTypeSetter` microservice. 

## Architecture

- **`librarian`** (*interface*):

  - uses [LiftBridge](https://liftbridge.io/) to ensure all (federated) 
  `artifact.have` messages are consistent. 
    
  - listens for all `artifact.have.>` messages, recording where to obtain 
  a given artifact. 
    
  - listens for all `artifact.wants.>` messagse.

  - uses randomized reply intervals (similar to those used by the 
  [Raft](https://raft.github.io/raft.pdf) leadership election"
  algorithm) to "choose" a `librarian` to respond to `artiface.wants` 
  requests.

    - If another `librarian` responds before the end of the given 
    `librarian`'s random reply interval, the given `librarian` remains 
    silent. 
  
    - If no other `librarian` responds before the end of the given 
    `librarian`'s reply interval, then the given `librarian` responds. 

    - Which ever `librarian` responds, it responds with the information it 
    knows about in anwser to the `artifact.wants` question. 

    - The `librarian` which responds, sets its random reply interval to 
    zero, all other `librarians` set their random reply intervals to a 
    small but strictly non-zero value. This ensures the last `librarian` 
    to reply will, unless it crashes, reply with no delay.

    - If more than one `librarian` responds to the same request, then 
    *all* `librarians` reset their random reply interval to a small but 
    strictly non-zero value. This ensures a *single* new leader will 
    be elected in following rounds.

    - Since all `librarians` have the *same* information about all 
    artifacts, if multiple `librarians` respond the client will simply get 
    redundant information which it should ignore. This reply protocol will 
    however guarantee a client gets at least one reply. 
    
  - provides an http route which clients can use to `GET` a given 
  artifact. This interface may respond with http code 303 'see other' to 
  redirect a client to another `cnLibrarian` in the federation. 

    This http route acts as a reverse proxy server to other cnNurser 
    microservices inside a given podman pod. 

  - provides an http route which a user can use, via a browser, to 
  monitor what is known about a given resource, 

  - (may) store all `artifact.have` message meta-data in a local SQLite 
  data base for later use/searching. 

- **`cnLibrarian`** (*microservice*): implements the above `librarian` 
interface. 

- **`cnClient`** (*microservice*):

  - implements (among other interfaces) the `librarian` interface. 

  - provides a tool for the user to inject `artifact.have` messages into 
  the cnNursery system, detailing where to find original source artifacts. 

  - provides a tool for the user to inject `artifact.want` messages into 
  the cnNursery system, detailing what needs to be built.

  - implements an http route which the user can use to list all files in a 
  given directory as `artifact.have`. 

  - implements an http route which the user can use to request a given 
  artifact to be built. 

## Reflections

### Ideas

- Is there a distinction between the 'request' (by an end user) using a 
http/RESTful interface and the 'messages' used by the underlying 
distributed system to communicate with itself? 

### Missing

- we are missing a way to create workspaces/projects (which contain 
collections of related artifact).

- we are missing a way to delete one or more artifacts.

- we are missing a way to delete a workspace/project.
