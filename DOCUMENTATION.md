
# Introduction

 
GitConvex is a web app for managing your git repositories. It is supported by Linux, Mac OS and Windows. The [latest version](https://github.com/neel1996/gitconvex-package/releases) of GitConvex supports git features such as initializing and tracking a git repo, staging and unstaging changes, committing changes, pushing changes to selected remote and so on.
The main goal of this platform is to act as a web based alternative for Github desktop, but we are still in the starting stages, so we are not completely there yet (baby steps, right?)
## Requirements
1. Node JS (Version 12.0+)

2. Git (Version 2.20+)

  

## Download Options

 
- **Option-1:** Directly clone the repo from github

`git clone https://github.com/neel1996/gitconvex-package`

  

- **Option-2** Downloading the zip file from the [releases](https://github.com/neel1996/gitconvex-package/releases).

  

- **Option-3** GitConvex is also available on `npm`. Install the package globally to run it directly from the command line,

  

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

  
  

## Setup

If either download **Option-1** or **Option-2** is followed, then the following steps need to be followed to setup GitConvex

  

1. For installing all dependencies,

`` $npm install ``

  

2. To start a server, either use normal node command

`` $node server.js ``

  

or use `pm2` by downloading it from npm - `npm i -g pm2` and start the module by executing the following command,

`pm2 start ecosystem.config.js`

  

## Features available

- Visualizing basic repo stats such as active branch, active remotes, number of files tracked etc

- Tracking modified files

- Creating new branches (provided there are no diverging changes)

- Initializing git inside a new repo and adding it to the platform tracker on the go

- File difference tracker with syntax highlighting for the [supported languages](language).

- Commit log viewer

- Basic git operations such as staging, un-staging, committing and pushing to remote repo.

  

## How to use

  

### Left Pane Menu

  

- **Repositories** - To check tracked files changes, line based changes with syntax highlighting and git operations (staging, unstaging, commiting changes and pushing changed to remote)

- **Settings** - To check and edit internal data file, to remove a repo from gitconvex and to update the active port

- **Help** - Includes documentation link and various options to report an issue or to submit feedback

  

### Adding a new repo

- Use "+" at the bottom right corner to add a repo.

![Picture1](https://user-images.githubusercontent.com/65342122/87232632-0eff7480-c3de-11ea-8a9f-f0a6cf9cd6ee.png)

<br>

- Enter repo name and paste the repo path. If the folder is not a git repo then check the "*Check this if the folder is not a git repo*" checkbox to initialize git.

![Picture2](https://user-images.githubusercontent.com/65342122/87232637-16268280-c3de-11ea-9f9d-708c5a3eb668.png)

  

- The newly added repo will be displayed as a card in the dashboard

![Picture3](https://user-images.githubusercontent.com/65342122/87232640-1aeb3680-c3de-11ea-8751-e47e5f64c8a1.png)

### Repo Details

- Click on the repo card to get the following details about the repo

```

- The list of branches

- Commit logs

- Latest commit

- Active branch and available local branches

- Remote repo URL and host

- Files and folders tracked by git

```

The repo detail view also provides features for performing the following git operations,

```

- Adding a new branch

- Pulling changes from remote

- Fetching changes from remote

- Adding a new remote repo

```

  

![Picture4](https://user-images.githubusercontent.com/65342122/87232642-1f175400-c3de-11ea-8ead-80cd5ab4c37c.png)

![enter image description here](https://user-images.githubusercontent.com/65342122/87232671-71587500-c3de-11ea-8ff5-25fb95b08553.png)

  

#### Add new branch

  

![Picture5](https://user-images.githubusercontent.com/65342122/87232643-23437180-c3de-11ea-9d9e-7e3c3789c32e.png)

Note: In a newly initialized git repo, the newly added branch will be considered by git only after a initial commit

  

- Next, left pane consists of repositories, settings and help. In "repositories" choose a saved repo. This module has three sections : File view, Git difference and Git operation. The header portion shows chosen repo name, active branch, number of tracked files and commits.

- File view shows the new/Untracked/Modified/deleted files.

![Picture6](https://user-images.githubusercontent.com/65342122/87232644-29d1e900-c3de-11ea-9adc-03fb4e690882.png)

- Git operations module consists of three basic git operations( Stage all changes, commit changes, push to remote ). Below this file status is display and the files can be staged with add action or click on "stage all changes" to stage all the untracked files from the chosen repo.

  

![Picture7](https://user-images.githubusercontent.com/65342122/87232645-30f8f700-c3de-11ea-8ddb-52f4d5ec7140.png)

- After add operation all staged files can be removed immediately if you wish to, else ignore it. Reload the page to see the difference in "File View".

![Picture8](https://user-images.githubusercontent.com/65342122/87232658-51c14c80-c3de-11ea-95e1-b9bbeeac82bb.png)

- All staged files can be committed using "Commit changes" module. After commit changes reload the page to see active branch name, number of tracked files and number of commits.

  

![Picture9](https://user-images.githubusercontent.com/65342122/87232659-56860080-c3de-11ea-9bc4-a19ad727b101.png)

- After **commit changes** operation,

  

![Picture10](https://user-images.githubusercontent.com/65342122/87232661-5c7be180-c3de-11ea-82b7-104792c2e3ec.png)

  

- "Push to remote" module pushes all commits to the remote host. If there is no remote in the chosen repo then ignore this module. After the push to remote is completed, check the remote host to see the pushed folder/file.

with Remote(Git Hub)

![Picture11](https://user-images.githubusercontent.com/65342122/87232662-61409580-c3de-11ea-8ad7-61c3871f0a4d.png)

If you try to "push to remote" without having a remote host, "push failed" error message will be displayed.

  

![Picture12](https://user-images.githubusercontent.com/65342122/87232666-6867a380-c3de-11ea-9903-5ea12200e994.png)
 

In "Git Difference" click on the modified file(s) to see the difference.

![Picture13](https://user-images.githubusercontent.com/65342122/87232669-6d2c5780-c3de-11ea-9739-f8181e4d0901.png)

- "Stage all changes" and "commit changes" for applying a change.

  

- Settings in the left pane has three sections( Server data file, saved repos, Active GitConvex port number ).

  

![Picture15](https://user-images.githubusercontent.com/65342122/87232673-75849280-c3de-11ea-9d01-fe3479282561.png)

  

- Server data file stores repo details. The data file can be updated. The data file must be an accessible JSON file with read / write permissions set to it. Also make sure you enter the full path for the file. E.g: /opt/my_data/data-file.json

- In the saved repos section, added repo(s) can be deleted permanently from GitConvex.

- The port number can be updated but make sure to restart the app and to change the port in the URL after updating it.

- Visit help section if you're facing an issue or need any help. If you have any queries discuss it in "Discord" or report an issue in GitHub.

![Picture16](https://user-images.githubusercontent.com/65342122/87232674-79b0b000-c3de-11ea-904a-c6a02c4d15f5.png)
