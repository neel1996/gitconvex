package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"strings"
)

type RemoteDataInterface interface {
	GetRemoteHost() *string
	GetRemoteName() string
	RemoteData(remoteChan chan RemoteDataModel)
}

type RemoteDataStruct struct {
	Repo      *git.Repository
	RemoteURL string
}

type RemoteDataModel struct {
	RemoteHost *string
	RemoteURL  []*string
}

// GetRemoteHost returns the remote repository host name based on the remote URL
// e.g. github.com/test.git => returns github
func (r RemoteDataStruct) GetRemoteHost() *string {
	remoteURL := r.RemoteURL
	var remoteHostReference []string
	remoteHostReference = []string{"github", "gitlab", "bitbucket", "azure", "codecommit"}

	for _, host := range remoteHostReference {
		if strings.Contains(remoteURL, host) {
			return &host
		}
	}
	return nil
}

// GetRemoteName function returns the name of the remote based on the supplied remote URL
func (r RemoteDataStruct) GetRemoteName() string {
	var remoteName string
	logger := global.Logger{}

	repo := r.Repo
	remoteURL := r.RemoteURL

	remotes, remoteErr := repo.Remotes()

	if remoteErr != nil {
		logger.Log(remoteErr.Error(), global.StatusError)
	} else {
		for _, remote := range remotes {
			if remote.Config().URLs[0] == remoteURL {
				remoteName = remote.Config().Name
			}
		}
	}
	return remoteName
}

// RemoteData returns the remote host name and the remote URL of the target repo
func (r RemoteDataStruct) RemoteData(remoteChan chan RemoteDataModel) {
	logger := global.Logger{}
	var remoteURL []*string
	repo := r.Repo

	remote, _ := repo.Remotes()
	remoteURL = func() []*string {
		var rUrl []*string
		for _, i := range remote {
			for _, tempUrl := range i.Config().URLs {
				logger.Log(fmt.Sprintf("Available remotes in repo : \n%v", tempUrl), global.StatusInfo)
				rUrl = append(rUrl, &tempUrl)
			}
		}
		return rUrl
	}()

	if len(remoteURL) == 0 {
		nilRemote := "No Remote Host Available"
		nilRemoteURL := ""
		remoteChan <- RemoteDataModel{
			RemoteHost: &nilRemote,
			RemoteURL:  []*string{&nilRemoteURL},
		}
	} else {
		r.RemoteURL = *remoteURL[0]
		remoteChan <- RemoteDataModel{
			RemoteHost: r.GetRemoteHost(),
			RemoteURL:  remoteURL,
		}
	}
	close(remoteChan)
}
