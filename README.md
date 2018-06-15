# imis 
[![Build Status](https://travis-ci.org/foxbot/imis.svg?branch=master)](https://travis-ci.org/foxbot/imis) 
[![codecov](https://codecov.io/gh/foxbot/imis/branch/master/graph/badge.svg)](https://codecov.io/gh/foxbot/imis)

in memory image server

### what is this?

its redis but it only does 3 things

- create
- get
- list

the create deletes itself. you can change how long it lasts if you
like.

### why

its purpose built for a discord bot. solve spaghetti code with golang.