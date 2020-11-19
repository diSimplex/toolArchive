<header><title>Complex use case</title></header>

# Complex use case

## Objective

Our objective is to typeset/compile/run the JoyLoL implementation 
documentation. 

This requires:

1. Typesetting the document as both pdf and html (using pdf2htmlEx).

2. Extracting and then compiling the source code for joylol.

3. Running joylol snippets from inside the ConTeXt document (which changes 
   the contents of the document). 

## User's actions

1. User logs into their local tool (acting as an artifact manager)

2. User selects a particular \*.tex ConTeXt *master* document and requests 
   the html version.

   NOTE: the choice of a particular \*.tex master document is really 
   choosing a dependency description file located in a 
   directory/archive/repository located in an artifact manager's 
   workspace. 

### Dependency description files

There are two principle types of dependency description files

1. ConTeXt document 

2. A Lua DSL

In both cases the lower-level interaction with the system will be via a 
Lua library. For the ConTeXt document this will be via a ConTeXt module 
written in Lua. 

The dependency description DSL will follow the ninja-build file 
principles. Conversely, the TaskManager interface will correspond closely 
to this dependency description DSL. 

## Assumed components

1. Artifact managers. One running in the user's file space will act as the 
   User's tool.

2. Federated messaging system (assumed and working in the background, that 
   is we do not want nor need to focus upon the lower-level messageing). 

3. Task managers (central resouce for a particular task request)

4. ConTeXt typesetter (can be called in parallel)

5. ConText compositor (central resource managing central update of one 
   particular document's \*.aux file data (one per task request)) 

6. GCC compiler (can be called in parallel)

7. GoLang compiler (can ONLY be called sequentially)

8. JoyLoL compiler/interpreter (should be callable in parallel)

9. pdf2htmlEX compiler (can ONLY be called sequentially? or can we stich 
   together runs of smaller pdfs? ) 

## System's actions

1. The User's tool acting as an artifact manager sends a message 
   requesting a Task Manager to manage the production of a html version of 
   the specified document. 

2. Any free TaskManagers repspond with offers to manage this new task.

3. The User's tool selects one TaskManager, who requests the dependency 
   document. 

4. The chosen TaskManager determines how to interpret the dependency 
   document. If the dependency document is a Lua DSL script, the 
   TaskManager directly runs the script. In this case, since the 
   dependency document is the ConTeXt document itself, the TaskManager 
   requests a free ConTeXt Typesetter. 

5. Any free ConTeXt Typesetters respond with offers to type set the 
   ConTeXt document (and hence interpret the dependency information). 

6. The TaskManger chooses a ConTeXt typesetter and informs the choosen 
   typesetter with details of the document to typeset.

7. The choosen ConTeXt typesetter asks the ArtefactManager (in this case 
   the User's Tool) for the ConTeXt document and begins typesetting it 
   with then context command in *draft* format. 

8. The ConTeXt typesetter will have to request any additional 
   sub-documents or ConTeXt modules from the collection of 
   ArtefactManagers via messages. 

9. The ConTeXt typesetter interacts with the TaskManager to describe all 
   dependencies described in the (sub)document(s). 

10. The TaskManager incrementally builds an expandable-truncated 
    topological sort of the dependency *graph* (which will be *acyclic*). 
    As each subtask completes the TaskManager issues requests to fulfil 
    subsequent tasks until the overall goal is reached. 

    **Here be dragons!!!!**
    
11. Once the overall goal is reached, the TaskManage informs the User's 
    Tool. 

## Problems

### How does the TaskManager connect goals with known dependency?

**Question:** How does the TaskManager connect overall goals, such as the 
html of a ConTeXt document, with the incrementally built dependency graph?

**Answer:** There are (potentially) *multiple* dependency descriptions for 
any overall problem. So the User's Tool is configured to supply an HTML 
dependency description (by interpreting Lua Dependency DSL as 
configuration) with a request to produce an HTML document *from* a given 
ConTeXt document. 

### Expanding partially specified requests

**Question:** Given the bare minimum of the user's request, which 
component interprets it? 

**Answer:** (In principle) the component "most" capable of understanding 
the dependency description file.

If the dependency description file is a ConTeXt document, then the ConTeXt 
typesetter will interpret the description and interact with the 
TaskManager. 

If the dependency description file is a Lua DSL script, then the 
ArtefactManager acting as the User's Tool will interpret the Lua 
dependency description directly, again interacting with the TaskManager. 
