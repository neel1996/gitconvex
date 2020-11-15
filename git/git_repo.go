package git

import (
	git "github.com/go-git/go-git/v5"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/utils"
)

type RepoDetails struct {
	RepoId   string
	RepoPath string
	GitRepo  *git.Repository
}

func handlePanic() {
	logger := global.Logger{}

	panicMsg := recover()
	if panicMsg != nil {
		logger.Log("Required fields not received", global.StatusWarning)
	}
}

// Repo function gets the repoId and returns the respective git.Repository object along with additional repo metadata

func Repo(repoId string, repoChan chan RepoDetails) {
	var repoData []utils.RepoData
	var repoPath string
	logger := global.Logger{}

	defer handlePanic()
	if repoId == "" || repoChan == nil {
		close(repoChan)
		panic("Required fields not received")
	}

	repoData = utils.DataStoreFileReader()

	for _, repo := range repoData {
		if repo.RepoId == repoId {
			repoPath = repo.RepoPath
			break
		}
	}

	repository, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
	} else {
		repoChan <- RepoDetails{
			RepoId:   repoId,
			RepoPath: repoPath,
			GitRepo:  repository,
		}
	}

	close(repoChan)
}
