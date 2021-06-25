package remote

import (
	"errors"
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

type List interface {
	GetAllRemotes() []*model.RemoteDetails
}

type listRemotes struct {
	repo *git2go.Repository
}

func (l listRemotes) GetAllRemotes() []*model.RemoteDetails {
	var remoteList []*model.RemoteDetails
	repo := l.repo

	if validationErr := l.validateRepo(repo); validationErr != nil {
		return nil
	}

	list, listErr := repo.Remotes.List()
	if listErr != nil {
		logger.Log(listErr.Error(), global.StatusError)
		return nil
	}

	for _, remoteEntry := range list {
		remote, remoteErr := repo.Remotes.Lookup(remoteEntry)
		if remoteErr != nil {
			logger.Log(remoteErr.Error(), global.StatusError)
			continue
		}

		data := model.RemoteDetails{
			RemoteName: remote.Name(),
			RemoteURL:  remote.Url(),
		}
		remoteList = append(remoteList, &data)
	}

	if len(remoteList) == 0 {
		logger.Log(fmt.Sprintf("No remotes available for the repo"), global.StatusWarning)
		return nil
	}

	logger.Log(fmt.Sprintf("Remote data fetched => %+v", remoteList), global.StatusInfo)

	return remoteList
}

func (l listRemotes) validateRepo(repo *git2go.Repository) error {
	if repo == nil {
		logger.Log("Repo is nil", global.StatusError)
		return errors.New("repo is nil")
	}

	if repo.Remotes == (git2go.RemoteCollection{}) {
		logger.Log("Remote collection is nil", global.StatusError)
		return errors.New("repo remote collection is nil")
	}
	return nil
}

func NewRemoteList(repo *git2go.Repository) List {
	return listRemotes{repo: repo}
}
