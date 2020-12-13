package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/utils"
	"runtime"
	"strings"
	"time"
)

// windowsCommit is used for committing changes using the git client if the platform is windows
// go-git commit fails with an access denied error for windows platform
func windowsCommit(repoPath string, msg string) string {
	args := []string{"commit", "-m", msg}
	cmd := utils.GetGitClient(repoPath, args)
	cmdStr, cmdErr := cmd.Output()

	if cmdErr != nil {
		logger.Log(fmt.Sprintf("Commit failed -> %s", cmdErr.Error()), global.StatusError)
		return global.CommitChangeError
	} else {
		logger.Log(fmt.Sprintf("Changes committed to the repo -> %s", cmdStr), global.StatusInfo)
		return global.CommitChangeSuccess
	}
}

// CommitChanges commits the staged changes to the repo
//
// The function falls back to the native git client for Windows platform due to an existing bug in the go-git library which
// blocks commits in windows platform
func CommitChanges(repo *git.Repository, commitMessage string) string {
	var formattedMessage = commitMessage
	logger := global.Logger{}
	w, wErr := repo.Worktree()

	if wErr != nil {
		logger.Log(fmt.Sprintf("Error occurred while fetching repo worktree -> %s", wErr.Error()), global.StatusError)
		return global.CommitChangeError
	} else {
		//Checking and splitting multi-line commit messages
		if strings.Contains(commitMessage, "||") {
			splitMessage := strings.Split(commitMessage, "||")
			formattedMessage = strings.Join(splitMessage, "\n")
		}

		// Checking OS platform for switching to git client for Windows systems
		platform := runtime.GOOS
		if platform == "windows" && w != nil {
			logger.Log(fmt.Sprintf("OS is %s -- Switching to native git client", platform), global.StatusWarning)
			return windowsCommit(w.Filesystem.Root(), formattedMessage)
		}

		// Checking if repo is a fresh repo with no branches
		// fallback function will be used to commit with git if no branches are present
		head, _ := repo.Head()
		if head == nil {
			logger.Log("Repo with no HEAD", global.StatusWarning)
			return windowsCommit(w.Filesystem.Root(), formattedMessage)
		}

		// Logic to check if the repo / global config has proper user information setup
		// Commit will be signed by default user if no user config is present
		globalConfig, gCfgErr := repo.ConfigScoped(config.GlobalScope)
		localConfig, lCfgErr := repo.ConfigScoped(config.LocalScope)
		var author string

		if gCfgErr == nil && lCfgErr == nil {
			fmt.Println(localConfig.User)
			fmt.Println(globalConfig.User)

			if globalConfig.User.Name != "" {
				author = globalConfig.User.Name
			} else if localConfig.User.Name != "" {
				author = localConfig.User.Name
			}
		} else {
			logger.Log(fmt.Sprintf("Unable to fetch repo config -> %v || %v", gCfgErr, lCfgErr), global.StatusError)
			return global.CommitChangeError
		}

		var commitOptions *git.CommitOptions
		var parentHash plumbing.Hash
		head, headErr := repo.Head()
		if headErr != nil {
			logger.Log(headErr.Error(), global.StatusError)
		} else {
			parentHash = head.Hash()
		}

		if author == "" {
			logger.Log(fmt.Sprintf("No author name is available for the repo.\nSetting default signature for the commit"), global.StatusWarning)
			logger.Log("You can set the author details using git config commands --> https://support.atlassian.com/bitbucket-cloud/docs/configure-your-dvcs-username-for-commits/", global.StatusWarning)
			commitOptions = &git.CommitOptions{
				All: false,
				Author: &object.Signature{
					Name:  "gitconvex",
					Email: "help@gitconvex.com",
					When:  time.Now(),
				},
				Parents: []plumbing.Hash{parentHash},
			}
		} else {
			logger.Log(fmt.Sprintf("Commiting changes with author -> %s, message -> %s", author, formattedMessage), global.StatusInfo)
			commitOptions = &git.CommitOptions{
				All:     false,
				Parents: []plumbing.Hash{parentHash},
			}
		}

		if formattedMessage == "" {
			return global.CommitChangeError
		}

		hash, err := w.Commit(formattedMessage, commitOptions)
		if err != nil {
			logger.Log(fmt.Sprintf("Error occurred while committing changes -> %s\n%v", err.Error(), err), global.StatusError)
			return global.CommitChangeError
		} else {
			logger.Log(fmt.Sprintf("Staged changes have been comitted - %s", hash.String()), global.StatusInfo)
			return global.CommitChangeSuccess
		}
	}
}
