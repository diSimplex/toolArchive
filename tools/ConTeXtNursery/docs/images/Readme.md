# Images

The `svg` images which have matching `odp` files ([LibreOffice 
Impress](https://www.libreoffice.org/discover/impress/)) have been created 
by exporting the `odp` file as a `pdf` file and then using the 
[`pdftocairo` command line 
tool](http://manpages.ubuntu.com/manpages/precise/man1/pdftocairo.1.html) 
(a [poppler](https://poppler.freedesktop.org/) utility) to convert the 
`pdf` file into an `svg` file. For example:

```
    pdftocairo -svg ConTeXtNursery.pdf ConTeXtNursery.svg
```

