# Loci: Testing remote CI scripts locally
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![Build Status](https://travis-ci.org/jkawamoto/loci.svg?branch=master)](https://travis-ci.org/jkawamoto/loci)
[![wercker status](https://app.wercker.com/status/25b462a013ed96bf51254862938e7659/s/master "wercker status")](https://app.wercker.com/project/byKey/25b462a013ed96bf51254862938e7659)
[![Release](https://img.shields.io/badge/release-0.4.3-brightgreen.svg)](https://github.com/jkawamoto/loci/releases/tag/v0.4.3)
[![Japanese](https://img.shields.io/badge/qiita-%E6%97%A5%E6%9C%AC%E8%AA%9E-brightgreen.svg)](http://qiita.com/jkawamoto/items/a409dd9cd6e63034aa28)

[![Loci Logo](https://jkawamoto.github.io/loci/img/image.png)](https://jkawamoto.github.io/loci/)

Loci runs CI tests locally to make sure your commits will pass such tests
*before* pushing to a remote repository.

Loci currently supports [Travis](https://travis-ci.org/)'s CI scripts
for [Python](https://www.python.org/) and [Go](https://golang.org/) projects.
Loci also requires [Docker](https://www.docker.com/) to run tests in a sandbox.


## Usage
If your current directory has `.travis.yml`, run just `loci`.

Here is the help text of the `loci` command:

~~~
loci [global options] [script file]

  If script file isn't given, .travis.yml will be used.

GLOBAL OPTIONS:
   --name NAME, -n NAME  creating a container named NAME to run tests,
                         and that container will not be deleted.
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
   --no-proxy LIST       Comma separated URL LIST for which proxies won't
                         be used. [$NO_PROXY]
   --no-build-cache      Do not use cache when building the image.
   --no-color            Omit to print color codes.
   --help, -h            show help
   --version, -v         print the version
~~~

Note that `$XYZ` means that environment variable `$XYZ` will be used
if the associated option value isn't given.


## Installation
Loci works with [docker](https://www.docker.com/).
If your environment doesn't have docker, install it first.

To build the newest version of Loci, use `go get` command:

```shell
$ go get github.com/jkawamoto/loci
```

If you're a [Homebrew](http://brew.sh/) user, you can install Loci by
the following commands:

```shell
$ brew tap jkawamoto/loci
$ brew install loci
```

Otherwise, compiled binaries are also available in
[Github](https://github.com/jkawamoto/loci/releases).


# License
This software is released under the MIT License, see [LICENSES](LICENSES.md).
