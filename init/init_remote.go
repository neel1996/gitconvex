package initialize

import (
	"context"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
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
	r := ctx.Value(Repo).(*git2go.Repository)
	repo := middleware.NewRepository(r)
	remoteName := ctx.Value(RemoteName).(string)
	remoteUrl := ctx.Value(RemoteUrl).(string)

	remoteValidation := remote.NewRemoteValidation(repo)

	listRemote := remote.NewRemoteList(repo, remoteValidation)
	addRemote := remote.NewAddRemote(repo, remoteName, remoteUrl, remoteValidation)
	deleteRemote := remote.NewDeleteRemote(repo, remoteName, remoteValidation)
	editRemote := remote.NewEditRemote(repo, remoteName, remoteUrl, remoteValidation, listRemote)
	name := remote.NewGetRemoteName(repo, remoteUrl, remoteValidation, listRemote)
	listRemoteURL := remote.NewRemoteUrlData(repo, remoteValidation, listRemote)

	return Remote{
		AddRemote:     addRemote,
		DeleteRemote:  deleteRemote,
		EditRemote:    editRemote,
		ListRemote:    listRemote,
		ListRemoteUrl: listRemoteURL,
		RemoteName:    name,
	}
}
