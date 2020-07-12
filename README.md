
# Gitconvex

![gitconvex_logo](https://user-images.githubusercontent.com/47709856/87170859-8bfff080-c2ef-11ea-9140-b9e5db1c17d8.png)

## Web application for managing your git repositories

# Download options

> Use any of the below options to get gitconvex for your system

[![gitconvex npm package](https://badgen.net/badge/gitconvex/v1.1.0/green?icon=npm)](https://www.npmjs.com/package/@itassistors/gitconvex)
[![github release](https://badgen.net/badge/gitconvex/v1.1.0/blue?icon=github)](https://github.com/neel1996/gitconvex-package/releases)
[![docker image](https://badgen.net/badge/gitconvex/v1.1.0/cyan?icon=docker)](https://hub.docker.com/repository/docker/itassistors/gitconvex)

- Cloning repo from **github**

`git clone https://github.com/neel1996/gitconvex-package.git`

- Downloading package from **npm**

`npm i -g @itassistors/gitconvex`

This will install **gitconvex** as a global module and it can be started straight away from the command line with `gitconvex` command

```
$ gitconvex

INFO: Checking for config file
INFO: Config file is present
INFO: Reading from config file /usr/lib/node_modules/@itassistors/gitconvex/env_config.json
GitConvex API connected!

Checking data file availability...
INFO: Data file /usr/lib/node_modules/@itassistors/gitconvex/database/repo-datastore.json is present and it will be used as the active data file!

You can change this under the settings menu

Gitconvex is running on port 9001

    Open http://localhost:9001/ to access gitconvex
```

- If you are into **docker**, then there is also a docker image available for gitconvex 

`docker pull itassistors/gitconvex`

- Downloading the zip file from the tagged github [**release**](https://github.com/neel1996/gitconvex-package/releases)


# Platforms

|supported platforms|
|--|
|Linux :penguin:  |
|Mac OS  :apple: |
|Windows :black_square_button: |

## Requirements

| <b>[Node js](https://nodejs.org/en/)</b> | <b>Tested on v12.0+ |
|--|--|
| <b>[Git](https://git-scm.com/)</b> | <b>Tested on v2.20+</b> |

## Build status

| <b>React UI build status</b>  | ![Gitconvex UI pipeline](https://github.com/neel1996/gitconvex/workflows/Gitconvex%20UI%20pipeline/badge.svg) |
|--|--|
| <b>Node Server Build Status</b> | ![Gitconvex Server pipeline](https://github.com/neel1996/gitconvex-server/workflows/Gitconvex%20Server%20pipeline/badge.svg) |

# Detailed documentation

Refer the detailed [Documentation](DOCUMENTATION.md) for how to setup and use the platform

# License

See [LICENSE ](LICENSE) info for more

