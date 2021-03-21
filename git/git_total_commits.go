package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
)

type AllCommitInterface interface {
	AllCommits(commitChan chan AllCommitData)
}

type AllCommitStruct struct {
	Repo *git2go.Repository
}

type AllCommitData struct {
	TotalCommits float64
	LatestCommit string
}

// AllCommits function returns the total number of commits from the repo and commit message of the most recent commit
func (t AllCommitStruct) AllCommits(commitChan chan AllCommitData) {
	repo := t.Repo
	logItr, itrErr := repo.Walk()
	var commits []git2go.Commit

	if itrErr != nil {
		logger.Log(fmt.Sprintf("Repo has no logs -> %s", itrErr.Error()), global.StatusError)
		commitChan <- AllCommitData{
			TotalCommits: 0,
			LatestCommit: "No Commits Available!",
		}
	} else {
		_ = logItr.PushHead()

		err := logItr.Iterate(func(commit *git2go.Commit) bool {
			if commit != nil {
				commits = append(commits, *commit)
				return true
			} else {
				return false
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

			if len(commits) > 0 {
				commitChan <- AllCommitData{
					TotalCommits: float64(len(commits)),
					LatestCommit: commits[0].Message(),
				}
			}
		}
	}
	close(commitChan)
}
