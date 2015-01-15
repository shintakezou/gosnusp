GoSNUSP
=======

[SNUSP](http://esolangs.org/wiki/SNUSP) is a “two-dimensional” esolanguage.

This code implement a Core, Modular and Bloated SNUSP interpreter in the [Go programming language](http://golang.org/).


Links and resources
-------------------

The repository contains also examples taken here and there. Other useful links or alike:

- [SNUSP on c2 Wiki](http://c2.com/cgi/wiki?SnuspLanguage)
- [Implementation of SNUSP interpreters on RosettaCode](http://rosettacode.org/wiki/Execute_SNUSP). The [C interpreter](http://rosettacode.org/wiki/RCSNUSP/C) should be mine; cfr. also [snuspi on sourceforge](http://sourceforge.net/projects/snuspi/) (if sourceforge works now…)
- [Examples on RosettaCode](http://rosettacode.org/wiki/Category:SNUSP)
- John Bauman's [esoteric language page](http://www.baumanfamily.com/john/esoteric.html) talks of SNUSP; he claims he has written the first complete full Bloated SNUSP interpreter — currently, following [the link](http://www.baumanfamily.com/john/snusp.py) raises a server error. It should be [this one](https://github.com/graue/esofiles/blob/master/snusp/impl/snusp.py)



Notes
------

- not fully tested (yet?)
- Modular SNUSP (which is the default) comes in two “flavours”; I call the second flavour *twisted* — which is now the default. If a modular SNUSP code does not work, try `-twist=false` flag… (examples in the [SNUSP page on esolangs.org](http://esolangs.org/wiki/SNUSP) are all “twisted”). If the code does not work anyway, you have found a bug — it would be nice if you let me know.
  - The difference between twisted and untwisted Modular SNUSP is in how the Enter (`@`) and Leave (`#`) command behave. Details in the code (sorry) and explicative example on the [blog](http://shintakezou.blogspot.it/2015/01/snusp-esolang-and-interpreter-in-go.html).
- current memory cell value should be given as program return code? It is not (yet?)
- I am a Go absolute beginner


Maybe-wanted features
---------------------

- Join for Bloated SNUSP?


