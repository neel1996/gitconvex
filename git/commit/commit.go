package commit

import "github.com/neel1996/gitconvex/global"

var logger global.Logger

type Commit interface {
	GitCommitChange() (string, error)
}

type Operation struct {
	Changes Changes
}

func (c Operation) GitCommitChange() (string, error) {
	err := c.Changes.Add()

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return global.CommitChangeError, err
	}

	return global.CommitChangeSuccess, nil
}
