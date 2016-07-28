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
