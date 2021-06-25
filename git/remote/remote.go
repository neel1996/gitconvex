package remote

import (
	"errors"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

var logger global.Logger

type Remote interface {
	GitAddRemote() (*model.RemoteMutationResult, error)
	GitDeleteRemote() (*model.RemoteMutationResult, error)
	GitGetAllRemote() ([]*model.RemoteDetails, error)
	GitGetRemoteHostName() (string, error)
	GitRemoteName() (string, error)
	GitGetAllRemoteUrl() ([]*string, error)
}

type Operation struct {
	Add           Add
	Delete        Delete
	Host          Host
	List          List
	Name          Name
	ListRemoteUrl ListRemoteUrl
}

func (r Operation) GitAddRemote() (*model.RemoteMutationResult, error) {
	addRemoteErr := r.Add.NewRemote()

	if addRemoteErr != nil {
		return nil, errors.New(global.RemoteAddError)
	}

	return &model.RemoteMutationResult{Status: global.RemoteAddSuccess}, nil
}

func (r Operation) GitDeleteRemote() (*model.RemoteMutationResult, error) {
	err := r.Delete.DeleteRemote()

	if err != nil {
		return nil, errors.New(global.RemoteDeleteError)
	}

	return &model.RemoteMutationResult{Status: global.RemoteDeleteSuccess}, nil
}

func (r Operation) GitGetRemoteHostName() (string, error) {
	remoteHost := r.Host.GetRemoteHostForUrl()

	if remoteHost == "" {
		return "", errors.New("unable to find a matching host name")
	}

	return remoteHost, nil
}

func (r Operation) GitGetAllRemote() ([]*model.RemoteDetails, error) {
	allRemotes := r.List.GetAllRemotes()

	if allRemotes == nil {
		return nil, errors.New("unable to fetch remotes from the repo")
	}

	return allRemotes, nil
}

func (r Operation) GitRemoteName() (string, error) {
	remoteHost := r.Name.GetRemoteNameWithUrl()

	if remoteHost == "" {
		return "", errors.New("no matching remote found for the URL")
	}

	return remoteHost, nil
}

func (r Operation) GitGetAllRemoteUrl() ([]*string, error) {
	remoteUrlList := r.ListRemoteUrl.GetAllRemoteUrl()

	if remoteUrlList == nil {
		return nil, errors.New("unable to fetch remote URLs from the repo")
	}

	return remoteUrlList, nil
}
