# Typesetting RESTful interface protocol

The typesetting interface is responsible for for initiating and controlling 
the typesetting of a given ConTeXt document. It is also responsible for 
assigning the "best" nursery for a given root ConTeXt binary's use.

It does this by doing:

1. Running previously installed shell scripts and/or commands
2. These commands can have command line arguments specified
3. These commands can have environment variables specified.
4. The output (both stdout and stderr) can be viewed in semi-real-time.
