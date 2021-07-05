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

	addRemote := remote.NewAddRemote(repo, remoteName, remoteUrl)
	deleteRemote := remote.NewDeleteRemote(repo, remoteName)
	editRemote := remote.NewEditRemote(repo, remoteName, remoteUrl)
	listRemote := remote.NewRemoteList(repo)
	name := remote.NewGetRemoteName(repo, remoteUrl)
	listRemoteURL := remote.NewRemoteUrlData(repo)

	return Remote{
		AddRemote:     addRemote,
		DeleteRemote:  deleteRemote,
		EditRemote:    editRemote,
		ListRemote:    listRemote,
		ListRemoteUrl: listRemoteURL,
		RemoteName:    name,
	}
}
