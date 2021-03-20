# Gitconvex

### Web application for managing your git repositories

<p align="center">
    <img src="https://user-images.githubusercontent.com/47709856/87170859-8bfff080-c2ef-11ea-9140-b9e5db1c17d8.png" width="280">
</p>

## Build Status

![Gitconvex UI pipeline](https://img.shields.io/github/workflow/status/neel1996/gitconvex-ui/Gitconvex%20UI%20pipeline/master?label=gitconvex%20ui%20build&logo=github&style=for-the-badge)
![Gitconvex Server pipeline](https://img.shields.io/github/workflow/status/neel1996/gitconvex-server/Gitconvex%20Server%20Pipeline/main?label=gitconvex%20server%20build&logo=github&style=for-the-badge)

# Download options

> Use any of the below options to get gitconvex for your system

[![github release](https://img.shields.io/static/v1?label=gitconvex&message=v2.1.0&color=green&style=for-the-badge&logo=github)](https://github.com/neel1996/gitconvex-package/releases)
[![gitconvex-choco](https://img.shields.io/chocolatey/v/gitconvex?logo=C&style=for-the-badge)](https://chocolatey.org/packages/gitconvex)
[![docker image](https://img.shields.io/static/v1?label=gitconvex&message=v2.1.0&color=blue&style=for-the-badge&logo=docker)](https://hub.docker.com/repository/docker/itassistors/gitconvex)

- [Windows](#windows)
- [Linux](#linux)
- [MacOS](#macos)
- [Docker](#docker)

## Windows

Gitconvex is available on [chocolatey](https://chocolatey.org/). Install [chocolatey](https://chocolatey.org/install) and run the following command to install gitconvex

```shell
> choco install gitconvex --version=2.1.0
```

If you do not prefer `chocolatey` then the [binaries for windows](https://github.com/neel1996/gitconvex/releases/download/2.1.0/gitconvex-v2.1.0-windows.zip) can be downloaded from the [releases](https://github.com/neel1996/gitconvex/releases/tag/2.1.0)

## Linux

Download the tarball for Linux from [releases](https://github.com/neel1996/gitconvex/releases/tag/2.1.0). Extract the tar and run the following script to install the required dependencies

``` shell
> ./lib/linux_x86_64/install_deps
```

After setting up the dependency libs, just run `gitconvex` from the terminal to start gitconvex

Gitconvex is also available as a [homebrew](https://brew.sh/) tap which can be easily installed using the following command (provided [brew](https://brew.sh/) is already installed on the linux machine)

``` shell
> brew install itassistors/taps/gitconvex
```

## MacOS

Gitconvex is readily available as a [Homebrew](https://brew.sh/) Tap on github. Run the following command to download & install gitconvex on the go

``` shell
> brew install itassistors/taps/gitconvex
```

## Docker 

If you are into **docker**, then there is a docker image available for gitconvex

``` shell
docker pull itassistors/gitconvex
```

**Note:** Make sure you mount the host volume to the container to access the git repos from the host system. If you have git repos stored within your containers then it is not required

## Build from scratch

To build gitconvex from scratch, [libgit2](https://github.com/libgit2/libgit2) is a requirement. The [LIBGIT_NOTES](LIBGIT_NOTES.md) file includes all the guidelines to download and setup libgit for different platforms

**To be Noted:** The `master` branch contains the latest and the stable build of the project. For a reliable experience, always clone the repo from the master branch.

``` shell

git clone https://github.com/neel1996/gitconvex.git
cd gitconvex

# for Mac & Linux
make build

# for Windows
./make.bat build

# After build completion...
./dist/gitconvex-server

2020/11/14 22:57:47 INFO: Starting Gitconvex server modules
2020/11/14 22:57:47 INFO: Using available env config file
2020/11/14 22:57:47 INFO: Gitconvex started on  http://localhost:9001

```

# Platforms

|supported platforms|
|---|
|Linux :penguin:  |
|Mac OS  :apple: |
|Windows :black_square_button: |

## Requirements

| Software | Purpose |
| --- | --- |
| <b>[Git](https://git-scm.com/)</b> | <b>For cloning gitconvex and to build the application from scratch</b> |
| <b>[Go](https://golang.org/)</b> | <b>For building the backend from the source</b> |
| <b>[Node JS](https://nodejs.org/en/)</b> | <b>For building the React UI bundle from scratch</b> |
    
# Detailed documentation

Refer the detailed [Documentation](DOCUMENTATION.md) for how to setup and use the platform

# Contributions 
Contributions are always welcome for the project. Please do refer the [Contribution Guidelines](CONTRIBUTING.md) for details

# Help and Feedback

For reporting issues or for requesting a new feature, use the following channels

[**Discord Channel**](https://discord.gg/PSd2Cq9)

[**Website**](https://gitconvex.com/)

[**Github Issue Reporting**](https://github.com/neel1996/gitconvex/issues)

![open issues](https://img.shields.io/github/issues/neel1996/gitconvex?color=orange&style=for-the-badge)

# License

[![License](https://img.shields.io/static/v1?label=LICENSE&message=Apache-2.0&color=yellow&style=for-the-badge)](LICENSE)
