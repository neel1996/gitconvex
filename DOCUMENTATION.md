# Introduction
GitConvex is a web app for managing your git repositories. It is supported by Linux, Mac OS and Windows. The [latest version](https://github.com/neel1996/gitconvex/releases) of GitConvex supports git features such as initializing and tracking a git repo, staging and unstaging changes, committing changes, pushing changes to the selected remote and so on.
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
1. [Git](https://git-scm.com/) (Version 2.20+)

If you wish to build the application from source then install,
   - [Node JS](https://nodejs.org/en/) (Version 12.0+)
   - [go](https://golang.org/dl/)

## Download Options
- **Option-1:** Directly clone the repo from GitHub
`git clone https://github.com/neel1996/gitconvex`

- **Option-2:** If you are into docker, then there is also a docker image available for gitconvex
`docker pull itassistors/gitconvex`

- **Option-3:** Downloading the zip file from the [releases](https://github.com/neel1996/gitconvex/releases).

```
## Extract the downloaded zip file and execute the commands

$ cd gitconvex

# for Mac & Linux
$ make build

# for Windows
$ ./make.bat build

```

## Setup
$ git clone https://github.com/neel1996/gitconvex.git
<br>
$ cd gitconvex

# For Mac & Linux
$ make build

# For Windows
$ ./make.bat build

## After build completion...
$ ./dist/gitconvex-server

## Features available
- Visualizing basic repo stats such as active branch, active remotes, number of files tracked etc
- Tracking modified files
- Creating new branches (provided there are no diverging changes)
- Initializing git inside a new repo and adding it to the platform tracker on the go
- Secure clone option with authentication
- File difference tracker with syntax highlighting
- File explorer with repository navigation features
- Code view capability from in-build repository explorer
- Commit log viewer
- Looking up desired commit logs using search feature
- Basic git operations such as staging, un-staging, committing and pushing to remote repo.

## How to use

### Left Pane Menu

- **Repositories** - To check tracked files changes, line-based changes with syntax highlighting and git operations (staging, un-staging, committing changes and pushing changed to remote)
- **Compare** - To compare the branches and commits for the selected repository.
- **Settings** - To check and edit internal data file, to remove a repo from GitConvex and to update the active port.
- **Help** - Includes documentation link and various options to report an issue or to submit feedback.

### Adding a new repo
- Use "+" at the bottom right corner to add a repo.
![add-a-repo](https://user-images.githubusercontent.com/65342122/88536126-db9d2680-d028-11ea-890f-c5fc11cd7cf0.png)

- Enter repo name and paste the repo path. If the folder is not a git repo then check the "*Check this if the folder is not a git repo*" checkbox to initialize git.
![repo-details](https://user-images.githubusercontent.com/65342122/101984594-80eb2b00-3ca8-11eb-97e5-5804ddfaed61.png)

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
The repo detail view also provides features for performing the following operations,
```
- Adding a new branch
- Pulling changes from remote
- Fetching changes from remote
- Adding a new remote repo
- Directory navigator to lookup all files and folders within the repo
- Code view for valid files from the file explorer view
- Loading commit logs dynamically
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

## Compare
- The **Compare** menu has two sections: **Branch compare** and **Commit compare**.

- **Branch compare** - Compares any two branches based on the selected base branch and compare branch. It displays difference between the selected branches. 

![branch-compare](https://user-images.githubusercontent.com/65342122/94800651-83a4e680-0402-11eb-9e34-1c7b53e66add.png)

- **Commit compare** - Compares any two commits based on the selected base commit and compare commit. It shows file difference between the selected commits.

![commit-compare](https://user-images.githubusercontent.com/65342122/94800717-9d462e00-0402-11eb-84aa-9890bfea1f78.png)

## Settings
- Settings in the left pane has three sections (Server data file, saved repos, Active GitConvex port number).
![settings](https://user-images.githubusercontent.com/65342122/87243003-a4871c80-c44f-11ea-9d1a-8350bdfb0da8.png)

- Server data file stores repo details such as the repo path, timestamp and the unique ID assigned to each repo. The data file must be an accessible JSON file with read / write permissions set to it. Also make sure you enter the full path for the file. E.g: /opt/my_data/data-file.json
- In the saved repos section, added repo(s) can be deleted permanently from GitConvex. 

>Note that, this will only remove the repo from GitConvex records and it will not perform an actual folder delete operation

- The port number can be updated to an available alternate port. Make sure that the port is not in use. The app needs to be restarted for the port change to take effect. 

## Help and Support

- Visit help section if you're facing an issue or need any help. If you have any queries or feedback, then discuss it in "Discord" or report an issue in GitHub.
![help-and-support](https://user-images.githubusercontent.com/65342122/101984821-d96ef800-3ca9-11eb-913e-f1cc062fae56.png)
