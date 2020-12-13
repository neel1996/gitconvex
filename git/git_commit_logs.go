package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"go/types"
	"strings"
	"time"
)

// commitOrganizer collects and organizes the commit related information in the GitCommits struct
func commitOrganizer(commits []object.Commit) []*model.GitCommits {
	logger := global.Logger{}
	var commitList []*model.GitCommits
	for _, commit := range commits {
		if !commit.Hash.IsZero() {
			commitHash := commit.Hash.String()
			commitAuthor := strings.Split(commit.Author.String(), " ")[0]
			commitMessage := strings.Split(commit.Message, "\n")[0]
			commitFilesItr, err := commit.Files()
			commitFileCount := 0
			commitDate := ""

			logger.Log(fmt.Sprintf("Fetching commit details for -> %s", commitHash), global.StatusInfo)

			var prevTree *object.Tree
			prevCommit, parentErr := commit.Parents().Next()
			currentTree, _ := commit.Tree()

			if parentErr != nil {
				commitFileCount = 0
			} else {
				prevTree, _ = prevCommit.Tree()
				diff, _ := currentTree.Diff(prevTree)
				commitFileCount = diff.Len()
			}

			// Logic to extract commit timestamp from commit string
			// go-git commit iterator does not provide an option to extract the timestamp directly

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
						commitDate = cTime.String()
					}
				}
			}

			if err != nil {
				logger.Log(err.Error(), global.StatusError)
			} else {
				_ = commitFilesItr.ForEach(func(file *object.File) error {
					return nil
				})
			}

			commitList = append(commitList, &model.GitCommits{
				Hash:             &commitHash,
				Author:           &commitAuthor,
				CommitTime:       &commitDate,
				CommitMessage:    &commitMessage,
				CommitFilesCount: &commitFileCount,
			})
		}
	}

	return commitList
}

// CommitLogs fetches the structured commit logs list for the target repo
// Limits the length of commits to 10 for a single function call
// The referenceCommit is used as a reference to fetch the commits in a linear manner
func CommitLogs(repo *git.Repository, referenceCommit string) *model.GitCommitLogResults {
	var commitLogs []object.Commit
	var logOptions *git.LogOptions

	allCommitChan := make(chan AllCommitData)
	go AllCommits(repo, allCommitChan)
	acc := <-allCommitChan
	totalCommits := acc.TotalCommits

	if referenceCommit == "" {
		head, _ := repo.Head()
		logOptions = &git.LogOptions{
			From:  head.Hash(),
			Order: git.LogOrderDFSPost,
			All:   false,
		}
	} else {
		hash := plumbing.NewHash(referenceCommit)
		logOptions = &git.LogOptions{
			From:  hash,
			Order: git.LogOrderDFSPost,
			All:   false,
		}
	}

	commitItr, commitErr := repo.Log(logOptions)

	var commitCounter int
	commitCounter = 0

	if commitErr == nil {
		_ = commitItr.ForEach(func(commit *object.Commit) error {
			if commitCounter >= 10 || commit == nil {
				return types.Error{Msg: "Commit limit reached"}
			}

			if commit.Hash.String() != referenceCommit {
				commitLogs = append(commitLogs, *commit)
				commitCounter++
			}
			return nil
		})
	}

	if len(commitLogs) == 0 {
		return &model.GitCommitLogResults{
			TotalCommits: &totalCommits,
			Commits:      nil,
		}
	}

	refinedCommits := commitOrganizer(commitLogs)
	return &model.GitCommitLogResults{
		TotalCommits: &totalCommits,
		Commits:      refinedCommits,
	}
}
