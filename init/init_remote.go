package initialize

import (
	"context"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/constants"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/git/remote"
	"github.com/neel1996/gitconvex/models"
)

func RemoteObjects(ctx context.Context) models.Remote {
	r := ctx.Value(constants.Repo).(*git2go.Repository)
	repo := middleware.NewRepository(r)
	remoteName := ctx.Value(constants.RemoteName).(string)
	remoteUrl := ctx.Value(constants.RemoteUrl).(string)

	remoteValidation := remote.NewRemoteValidation(repo)

	listRemote := remote.NewRemoteList(repo, remoteValidation)
	addRemote := remote.NewAddRemote(repo, remoteName, remoteUrl, remoteValidation)
	deleteRemote := remote.NewDeleteRemote(repo, remoteName, remoteValidation)
	editRemote := remote.NewEditRemote(repo, remoteName, remoteUrl, remoteValidation, listRemote)
	name := remote.NewGetRemoteName(repo, remoteUrl, remoteValidation, listRemote)
	listRemoteURL := remote.NewRemoteUrlData(repo, remoteValidation, listRemote)

	return models.Remote{
		AddRemote:     addRemote,
		DeleteRemote:  deleteRemote,
		EditRemote:    editRemote,
		ListRemote:    listRemote,
		ListRemoteUrl: listRemoteURL,
		RemoteName:    name,
	}
}
