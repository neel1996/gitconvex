package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/neel1996/gitconvex-server/global"
	"go/types"
	"strings"
	"time"
)

// commitModel function generates a pipe separated string with the required commit details
// The resultant string will be sent to the client
func commitModel(commit object.Commit) string {
	commitHash := commit.Hash.String()
	commitAuthor := strings.Split(commit.Author.String(), " ")[0]
	commitMessage := strings.Split(commit.Message, "\n")[0]
	commitDate := ""

	for _, cString := range strings.Split(commit.String(), "\n") {
		if strings.Contains(cString, "Date:") {
			str := strings.Split(cString, "Date:")[1]
			tempDate := strings.TrimSpace(str)

			if strings.Contains(tempDate, "+") {
				tempDate = strings.TrimSpace(strings.Split(tempDate, "+")[0])
			} else if strings.Contains(tempDate, "-") {
				tempDate = strings.TrimSpace(strings.Split(tempDate, "-")[0])
			}

			cTime, convErr := time.Parse(time.ANSIC, tempDate)
			if convErr != nil {
				logger.Log(convErr.Error(), global.StatusError)
			} else {
				tempDate = cTime.String()
				if strings.Contains(tempDate, "+") {
					tempDate = strings.TrimSpace(strings.Split(tempDate, "+")[0])
				} else if strings.Contains(tempDate, "-") {
					tempDate = strings.TrimSpace(strings.Split(tempDate, "-")[0])
				}
				commitDate = tempDate
			}
		}
	}

	commitString := commitHash[:7] + "||" + commitAuthor + "||" + commitDate + "||" + commitMessage
	return commitString
}

func nilCommit() []*string {
	logger.Log("No new commits available to push", global.StatusWarning)
	return nil
}

// UnPushedCommits compares the local branch and the remote branch to extract the commits which are not pushed to the remote
func UnPushedCommits(repo *git.Repository, remoteRef string) []*string {
	var commitArray []*string
	var isAncestor bool
	isAncestor = true

	revHash, _ := repo.ResolveRevision(plumbing.Revision(remoteRef))
	remoteCommit, _ := repo.CommitObject(*revHash)

	head, _ := repo.Head()

	if head == nil {
		return nilCommit()
	}

	localCommit, _ := repo.CommitObject(head.Hash())
	isAncestor, _ = localCommit.IsAncestor(remoteCommit)

	if !isAncestor {
		logItr, _ := repo.Log(&git.LogOptions{
			From:  localCommit.Hash,
			Order: git.LogOrderDFSPost,
			All:   false,
		})

		if logItr == nil {
			if localCommit != nil {
				commitString := commitModel(*localCommit)
				logger.Log(fmt.Sprintf("New commit available for pushing to remote -> %s", commitString), global.StatusInfo)
				commitArray = append(commitArray, &commitString)
				return commitArray
			} else {
				return nilCommit()
			}
		}

		_ = logItr.ForEach(func(commit *object.Commit) error {
			if commit == nil {
				logger.Log("Commit object is nil", global.StatusError)
				return types.Error{Msg: "Commit object is nil"}
			}
			if commit.Hash == remoteCommit.Hash {
				logger.Log(fmt.Sprintf("Same commits found in remote and local trees -> %s", commit.Hash.String()), global.StatusInfo)
				return types.Error{Msg: "Same commit"}
			}

			commitString := commitModel(*commit)
			logger.Log(fmt.Sprintf("New commit available for pushing to remote -> %s", commitString), global.StatusInfo)
			commitArray = append(commitArray, &commitString)
			return nil
		})
		logItr.Close()
		return commitArray
	} else {
		return nilCommit()
	}

}
