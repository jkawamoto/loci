## 0.5.0 (2017-06-24)
### Update
- Support parallel running tests.
- Support `--select`/`-s` flag to specify a runtime version tests will run on.
- Delete `--verbose` flag but add `--log` flag to store logging information to files.

### Fixed
- no-color mode.


## 0.4.5 (2017-06-19)
### Update
- Improve test preparation time.
- Use sub shell to execute each command.


## 0.4.4 (2017-06-09)
### Update
- To follow the update of docker client library.

### Fixed
- Provide sudo command in travis scripts.
- Parse quoted env.


## 0.4.3 (2017-02-10)
### Added
- Install sudo command to support old style Travis's configuration files.

### Fixed
- Parsing quoted environment variables.


## 0.4.2 (2017-02-01)
### Added
- Support no-color and no-build-cache options.

### Fixed
- before_install, install, before_script, and script attributes support both a single string and a list of strings.


## 0.4.1 (2017-01-31)
### Update
- Allow unbounded variables in Travis's scripts.


## 0.4.0 (2017-01-31)
### Added
- Build images based on installed runtime to reduce preparation time.

### Fixed
- Avoid installing apt packages when another version number is given.
- env attribute in .travis.yml allows both a string and a list of strings,
- Before install step allows define environment variables,
- Garbled characters in outputs of containers.


## 0.3.5 (2017-01-26)
### Fixed
- To lower repository names to create a docker image since docker tags allow only lower characters.


## 0.3.4 (2017-01-25)
### Fixed
- Stop testing when building a test image is failed.

### Update
- Use aria2 to prepare virtual environments for matrix evaluations for Python.


## 0.3.3 (2017-01-10)
### Fixed
- Declaring environment variables in startup scripts.


## 0.3.2 (2017-01-09)
### Fixed
- Running containers weren't stopped when loci is canceled,
- Use newer python version when an ambiguous version given,
- Fix garbled characters from outputs.


## 0.3.1 (2017-01-07)
### Added
- Support build matrices for go projects.
- Use repository names as a part of tag names of container images.

### Fixed
- Use correct archived source files.


## 0.3.0 (2016-12-14)
### Added
- Support matrix build for python projects.

### Fixed
- Set `GOPATH` and add $GOPATH/bin to `PATH`.


## 0.2.1 (2016-08-13)
### Fixed
- Dead lock when tarballing raises errors,
- List up functions gave directories to tarballing.


## 0.2.0 (2016-07-28)
### Added
Proxy supports:
- `--apt-proxy` flag sets a proxy URL for `apt-get` command,
- `--pypi-proxy` flag sets a proxy URL for `pip` command,
- `--http-proxy`, `--https-proxy`, and `--no-proxy` flags set http and https
  proxies.


## 0.1.5 (2016-07-22)
### Fixed
- Git directory will be added to a container so that test program can access
  repository information.


## 0.1.4 (2016-07-21)
### Added
- Support golang.


## 0.1.3 (2016-07-19)
### Added
- Support verbose mode,
- Use customized Dockerfile for each language,
- Support base flag to switch base image.

If choose an image previous run, some installations might be omitted.
It could reduce running time.


## 0.1.2 (2016-07-19)
### Added
- Support name and tag option to set container name and image tag.


## 0.1.1 (2016-07-18)
### Fixed
- Bugs about temporary directories, forget creating and deleting them.


## 0.1.0 (2016-07-18)
Initial release
