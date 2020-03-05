# ConTeXt Nursery configuration

## Principles

- All configuration is immutable and is changed by the administrator via a 
  collection of file system configuration files.

- Exactly ONE Nursery is configured as the Primary Nursery, with a 
  well-know address and port.

- All other Nurseries are configured with the Primary Nursery's well-know 
  address and port.

- Each Nursery will be assigned a "Server" Certificate Specific to them. 
  This "Server" certificate will be used as a "Server" certificate when 
  another Nursery is contacting this Nursery. This "Server" certificate 
  will be used as a "Client" certificate when this Nursery is connecting to 
  an other Nursery.

- All Nurseries will share the Certificate Authority's public certificate.

- The ConTeXtNurseryCA tool will generate "Server" certificates on demand.

- Each Nursery MAY be configured with a heartbeat interval.

## Questions


## Resources

- [The Laws of Reflection](https://blog.golang.org/laws-of-reflection)

- [Package reflect](https://golang.org/pkg/reflect/)

- [golang get a struct from an interface via 
  reflection](https://stackoverflow.com/questions/34272837/golang-get-a-struct-from-an-interface-via-reflection)

- [Type assertions and type 
  switches](https://yourbasic.org/golang/type-assertion-switch/)

- [Part 18: Interfaces - I](https://golangbot.com/interfaces-part-1/)

- [Part 19: Interfaces - II](https://golangbot.com/interfaces-part-2/)

- [ jinzhu / configor ](https://github.com/jinzhu/configor)

- [Parsing JSON files With Golang : Working with Unstructured 
  Data](https://tutorialedge.net/golang/parsing-json-with-golang/#working-with-unstructured-data)

- [How to handle configuration in Go : Just use standard go flags with 
  iniflags.](https://stackoverflow.com/a/25324191)
