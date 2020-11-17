package git

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/neel1996/gitconvex-server/global"
	"go/types"
)

type AllCommitData struct {
	TotalCommits float64
	LatestCommit string
}

// AllCommits function returns the total number of commits from the repo and commit message of the most recent commit
func AllCommits(repo *git.Repository, commitChan chan AllCommitData) {
	logIter, itrErr := repo.Log(&git.LogOptions{})
	logger := global.Logger{}
	var commits []*object.Commit

	if itrErr != nil {
		logger.Log(fmt.Sprintf("Repo has no logs -> %s", itrErr.Error()), global.StatusError)
		commitChan <- AllCommitData{
			TotalCommits: 0,
			LatestCommit: "No Commits Available!",
		}
	} else {
		err := logIter.ForEach(func(commit *object.Commit) error {
			if commit != nil {
				commits = append(commits, commit)
				return nil
			} else {
				return types.Error{Msg: "Empty commit"}
			}
		})

		if err != nil {
			logger.Log(fmt.Sprintf("Unable to obtain commits for the repo"), global.StatusError)
			commitChan <- AllCommitData{
				TotalCommits: 0,
				LatestCommit: "No Commits Available!",
			}

		} else {
			logger.Log(fmt.Sprintf("Total commits in the repo -> %v", len(commits)), global.StatusInfo)

			commitChan <- AllCommitData{
				TotalCommits: float64(len(commits)),
				LatestCommit: commits[0].Message,
			}
		}
	}
	close(commitChan)
}
