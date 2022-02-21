# The ConTeXtNursery configuration tool


## Using openssl to verify Certificates and Keys

### Verify CA certificate

```
    openssl verify -CAfile <<CA path>>-crt.pem <<CA path>>-crt.pem
```

### Verify a certificate chain

```
    openssl verify -CAfile <<CA path>>-crt.pem <<cert path>>-crt.pem
```

### Verify private keys

```
    openssl rsa -inform PEM -in <<private key path>>-key.pem -noout
    if [ $? == 0 ] ; then ; echo "VERIFIED" ; fi
```

### Dump a certificate

```
    openssl x509 -in <<cert path>> -text
```

### Dump a PKCS#12 file

```
    openssl pkcs12 -info -in <<pkcs12 path>>
```
