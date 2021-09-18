package api

import (
	"github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/neel1996/gitconvex/git/commit"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/git/remote"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/neel1996/gitconvex/utils"
	"strings"
)

//TODO: This is a temporary function and will be replaced with an clean interface over the course of refactoring

// RepoStatus collects the basic details of the target repo and returns the consolidated result
func RepoStatus(repoId string) *model.GitRepoStatusResults {
	logger.Log("Collecting repo status information", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	//commitChan := make(chan commit.AllCommitData)
	trackedFileCountChan := make(chan int)

	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoId}
	go repoObject.Repo(repoChan)

	var repoName *string
	r := <-repoChan
	repo := r.GitRepo

	if repo == nil {
		return &model.GitRepoStatusResults{
			GitRemoteData:        nil,
			GitRepoName:          nil,
			GitBranchList:        nil,
			GitAllBranchList:     nil,
			GitCurrentBranch:     nil,
			GitRemoteHost:        nil,
			GitTotalCommits:      nil,
			GitLatestCommit:      nil,
			GitTotalTrackedFiles: nil,
		}
	}

	tempRemote := ""
	var (
		remoteURL  *string
		remoteName string
	)
	remoteURL = &tempRemote

	repoForRemote := middleware.NewRepository(repo)

	remoteValidation := remote.NewRemoteValidation(repoForRemote)
	remoteList := remote.NewRemoteList(repoForRemote, remoteValidation)
	remoteUrlList := remote.NewRemoteUrlData(repoForRemote, remoteValidation, remoteList)
	listRemoteUrl := remote.Operation{ListRemoteUrl: remoteUrlList}

	remotes, remoteErr := listRemoteUrl.GitGetAllRemoteUrl()
	if remoteErr != nil {
		logger.Log(remoteErr.Error(), global.StatusError)
		return nil
	}

	if len(remotes) > 0 && *remotes[0] != "" {
		remoteNameObject := remote.NewGetRemoteName(
			repoForRemote,
			*remotes[0],
			remoteValidation,
			remoteList,
		)
		remoteName = remoteNameObject.GetRemoteNameWithUrl()
		sRemote := strings.Split(*remotes[0], "/")
		repoName = &sRemote[len(sRemote)-1]
	} else {
		nilRemote := "No Remotes Available"
		remotes[0] = &nilRemote
		repoData := utils.DataStoreFileReader()

		for _, repo := range repoData {
			if repo.RepoId == repoId {
				repoName = &repo.RepoName
			}
		}
	}

	if len(remotes) > 1 {
		var tempRemoteArray []string
		for _, ptrRemote := range remotes {
			tempRemoteArray = append(tempRemoteArray, *ptrRemote)
		}
		*remoteURL = strings.Join(tempRemoteArray, "||")
	} else {
		*remoteURL = *remotes[0]
	}

	branchListObj := branch.NewBranchList(middleware.NewRepository(repo))
	branchList, _ := branchListObj.ListBranches()

	currentBranch := branchList.CurrentBranch
	branches := branchList.BranchList
	allBranches := branchList.AllBranchList

	listCommitLogs := commit.NewListAllLogs(middleware.NewRepository(repo), nil, nil)
	totalCommits := commit.NewTotalCommits(listCommitLogs).Get()

	latestCommitMessage := commit.NewLatestMessage(middleware.NewRepository(repo)).Get()

	var listFilesObject git.ListFilesInterface
	listFilesObject = git.ListFilesStruct{
		Repo: repo,
	}

	go listFilesObject.TrackedFileCount(trackedFileCountChan)
	trackedFileCount := <-trackedFileCountChan
	trackedFilePtr := &trackedFileCount

	return &model.GitRepoStatusResults{
		GitRemoteData:        remoteURL,
		GitRepoName:          repoName,
		GitBranchList:        utils.GeneratePointerArrayFrom(branches),
		GitAllBranchList:     utils.GeneratePointerArrayFrom(allBranches),
		GitCurrentBranch:     &currentBranch,
		GitRemoteHost:        &remoteName,
		GitTotalCommits:      &totalCommits,
		GitLatestCommit:      &latestCommitMessage,
		GitTotalTrackedFiles: trackedFilePtr,
	}
}
