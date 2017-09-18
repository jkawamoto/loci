# Loci: Testing remote CI scripts locally
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Build Status](https://travis-ci.org/jkawamoto/loci.svg?branch=master)](https://travis-ci.org/jkawamoto/loci)
[![wercker status](https://app.wercker.com/status/25b462a013ed96bf51254862938e7659/s/master "wercker status")](https://app.wercker.com/project/byKey/25b462a013ed96bf51254862938e7659)
[![Release](https://img.shields.io/badge/release-0.5.3-brightgreen.svg)](https://github.com/jkawamoto/loci/releases/tag/v0.5.3)
[![Japanese](https://img.shields.io/badge/qiita-%E6%97%A5%E6%9C%AC%E8%AA%9E-brightgreen.svg)](http://qiita.com/jkawamoto/items/a409dd9cd6e63034aa28)

[![Loci Logo](https://jkawamoto.github.io/loci/img/image.png)](https://jkawamoto.github.io/loci/)

Loci runs CI tests locally to make sure your commits will pass such tests
*before* pushing to a remote repository.

Loci currently supports [Travis](https://travis-ci.org/)'s CI scripts
for [Python](https://www.python.org/) and [Go](https://golang.org/) projects.
Loci also requires [Docker](https://www.docker.com/) to run tests in a sandbox.

## Demo
[![asciicast](https://asciinema.org/a/126089.png)](https://asciinema.org/a/126089)

## Usage
If your current directory has `.travis.yml`, run just `loci` like

```shell
$ loci
```

If your `.travis.yml` specifies more than two runtime versions, Loci will run
those tests palatally. If you want to run tests on a selected one runtime
version, use `--select`/`-s` flag. For example, the following command runs tests
on only Python 3.6:

```shell
$ loci -s 3.6
```  

Here is the help text of the `loci` command:

~~~
loci [global options] [script file]

  If script file isn't given, .travis.yml will be used.

GLOBAL OPTIONS:
   --name NAME, -n NAME              base NAME of containers running tests.
                                     If not given, containers will be deleted.
   --select VERSION, -s VERSION      select specific runtime VERSION where tests
                                     running on.
   --tag TAG, -t TAG                 specify a TAG name of the docker image to
                                     be build.
   --max-processors value, -p value  max processors used to run tests.
   --log, -l                         store logging information to files.
   --base TAG, -b TAG                use image TAG as the base image.
                                     (default: "ubuntu:latest")
   --apt-proxy URL                   URL for a proxy server of apt repository.
                                     If environment variable APT_PROXY exists,
                                     that value will be used by default.
   --pypi-proxy URL                  URL for a proxy server of pypi repository.
                                     If environment variable PYPI_PROXY exists,
                                     that value will be used by default.
   --http-proxy URL                  URL for a http proxy server.
                                     If environment variable HTTP_PROXY exists,
                                     that value will be used by default.
   --https-proxy URL                 URL for a https proxy server.
                                     If environment variable HTTPS_PROXY exists,
                                     that value will be used by default.
   --no-proxy LIST                   Comma separated URL LIST for which proxies
                                     won't be used.
                                     If environment variable NO_PROXY exists,
                                     that value will be used by default.
   --no-build-cache                  Do not use cache when building the image.
   --no-color                        Omit to print color codes.
   --help, -h                        show help
   --version, -v                     print the version
~~~


## Installation
Loci works with [docker](https://www.docker.com/).
If your environment doesn't have docker, install it first.
The minimum required docker version is 1.12.0 (API version: 1.24).

If you're a [Homebrew](http://brew.sh/) or [Linuxbrew](http://linuxbrew.sh/)
user, you can install Loci by the following commands:

```shell
$ brew tap jkawamoto/loci
$ brew install loci
```

To build the newest version of Loci, use `go get` command:

```shell
$ go get github.com/jkawamoto/loci
```

Otherwise, compiled binaries are also available in
[Github](https://github.com/jkawamoto/loci/releases).


# License
This software is released under the MIT License, see [LICENSES](LICENSES.md).
