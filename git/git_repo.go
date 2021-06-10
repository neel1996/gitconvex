package git

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/utils"
)

type RepoInterface interface {
	Repo(repoChan chan RepoDetails)
}

type RepoStruct struct {
	RepoId string
}

type RepoDetails struct {
	RepoId     string
	RepoPath   string
	RepoName   string
	TimeStamp  string
	AuthOption string
	UserName   string
	Password   string
	SSHKeyPath string
	GitRepo    *git2go.Repository
}

func handlePanic() {
	panicMsg := recover()
	if panicMsg != nil {
		logger.Log("Required fields not received", global.StatusWarning)
	}
}

// Repo function gets the repoId and returns the respective git.Repository object along with additional repo metadata
func (r RepoStruct) Repo(repoChan chan RepoDetails) {
	var repoData []utils.RepoData
	var (
		repoName   string
		repoPath   string
		authOption string
		sshKeyPath string
		userName   string
		password   string
		timeStamp  string
	)
	repoId := r.RepoId

	defer handlePanic()
	if repoId == "" || repoChan == nil {
		close(repoChan)
		panic("Required fields not received")
	}

	repoData = utils.DataStoreFileReader()

	for _, repo := range repoData {
		if repo.RepoId == repoId {
			repoName = repo.RepoName
			authOption = repo.AuthOption
			sshKeyPath = repo.SSHKeyPath
			userName = repo.UserName
			password = repo.Password
			repoPath = repo.RepoPath
			timeStamp = repo.TimeStamp
			break
		}
	}

	git2goRepo, err := git2go.OpenRepository(repoPath)

	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		repoChan <- RepoDetails{
			RepoId:   repoId,
			RepoPath: repoPath,
			GitRepo:  nil,
		}
	} else {
		repoChan <- RepoDetails{
			RepoId:     repoId,
			RepoPath:   repoPath,
			RepoName:   repoName,
			TimeStamp:  timeStamp,
			AuthOption: authOption,
			UserName:   userName,
			Password:   password,
			SSHKeyPath: sshKeyPath,
			GitRepo:    git2goRepo,
		}
	}
	close(repoChan)
}
