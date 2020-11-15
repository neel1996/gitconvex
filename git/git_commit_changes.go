package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/neel1996/gitconvex-server/global"
)

func CommitChanges(repo *git.Repository, commitMessage string) string {
	logger := global.Logger{}
	w, wErr := repo.Worktree()

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
	}

	if wErr != nil {
		logger.Log(fmt.Sprintf("Error occurred while fetching repo worktree -> %s", wErr.Error()), global.StatusError)
		return "COMMIT_FAILED"
	} else {
		var commitOptions *git.CommitOptions
		if author == "" {
			logger.Log(fmt.Sprintf("No author name is available for the repo.\nSetting default signature for the commit"), global.StatusWarning)
			logger.Log("You can set the author details using git config commands --> https://support.atlassian.com/bitbucket-cloud/docs/configure-your-dvcs-username-for-commits/", global.StatusWarning)
			commitOptions = &git.CommitOptions{
				All: false,
				Author: &object.Signature{
					Name:  "gitconvex",
					Email: "help@gitconvex.com",
				},
			}
		} else {
			logger.Log(fmt.Sprintf("Commiting changes with author -> %s", author), global.StatusInfo)
			commitOptions = &git.CommitOptions{
				All: false,
			}
		}
		hash, err := w.Commit(commitMessage, commitOptions)
		if err != nil {
			logger.Log(fmt.Sprintf("Error occurred while committing changes -> %s", err.Error()), global.StatusError)
			return "COMMIT_FAILED"
		} else {
			logger.Log(fmt.Sprintf("Staged changes have been comitted - %s", hash.String()), global.StatusInfo)
			return "COMMIT_DONE"
		}
	}
}
