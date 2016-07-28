# Loci: Testing remote CI scripts locally
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Build Status](https://travis-ci.org/jkawamoto/loci.svg?branch=master)](https://travis-ci.org/jkawamoto/loci)
[![Code Climate](https://codeclimate.com/github/jkawamoto/loci/badges/gpa.svg)](https://codeclimate.com/github/jkawamoto/loci)
[![Release](https://img.shields.io/badge/release-0.1.5-lightgrey.svg)](https://github.com/jkawamoto/loci/releases/tag/v0.1.5)

loci runs CI tests locally to help your commits pass such tests
before pushing remote repository.


## Usage
~~~
loci [global options] [script file]

   If script file is omitted, .travis.yml will be used.

GLOBAL OPTIONS:
   --name NAME, -n NAME  creating a container named NAME to run tests.
                         If name is given, continer will not be deleted.
   --tag TAG, -t TAG     creating an image named TAG.
   --base TAG, -b TAG    use image TAG as the base image.
   --verbose             verbose mode, which prints Dockerfile and
                         entrypoint.sh.
   --help, -h            show help
   --version, -v         print the version
~~~
