package controller

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/git/commit"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

type CommitLogController interface {
	GetCommitLogs() (*model.GitCommitLogResults, error)
}
type commitLogController struct {
	repoId          string
	referenceCommit string
}

func (c commitLogController) GetCommitLogs() (*model.GitCommitLogResults, error) {
	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: c.repoId}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	repository := middleware.NewRepository(repo.GitRepo)

	list := commit.NewListAllLogs(repository, nil, c.referenceId())
	total := commit.NewTotalCommits(list)
	fileHistory := commit.NewFileHistory(repository)

	operation := commit.Operation{
		ListAllLogs: list,
		Total:       total,
		FileHistory: fileHistory,
	}

	commits, err := operation.GitCommitLogs()
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return nil, err
	}
	totalCommits := operation.GitTotalCommits()

	return &model.GitCommitLogResults{
		TotalCommits: &totalCommits,
		Commits:      commits,
	}, nil
}

func (c commitLogController) referenceId() *git2go.Oid {
	if c.referenceCommit == "" {
		return nil
	}

	referenceId, oidErr := git2go.NewOid(c.referenceCommit)
	if oidErr != nil {
		logger.Log(oidErr.Error(), global.StatusError)
		return nil
	}

	return referenceId
}

func NewCommitLogController(repoId string, referenceCommit string) CommitLogController {
	return commitLogController{
		repoId:          repoId,
		referenceCommit: referenceCommit,
	}
}
