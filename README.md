# Gitconvex

### Web application for managing your git repositories

<p align="center">
    <img src="https://user-images.githubusercontent.com/47709856/87170859-8bfff080-c2ef-11ea-9140-b9e5db1c17d8.png" width="280">
    <p align="center">
        <a href="https://www.producthunt.com/posts/gitconvex-2?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-gitconvex-2" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=241240&theme=dark" alt="gitconvex - A web based UI client for managing git repositories | Product Hunt Embed" style="width: 250px; height: 54px;" width="250px" height="54px" /></a>
    </p>
</p>

## Build Status

![Gitconvex UI pipeline](https://img.shields.io/github/workflow/status/neel1996/gitconvex-ui/Gitconvex%20UI%20pipeline/master?label=gitconvex%20ui%20build&logo=github&style=for-the-badge)
![Gitconvex Server pipeline](https://img.shields.io/github/workflow/status/neel1996/gitconvex-server/Gitconvex%20Server%20Pipeline/main?label=gitconvex%20server%20build&logo=github&style=for-the-badge)

# Download options

> Use any of the below options to get gitconvex for your system

[![github release](https://img.shields.io/static/v1?label=gitconvex&message=v2.0.1&color=green&style=for-the-badge&logo=github)](https://github.com/neel1996/gitconvex-package/releases)
[![docker image](https://img.shields.io/static/v1?label=gitconvex&message=v2.0.1&color=blue&style=for-the-badge&logo=docker)](https://hub.docker.com/repository/docker/itassistors/gitconvex)


- **Option - 1** Cloning the repo from **github**

To build gitconvex from scratch, [libgit2](https://github.com/libgit2/libgit2) is required. The [LIBGIT_NOTES](LIBGIT_NOTES.md) includes the guidelines to download and setup libgit for different platforms

**To be Noted:** The `master` branch contains the latest and the stable build of the project. For a reliable experience, always clone the repo from the master branch.

``` shell

git clone https://github.com/neel1996/gitconvex.git
cd gitconvex

# for Mac & Linux
make build

# for Windows
./make.bat build

## After build completion...
./dist/gitconvex-server

2020/11/14 22:57:47 INFO: Starting Gitconvex server modules
2020/11/14 22:57:47 INFO: Using available env config file
2020/11/14 22:57:47 INFO: Gitconvex started on  http://localhost:9001

```

- **Option - 2** If you are into **docker**, then there is a docker image available for gitconvex 

``` shell
docker pull itassistors/gitconvex
```

**Note:** Make sure you mount the host volume to the container, to access the git repos from the host system. If you have git repos stored within your containers, then this is not required

- **Option - 3** Downloading the zip file from the tagged github [**release**](https://github.com/neel1996/gitconvex/releases)

``` shell
## Extract the downloaded zip file and execute the commands
cd gitconvex

# for Mac & Linux
make build

# for Windows
./make.bat build

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
| <b>[Git](https://git-scm.com/)</b> | <b>Required for handling some git operations with the target repo</b> |
| <b>[Go](https://golang.org/)</b> | <b>For building the backend from the source</b> |
| <b>[Node JS](https://nodejs.org/en/)</b> | <b>For building the React UI bundle from scratch</b> |
    
> **Important note for windows users**

If you are a windows user, then make sure that `git` is accessible from the command line without using git bash.

- Open command prompt or powershell and enter `git --version` and press enter. If it displays the following output, then it is fine

``` cmd
C:\> git --version

git version 2.28.0.windows.1
```

If this output is not displayed and if the command throws the following error, then it shows that `git` is not added to the 'path' environment variable and it will not be accessible directly from the command line,

``` cmd
C:\> git --version

'git' is not recognized as an internal or external command,
operable program or batch file.
```

This can be fixed by adding `git` to the PATH environment variable in windows. The steps for setting this up is available [here](https://stackoverflow.com/questions/26620312/git-installing-git-in-path-with-github-client-for-windows#answer-53706956:~:text=comment-,27,Here%20is%20the%20magic)

> **🍎 Important note for MacOS users**

The pre-built bundle for MacOS is not a verified or signed bundle, so gatekeeper could warn you or even prevent you from using gitconvex on your Mac devices. If this is the case, then I recommend building the application from scratch using the `Makefile` included in the repo. Follow **[Option - 1](#download-options)** mentioned above to build the application from scratch.

**Reason** - Enrolling in the apple developer program for making the application a verified one will cost me 100 USD annually, so I will do it once the project gets enough reach 

# Detailed documentation

Refer the detailed [Documentation](DOCUMENTATION.md) for how to setup and use the platform


# Contributions 
Contributions are always welcome for the project. Please do refer the [Contribution Guidelines]( CONTRIBUTING.md) for details
# Help and Feedback

For reporting issues or for requesting a new feature, use the following channels

[**Discord Channel** ](https://discord.gg/PSd2Cq9)

[**Website**](https://gitconvex.com/)

[**Github Issue Reporting**](https://github.com/neel1996/gitconvex/issues)

![open issues](https://img.shields.io/github/issues/neel1996/gitconvex?color=orange&style=for-the-badge)

# License

[![License](https://img.shields.io/static/v1?label=LICENSE&message=Apache-2.0&color=yellow&style=for-the-badge)](LICENSE)


