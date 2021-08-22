package controller

import (
	"github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/git/commit"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/graph/model"
)

type CommitLogSearchController interface {
	GetMatchingCommits() ([]*model.GitCommits, error)
}

type commitLogSearchController struct {
	repoId     string
	searchType string
	searchKey  string
}

func (c commitLogSearchController) GetMatchingCommits() ([]*model.GitCommits, error) {
	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: c.repoId}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	repository := middleware.NewRepository(repo.GitRepo)

	list := commit.NewListAllLogs(repository, nil, nil)
	fileHistory := commit.NewFileHistory(repository)
	mapper := commit.NewMapper(fileHistory)
	searchLogs := commit.NewSearchLogs(repository, list, mapper)

	operation := commit.Operation{SearchLogs: searchLogs}

	return operation.GitSearchCommitLogs(c.searchType, c.searchKey)
}

func NewCommitLogSearchController(repoId string, searchType string, searchKey string) CommitLogSearchController {
	return commitLogSearchController{
		repoId:     repoId,
		searchType: searchType,
		searchKey:  searchKey,
	}
}
