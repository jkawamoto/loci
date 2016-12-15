# Loci: Testing remote CI scripts locally
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Build Status](https://travis-ci.org/jkawamoto/loci.svg?branch=master)](https://travis-ci.org/jkawamoto/loci)
[![Code Climate](https://codeclimate.com/github/jkawamoto/loci/badges/gpa.svg)](https://codeclimate.com/github/jkawamoto/loci)
[![Release](https://img.shields.io/badge/release-0.3.0-brightgreen.svg)](https://github.com/jkawamoto/loci/releases/tag/v0.3.0)

loci runs CI tests locally to make sure your commits will pass such tests
*before* pushing to a remote repository.

## Usage
Loci currently supports [Travis](https://travis-ci.org/)'s CI script,
and python and golang projects.

If your current directory has `.travis.yml`, run just `loci`.
Note that loci requires [docker](https://www.docker.com/).

Here is the help text of the `loci` command:
~~~
loci [global options] [script file]

   If script file is omitted, .travis.yml will be used.

GLOBAL OPTIONS:
   --name NAME, -n NAME  creating a container named NAME to run tests.
                         If name is given, container will not be deleted.
   --tag TAG, -t TAG     creating an image named TAG.
   --base TAG, -b TAG    use image TAG as the base image.
                         (default: "ubuntu:latest")
   --verbose             verbose mode, which prints Dockerfile and
                         entrypoint.sh.
   --apt-proxy URL       URL for a proxy server of apt repository.
                         [$APT_PROXY]
   --pypi-proxy URL      URL for a proxy server of pypi repository.
                         [$PYPI_PROXY]
   --http-proxy URL      URL for a http proxy server. [$HTTP_PROXY]
   --https-proxy URL     URL for a https proxy server. [$HTTPS_PROXY]
   --no-proxy LIST       Comma separated URL LIST for which proxies won't be
                         used. [$NO_PROXY]
   --help, -h            show help
   --version, -v         print the version
~~~

Note that `$XXX` means environment variable `$XXX` will be used
if the associated option value isn't given.


## Installation
```sh
$ go get github.com/jkawamoto/loci
```
or if you're [Homebrew](http://brew.sh/) user,

```sh
$ brew tap jkawamoto/loci
$ brew install loci
```

Compiled binaries are also available in
Github's [release page](https://github.com/jkawamoto/loci/releases).


# License
This software is released under the MIT License, see [LICENSE](LICENSE).
