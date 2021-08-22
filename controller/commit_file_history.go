package controller

import (
	"github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/git/commit"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/graph/model"
)

type CommitFileHistoryController interface {
	GetCommitFileHistory() ([]*model.GitCommitFileResult, error)
}

type commitFileHistoryController struct {
	repoId     string
	commitHash string
}

func (c commitFileHistoryController) GetCommitFileHistory() ([]*model.GitCommitFileResult, error) {
	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: c.repoId}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	repository := middleware.NewRepository(repo.GitRepo)

	fileHistory := commit.NewFileHistory(repository)
	commitLookup := commit.NewLookup(repository)
	operation := commit.Operation{FileHistory: fileHistory, Lookup: commitLookup}

	return operation.GitCommitFileHistory(c.commitHash)
}

func NewCommitFileHistoryController(repoId string, commitHash string) CommitFileHistoryController {
	return commitFileHistoryController{
		repoId:     repoId,
		commitHash: commitHash,
	}
}
