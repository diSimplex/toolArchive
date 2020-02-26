# ConTeXt Nursery interface protocols

A ConTeXt Nursery will implement five interfaces:

1. [**Typesetting**](Typesetting.md) Responsible for initiating and 
   controlling the typesetting of a given ConTeXt document. It is also 
   responsible for assigning the "best" nursery for a given root ConTeXt 
   binary's use.

2. [**Artifact**](Artifact.md) Responsible for both listing the current 
   cached status of an artifact as well as downloading a given artifact.

3. [**Status**](Status.md) Responsible for querying *build* status:

    - build dependencies

    - build logs

4. [**Control**](Control.md) Responsible for managing the up, down, and 
   pause state of either a given Nursery or the whole federation.

5. [**Discovery**](Discovery.md) Responsible for communicating regular 
   load average, discovery, and heartbeat messages.

## Principles

- There will only be ONE Nursery on any given host

- Nurseries are long running services typcically started via systemd

- All interfaces will be RESTfull over HTTPS over the *same* port.

- The RESTful over HTTP messages will be formated as either HTML or JSON 
  depending upon the values listed in the request Accept header.

  - If the Accept header contains either of the 
    ["application/json"](https://www.iana.org/assignments/media-types/application/json) 
    or "text/json" (non-standard) values, then the response will be 
    formated as JSON (binary or text respectively).

  - If the Accpet header does not contain either of the "application/json" 
    or "text/json" values, then the response will be formated as HTML
    (text/html).

- The RESTful over HTTP message will be sent using TLS using both Client 
  and Server Certificates.

- Nurseries will use an artifact PULL model.

- When starting a typesetting, the root ConTeXt binary will start listening 
  on an TCP port for requests for files in its local directories. It will 
  inform the typesetting interface of the port that it is listening on, as 
  well as provide the url of its root document and then wait for requests 
  for documents.

## Questions

Do the Lua and GoLang network libraries support both Client and Server TLS 
certificates?

Do Lua and GoLang have encryption/decryption libraries?

Is the root ConTeXt binary (as used directly by the user) really a 
transient GoLang file server which places a request on a Nursery for its 
"root" document to be typeset?

How does the root ConTeXt binary find a Nursery? Does it listen for 
discovery messages and use the first Nersury it finds?

## Resources

### GoLang based Certificate Authority

- [go / go / refs/heads/master / . / src / crypto / tls / 
  generate_cert.go](https://go.googlesource.com/go/+/refs/heads/master/src/crypto/tls/generate_cert.go)

- [Creating a Certificate Authority + Signing Certificates in 
  Go](https://shaneutt.com/blog/golang-ca-and-signed-cert-go/)

- [TLS with selfsigned 
  certificate](https://stackoverflow.com/questions/22666163/tls-with-selfsigned-certificate)

- [square / certstrap : Tools to bootstrap CAs, certificate requests, and 
  signed certificates.](https://github.com/square/certstrap)

- [Create a PKI in 
  GoLang](https://fale.io/blog/2017/06/05/create-a-pki-in-golang/)

- [Using your own PKI for TLS in 
  Golang](http://www.hydrogen18.com/blog/your-own-pki-tls-golang.html)

- [Shyp / generate-tls-cert : Generating self signed 
  certificates](https://github.com/Shyp/generate-tls-cert)

- [[go-nuts] self-signed 
  certificate](https://grokbase.com/t/gg/golang-nuts/12b1y46sh1/go-nuts-self-signed-certificate)

- [[Golang] Build A Simple Web Service part.7 — Learn the SSL/TLS 
  connection](https://medium.com/a-layman/golang-build-a-simple-web-service-part-7-learn-the-ssl-tsl-connection-713b39f11eac)

### OpenSSL based Certificate Authority

- [How to Create Your Own SSL Certificate Authority for Local HTTPS 
  Development](https://deliciousbrains.com/ssl-certificate-authority-for-local-https-development/)

### Authentication by Double ended TLS / HTTPS

- [michaljemala/tls-client.go](https://gist.github.com/michaljemala/d6f4e01c4834bf47a9c4)

- [ denji / golang-tls ](https://github.com/denji/golang-tls)

- [GoLang TLS](https://golang.org/pkg/crypto/tls/)

### Authentication by Cookies

- [Swagger: Cookie 
  Authentication](https://swagger.io/docs/specification/authentication/cookie-authentication/)

- [Securing Cookie Based 
  Authentication](https://stackoverflow.com/questions/1283594/securing-cookie-based-authentication)

