package git

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"strings"
)

type CommitLogInterface interface {
	CommitLogs() *model.GitCommitLogResults
}

type CommitLogStruct struct {
	Repo            *git2go.Repository
	ReferenceCommit string
}

// commitOrganizer collects and organizes the commit related information in the GitCommits struct
func commitOrganizer(repo *git2go.Repository, commits []git2go.Commit) []*model.GitCommits {
	var commitList []*model.GitCommits
	for _, commit := range commits {
		if commit.Id().String() != "" {
			commitHash := commit.Id().String()
			commitAuthor := strings.Split(commit.Author().Name, " ")[0]
			commitMessage := strings.Split(commit.Message(), "\n")[0]

			parentCommit := commit.Parent(0)
			if parentCommit == nil {
				logger.Log("NIL Commit encountered", global.StatusWarning)
				continue
			}
			parentTree, _ := parentCommit.Tree()
			currentTree, _ := commit.Tree()

			commitFileCount := 0
			commitDate := commit.Committer().When.String()

			if parentTree != nil && currentTree != nil {
				diff, diffErr := repo.DiffTreeToTree(parentTree, currentTree, nil)

				if diffErr == nil {
					commitFileCount, _ = diff.NumDeltas()
				} else {
					logger.Log(diffErr.Error(), global.StatusError)
				}
			} else {
				logger.Log("Unable to fetch commit tree", global.StatusError)
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
func (c CommitLogStruct) CommitLogs() *model.GitCommitLogResults {
	var commitLogs []git2go.Commit

	repo := c.Repo
	referenceCommit := c.ReferenceCommit

	allCommitChan := make(chan AllCommitData)

	var allCommitsObject AllCommitInterface
	allCommitsObject = AllCommitStruct{Repo: repo}
	go allCommitsObject.AllCommits(allCommitChan)

	acc := <-allCommitChan
	totalCommits := acc.TotalCommits

	var counter int
	counter = 0

	if referenceCommit == "" {
		head, _ := repo.Head()
		if head != nil {
			nxt, _ := repo.LookupCommit(head.Target())
			for nxt != nil && counter <= 10 {
				commitLogs = append(commitLogs, *nxt)
				nxt = nxt.Parent(0)
				counter++
			}
		} else {
			logger.Log("Unable to fetch repo HEAD", global.StatusError)
			return &model.GitCommitLogResults{
				TotalCommits: &totalCommits,
				Commits:      nil,
			}
		}
	} else {
		refId, _ := git2go.NewOid(referenceCommit)
		refCommit, refCommitErr := repo.LookupCommit(refId)

		if refCommitErr == nil {
			nxt := refCommit.Parent(0)
			for nxt != nil && counter <= 10 {
				commitLogs = append(commitLogs, *nxt)
				nxt = nxt.Parent(0)
				counter++
			}
		} else {
			logger.Log(refCommitErr.Error(), global.StatusError)
			return &model.GitCommitLogResults{
				TotalCommits: &totalCommits,
				Commits:      nil,
			}
		}
	}

	if len(commitLogs) == 0 {
		return &model.GitCommitLogResults{
			TotalCommits: &totalCommits,
			Commits:      nil,
		}
	}

	refinedCommits := commitOrganizer(repo, commitLogs)
	return &model.GitCommitLogResults{
		TotalCommits: &totalCommits,
		Commits:      refinedCommits,
	}
}
