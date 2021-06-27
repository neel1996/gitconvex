package api

import (
	"github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/neel1996/gitconvex/git/remote"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/neel1996/gitconvex/utils"
	"strings"
)

// RepoStatus collects the basic details of the target repo and returns the consolidated result
func RepoStatus(repoId string) *model.GitRepoStatusResults {
	logger.Log("Collecting repo status information", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	branchChan := make(chan branch.ListOfBranches)
	commitChan := make(chan git.AllCommitData)
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

	remoteValidation := remote.NewRemoteValidation()
	remoteUrlList := remote.NewRemoteUrlData(repo, remoteValidation)
	listRemoteUrl := remote.Operation{ListRemoteUrl: remoteUrlList}

	remotes, remoteErr := listRemoteUrl.GitGetAllRemoteUrl()
	if remoteErr != nil {
		logger.Log(remoteErr.Error(), global.StatusError)
		return nil
	}

	if len(remotes) > 0 && *remotes[0] != "" {
		remoteNameObject := remote.NewGetRemoteName(repo, *remotes[0], remoteValidation)
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

	branchListObj := branch.NewBranchListInterface(repo)
	b := branch.Operation{List: branchListObj}
	go b.GitListBranches(branchChan)

	branchList := <-branchChan
	currentBranch := &branchList.CurrentBranch
	branches := branchList.BranchList
	allBranches := branchList.AllBranchList

	var latestCommit *string

	var allCommitObject git.AllCommitInterface

	allCommitObject = git.AllCommitStruct{Repo: repo}
	go allCommitObject.AllCommits(commitChan)
	commitData := <-commitChan
	latestCommit = &commitData.LatestCommit
	totalCommits := commitData.TotalCommits
	totalCommitsPtr := &totalCommits

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
		GitBranchList:        branches,
		GitAllBranchList:     allBranches,
		GitCurrentBranch:     currentBranch,
		GitRemoteHost:        &remoteName,
		GitTotalCommits:      totalCommitsPtr,
		GitLatestCommit:      latestCommit,
		GitTotalTrackedFiles: trackedFilePtr,
	}
}
