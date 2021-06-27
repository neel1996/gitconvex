package initialize

import (
	"context"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/remote"
)

const (
	Repo       = "GIT_REPO"
	RemoteName = "REMOTE_NAME"
	RemoteUrl  = "REMOTE_URL"
)

type Remote struct {
	AddRemote     remote.Add
	DeleteRemote  remote.Delete
	EditRemote    remote.Edit
	ListRemote    remote.List
	RemoteName    remote.Name
	ListRemoteUrl remote.ListRemoteUrl
}

func RemoteObjects(ctx context.Context) Remote {
	repo := ctx.Value(Repo).(*git2go.Repository)
	remoteName := ctx.Value(RemoteName).(string)
	remoteUrl := ctx.Value(RemoteUrl).(string)

	remoteValidation := remote.NewRemoteValidation()

	addRemote := remote.NewAddRemote(repo, remoteName, remoteUrl, remoteValidation)
	deleteRemote := remote.NewDeleteRemote(repo, remoteName, remoteValidation)
	editRemote := remote.NewEditRemote(repo, remoteName, remoteUrl, remoteValidation)
	listRemote := remote.NewRemoteList(repo, remoteValidation)
	name := remote.NewGetRemoteName(repo, remoteUrl, remoteValidation)
	listRemoteURL := remote.NewRemoteUrlData(repo, remoteValidation)

	return Remote{
		AddRemote:     addRemote,
		DeleteRemote:  deleteRemote,
		EditRemote:    editRemote,
		ListRemote:    listRemote,
		ListRemoteUrl: listRemoteURL,
		RemoteName:    name,
	}
}
