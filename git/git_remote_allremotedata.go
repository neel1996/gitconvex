package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"strings"
)

type RemoteDataInterface interface {
	GetRemoteHost() *string
	GetRemoteName() string
	GetAllRemotes() []*model.RemoteDetails
	RemoteData(remoteChan chan RemoteDataModel)
}

type RemoteDataStruct struct {
	Repo      *git2go.Repository
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
	repo := r.Repo
	remoteURL := r.RemoteURL
	remoteList, _ := repo.Remotes.List()

	for _, remoteEntry := range remoteList {
		remoteCollection, _ := repo.Remotes.Lookup(remoteEntry)
		if remoteCollection != nil {
			url := remoteCollection.Url()
			if url == remoteURL {
				remoteName := remoteCollection.Name()
				logger.Log(fmt.Sprintf("Remote Name - %s for the url - %s", remoteName, remoteURL), global.StatusInfo)
				return remoteName
			}
		}
	}
	logger.Log(fmt.Sprintf("Unable to find a suitable remote entry for the url - %s", remoteURL), global.StatusError)
	return ""
}

// RemoteData returns the remote host name and the remote URL of the target repo
func (r RemoteDataStruct) RemoteData(remoteChan chan RemoteDataModel) {
	var remoteURL []*string
	repo := r.Repo

	remotes, _ := repo.Remotes.List()

	for _, i := range remotes {
		remote, _ := repo.Remotes.Lookup(i)
		if remote != nil {
			url := remote.Url()
			remoteURL = append(remoteURL, &url)
		}
	}

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

// GetAllRemotes returns all the remotes and their corresponding URLs from the target repo
func (r RemoteDataStruct) GetAllRemotes() []*model.RemoteDetails {
	var allRemoteData []*model.RemoteDetails

	repo := r.Repo
	remoteList, remoteListErr := repo.Remotes.List()
	if remoteListErr != nil {
		logger.Log(remoteListErr.Error(), global.StatusError)
		return nil
	}

	for _, remoteEntry := range remoteList {
		if remoteEntry != "" {
			remote, remoteErr := repo.Remotes.Lookup(remoteEntry)
			if remoteErr != nil {
				logger.Log(remoteErr.Error(), global.StatusError)
				continue
			}
			data := model.RemoteDetails{
				RemoteName: remote.Name(),
				RemoteURL:  remote.Url(),
			}
			logger.Log(fmt.Sprintf("Remote data fetched => %+v", data), global.StatusInfo)
			allRemoteData = append(allRemoteData, &data)
		}
	}
	return allRemoteData
}
