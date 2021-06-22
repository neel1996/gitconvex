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
}

type Operation struct {
	Add    Add
	Delete Delete
}

func (r Operation) GitAddRemote() (*model.RemoteMutationResult, error) {
	addRemoteResult := r.Add.NewRemote()

	if addRemoteResult.Status == global.RemoteAddError {
		return nil, errors.New(global.RemoteAddError)
	}

	return addRemoteResult, nil
}

func (r Operation) GitDeleteRemote() (*model.RemoteMutationResult, error) {
	deleteRemoteResult := r.Delete.DeleteRemote()

	if deleteRemoteResult.Status == global.RemoteDeleteError {
		return nil, errors.New(global.RemoteDeleteError)
	}

	return deleteRemoteResult, nil
}
