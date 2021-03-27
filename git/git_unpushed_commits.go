package git

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
)

type UnPushedCommitInterface interface {
	UnPushedCommits() *model.UnPushedCommitResult
}

type UnPushedCommitStruct struct {
	Repo      *git2go.Repository
	RemoteRef string
	LocalRef  string
}

// commitModel function generates a pipe separated string with the required commit details
// The resultant string will be sent to the client
func commitModel(commit *git2go.Commit) *model.GitCommits {
	commitHash := commit.Id().String()[:7]
	commitAuthor := commit.Author().Name
	commitMessage := commit.Message()
	commitDate := commit.Author().When.String()

	return &model.GitCommits{
		Hash:          &commitHash,
		Author:        &commitAuthor,
		CommitTime:    &commitDate,
		CommitMessage: &commitMessage,
	}
}

// UnPushedCommits compares the local branch and the remote branch to extract the commits which are not pushed to the remote
func (u UnPushedCommitStruct) UnPushedCommits() *model.UnPushedCommitResult {
	repo := u.Repo
	remoteRef := u.RemoteRef
	var commitArray []*model.GitCommits

	// Returning nil commit response if repo has no HEAD
	head, _ := repo.Head()
	if head == nil {
		logger.Log("HEAD is NULL", global.StatusError)
		return &model.UnPushedCommitResult{
			IsNewBranch: false,
			GitCommits:  []*model.GitCommits{},
		}
	}

	remoteBranch, remoteBranchErr := repo.LookupBranch(remoteRef, git2go.BranchRemote)
	if remoteBranchErr != nil {
		logger.Log(remoteBranchErr.Error(), global.StatusError)
		logger.Log("Treating remote branch as a newly created branch", global.StatusWarning)

		return &model.UnPushedCommitResult{
			IsNewBranch: true,
			GitCommits:  nil,
		}
	}

	// Checking if both branches have any varying commits
	diff := head.Cmp(remoteBranch.Reference)
	if diff == 0 {
		logger.Log("No difference between remote and local branches", global.StatusError)
		return &model.UnPushedCommitResult{
			IsNewBranch: false,
			GitCommits:  []*model.GitCommits{},
		}
	}

	localCommit, _ := repo.LookupCommit(head.Target())
	remoteCommit, _ := repo.LookupCommit(remoteBranch.Target())

	if localCommit != nil && remoteCommit != nil {
		commonAncestor, _ := repo.MergeBase(localCommit.Id(), remoteCommit.Id())
		if commonAncestor != nil {
			commitArray = append(commitArray, commitModel(localCommit))
			n := localCommit.ParentCount()
			var i uint
			for i = 0; i < n; i++ {
				currentCommit := localCommit.Parent(i)
				if currentCommit != nil && currentCommit.Id().String() != commonAncestor.String() {
					commitArray = append(commitArray, commitModel(currentCommit))
				} else {
					break
				}
			}
		} else {
			logger.Log("No new commits available to push", global.StatusWarning)
			return nil
		}
		return &model.UnPushedCommitResult{
			IsNewBranch: false,
			GitCommits:  commitArray,
		}
	} else {
		logger.Log("No new commits available to push", global.StatusWarning)
		return &model.UnPushedCommitResult{
			IsNewBranch: false,
			GitCommits:  []*model.GitCommits{},
		}
	}
}
