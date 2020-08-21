# Introduction
GitConvex is a web app for managing your git repositories. It is supported by Linux, Mac OS and Windows. The [latest version](https://github.com/neel1996/gitconvex-package/releases) of GitConvex supports git features such as initializing and tracking a git repo, staging and unstaging changes, committing changes, pushing changes to the selected remote and so on.
The main goal of this platform is to act as a web-based alternative for Github desktop, but we are still in the starting stages, so we are not completely there yet (baby steps, right?)
## Table of Contents
[Requirements](#requirements)<br>
[Download Options](#download-options)<br>
[Setup](#setup)<br>
[Features available](#features-available)<br>
[How to use](#how-to-use)
- [Adding a new repo](#adding-a-new-repo)
- [Repo Details](#repo-details)
- [Add a new branch](#add-new-branch)
 
## Requirements
1. [Node JS](https://nodejs.org/en/) (Version 12.0+)
2. [Git](https://git-scm.com/) (Version 2.20+)
## Download Options
- **Option-1:** Directly clone the repo from GitHub
`git clone https://github.com/neel1996/gitconvex-package`

- **Option-2:** Downloading the zip file from the [releases](https://github.com/neel1996/gitconvex-package/releases).

- **Option-3:** GitConvex is also available on `npm`. Install the package globally to run it directly from the command line
`npm i -g @itassistors/gitconvex`

This will install **GitConvex** as a global module and it can be started straight away from the command line with `gitconvex` command
```
$ gitconvex
  
INFO: Checking for config file
INFO: Config file is present
INFO: Reading from config file /usr/lib/node_modules/@itassistors/gitconvex/env_config.json
GitConvex API connected!
  
Checking data file availability...
INFO: Data file /usr/lib/node_modules/@itassistors/gitconvex/database/repo-datastore.json is present and it will be used as the active data file!
You can change this under the settings menu
GitConvex is running on port 9001
Open http://localhost:9001/ to access GitConvex
```
## Setup

If either download **Option-1** or **Option-2** is opted, then the following steps need to be followed to setup GitConvex
1. For installing all dependencies,
`` $npm install ``

2. To start the server, either use normal node command
`` $node server.js ``
or use `pm2` by downloading it from npm - `npm i -g pm2` and start the module by executing the following command,
`pm2 start ecosystem.config.js`

## Features available
- Visualizing basic repo stats such as active branch, active remotes, number of files tracked etc
- Tracking modified files
- Creating new branches (provided there are no diverging changes)
- Initializing git inside a new repo and adding it to the platform tracker on the go
- File difference tracker with syntax highlighting for the [supported languages](LANGUAGES.md).
- Commit log viewer
- Basic git operations such as staging, un-staging, committing and pushing to remote repo.

## How to use

### Left Pane Menu

- **Repositories** - To check tracked files changes, line-based changes with syntax highlighting and git operations (staging, un-staging, committing changes and pushing changed to remote)
- **Settings** - To check and edit internal data file, to remove a repo from GitConvex and to update the active port.
- **Help** - Includes documentation link and various options to report an issue or to submit feedback.

### Adding a new repo
- Use "+" at the bottom right corner to add a repo.
![add-a-repo](https://user-images.githubusercontent.com/65342122/88536126-db9d2680-d028-11ea-890f-c5fc11cd7cf0.png)

- Enter repo name and paste the repo path. If the folder is not a git repo then check the "*Check this if the folder is not a git repo*" checkbox to initialize git.
![repo-details](https://user-images.githubusercontent.com/65342122/88486721-70f2d900-cf9d-11ea-802d-00b0b2aaf638.png)

- The newly added repo will be displayed as a card in the dashboard
![repo-card](https://user-images.githubusercontent.com/65342122/89167157-ab113b80-d598-11ea-8985-2469e7ad261e.png)

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
- Directory navigator
- Loading commit logs on demand
```

![repo-card-details](https://user-images.githubusercontent.com/65342122/89164057-30deb800-d594-11ea-94d6-d3a330260044.png)

### Commit logs

- With commit log searchbar, any commit log can be looked up using its commit message or commit hash or author name who created that commit.

![commit-logs](https://user-images.githubusercontent.com/65342122/90782955-1723cb80-e31d-11ea-9c42-d1d5a6306e6f.png)

#### Add new branch
![new-branch](https://user-images.githubusercontent.com/65342122/87232643-23437180-c3de-11ea-9d9e-7e3c3789c32e.png)

Note: In a newly initialized git repo, the newly added branch will be considered by git only after an initial commit

## Repository Menu
- The **Repository** menu has three sections: **File view, Git difference and Git operations**, which will be displayed based on the selected repo. 

### File View
- File view shows the New/Untracked/Modified/Deleted files.
![file-view](https://user-images.githubusercontent.com/65342122/87232644-29d1e900-c3de-11ea-9adc-03fb4e690882.png)

### Git Operations
- Git operations module lets you handle three basic git operations (**Stage all changes, commit changes, push to remote**). Below this option, the files will be displayed and the files can be staged individually using the "Add" button or as a whole using the "stage all changes".
![git-operations](https://user-images.githubusercontent.com/65342122/87232645-30f8f700-c3de-11ea-8ddb-52f4d5ec7140.png)

- The staged files can be removed individually or it can be removed all at once.
![staged-files](https://user-images.githubusercontent.com/65342122/89167388-ff1c2000-d598-11ea-8150-fc96b6aa41b7.png)

- The staged changes can be committed using the **Commit changes** option. This will display a pop-up with all the staged files and it requires a commit message to successfully commit the changes. The commit messages can either be a single line message or a multi-line message 
![commit-changes](https://user-images.githubusercontent.com/65342122/87232659-56860080-c3de-11ea-9bc4-a19ad727b101.png)

- **Push to remote** option pushes all commits to the selected remote host. The pop-up displayed will display the commits which are in queue to be pushed to the remote repository

![push-operation-with-remote](https://user-images.githubusercontent.com/65342122/89562600-d796b300-d837-11ea-969c-9abb93f24c1c.png)

This section will let you know if the selected remote is not valid or if the push operation fails
![push-opeartion-without-remote](https://user-images.githubusercontent.com/65342122/89562476-a1f1ca00-d837-11ea-9f3a-bb566aeb327e.png)

### Git Difference

In "Git Difference" click on the modified file to see the difference. The platform has syntax highlighting available for a limited set of [languages](LANGUAGES.md)

![git-difference](https://user-images.githubusercontent.com/65342122/89559704-ab793300-d833-11ea-88a8-28afea4b461b.png)

## Settings
- Settings in the left pane has three sections (Server data file, saved repos, Active GitConvex port number).
![settings](https://user-images.githubusercontent.com/65342122/87243003-a4871c80-c44f-11ea-9d1a-8350bdfb0da8.png)

- Server data file stores repo details such as the repo path, timestamp and the unique ID assigned to each repo. The data file must be an accessible JSON file with read / write permissions set to it. Also make sure you enter the full path for the file. E.g: /opt/my_data/data-file.json
- In the saved repos section, added repo(s) can be deleted permanently from GitConvex. 

>Note that, this will only remove the repo from GitConvex records and it will not perform an actual folder delete operation

- The port number can be updated to an available alternate port. Make sure that the port is not in use. The app needs to be restarted for the port change to take effect. 

## Help and Support

- Visit help section if you're facing an issue or need any help. If you have any queries or feedback, then discuss it in "Discord" or report an issue in GitHub.
![help-and-support](https://user-images.githubusercontent.com/65342122/87242999-8f11f280-c44f-11ea-9a81-f6cde7b4b419.png)
