package commit

import (
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

var logger global.Logger

type Commit interface {
	GitCommitChange() (string, error)
	GitTotalCommits() int
	GitCommitLogs() ([]*model.GitCommits, error)
	GitCommitFileHistory(string) ([]*model.GitCommitFileResult, error)
	GitSearchCommitLogs(string, string) ([]*model.GitCommits, error)
}

type Operation struct {
	Changes     Changes
	Total       Total
	ListAllLogs ListAllLogs
	FileHistory FileHistory
	Lookup      Lookup
	SearchLogs  SearchLogs
}

func (c Operation) GitCommitChange() (string, error) {
	err := c.Changes.Add()

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return "", ChangeError
	}

	return global.CommitChangeSuccess, nil
}

func (c Operation) GitTotalCommits() int {
	return c.Total.Get()
}

func (c Operation) GitCommitLogs() ([]*model.GitCommits, error) {
	commits, logsErr := c.ListAllLogs.Get()
	if logsErr != nil {
		logger.Log(logsErr.Error(), global.StatusError)
		return nil, LogsError
	}

	return NewMapper(c.FileHistory).Map(commits), nil
}

func (c Operation) GitCommitFileHistory(commitHash string) ([]*model.GitCommitFileResult, error) {
	if commitHash == "" {
		return nil, EmptyCommitHashError
	}

	commit, err := c.Lookup.WithReferenceId(commitHash)
	if err != nil {
		return nil, err
	}

	return c.FileHistory.Get(commit)
}

func (c Operation) GitSearchCommitLogs(searchType string, searchKey string) ([]*model.GitCommits, error) {
	return c.SearchLogs.GetMatchingLogs(searchType, searchKey)
}
