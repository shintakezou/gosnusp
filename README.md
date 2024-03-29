GoSNUSP
=======

[SNUSP](http://esolangs.org/wiki/SNUSP) is a “two-dimensional” esolanguage.

This code implement a Core, Modular and Bloated SNUSP interpreter in the [Go programming language](http://golang.org/).
Updated to use the new (and sick, I daresay) `go.mod` thing.

[Github.io page](http://shintakezou.github.io/gosnusp) could contain a reference, in case you won't take a look at the given links (why not?)… Actually the page is just a test, and likely it will stay as it is forever (but never say never — though there's not tooo much to say about SNUSP you can't find elsewhere, and written better…)

Building
========

From inside the `src` directory, as simple as

    go build -o gosnusp
    
The `go.mod` states the Go version must be `1.20`, but any version with the `go.mod` system should work fine.

Links and resources
-------------------

The repository contains also examples taken here and there. Other useful links or alike:

- [SNUSP on c2 Wiki](http://c2.com/cgi/wiki?SnuspLanguage)
- [Implementation of SNUSP interpreters on RosettaCode](http://rosettacode.org/wiki/Execute_SNUSP). The [C interpreter](http://rosettacode.org/wiki/RCSNUSP/C) should be mine; cfr. also [snuspi on sourceforge](http://sourceforge.net/projects/snuspi/) (but moved [here on GitHub](https://github.com/shintakezou/snuspi))
- [Examples on RosettaCode](http://rosettacode.org/wiki/Category:SNUSP)
- John Bauman's [esoteric language page](http://www.baumanfamily.com/john/esoteric.html) talks of SNUSP; he claims he has written the first complete full Bloated SNUSP interpreter — currently, following [the link](http://www.baumanfamily.com/john/snusp.py) raises a server error. It should be [this one](https://github.com/graue/esofiles/blob/master/snusp/impl/snusp.py)



Notes
------

- not fully tested (yet?)
- Modular SNUSP (which is the default) comes in two “flavours”; I call the second flavour *twisted* — which is now the default. If a modular SNUSP code does not work, try `-twist=false` flag… (examples in the [SNUSP page on esolangs.org](http://esolangs.org/wiki/SNUSP) are all “twisted”). If the code does not work anyway, you have found a bug — it would be nice if you let me know.
  - The difference between twisted and untwisted Modular SNUSP is in how the Enter (`@`) and Leave (`#`) command behave. Details in the code (sorry) and explicative example on the [blog](http://shintakezou.blogspot.it/2015/01/snusp-esolang-and-interpreter-in-go.html).
- current memory cell value should be given as program return code? It is not so (yet?)
- I am a Go absolute beginner
  - and so far I don't like very much what they did with `go.mod`; maybe a simple project like this can be organized so that the new module system won't bother you? Anyway I dislike how it works, and how I can't immediately refer to local packages inside the very same project without going through that `go.mod` file.


Maybe-wanted features
---------------------

- Join for Bloated SNUSP?


