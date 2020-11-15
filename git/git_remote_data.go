package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"strings"
)

type RemoteDataModel struct {
	RemoteHost *string
	RemoteURL  []*string
}

// GetRemoteHost returns the remote repository host name based on the remote URL
// e.g. github.com/test.git => returns github

func GetRemoteHost(remoteURL string) *string {
	var remoteHostReference []string
	remoteHostReference = []string{"github", "gitlab", "bitbucket", "azure", "codecommit"}

	for _, host := range remoteHostReference {
		if strings.Contains(remoteURL, host) {
			return &host
		}
	}
	return nil
}

// RemoteData returns the remote host name and the remote URL of the target repo

func RemoteData(repo *git.Repository, remoteChan chan RemoteDataModel) {
	logger := global.Logger{}
	var remoteURL []*string

	remote, _ := repo.Remotes()
	remoteURL = func() []*string {
		var rUrl []*string
		for _, i := range remote {
			for _, tempUrl := range i.Config().URLs {
				rUrl = append(rUrl, &tempUrl)
			}
		}
		return rUrl
	}()

	logger.Log(fmt.Sprintf("Available remotes in repo : \n%v", remoteURL), global.StatusInfo)
	remoteChan <- RemoteDataModel{
		RemoteHost: GetRemoteHost(*remoteURL[0]),
		RemoteURL:  remoteURL,
	}

	close(remoteChan)
}
