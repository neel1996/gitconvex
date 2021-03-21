package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/neel1996/gitconvex-server/api"
	"github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/generated"
	"github.com/neel1996/gitconvex-server/graph/model"
)

func (r *mutationResolver) AddRepo(ctx context.Context, repoName string, repoPath string, cloneSwitch bool, repoURL *string, initSwitch bool, authOption string, sshKeyPath *string, userName *string, password *string) (*model.AddRepoParams, error) {
	var addRepoObject api.AddRepoInterface
	addRepoObject = api.AddRepoInputs{
		RepoName:    repoName,
		RepoPath:    repoPath,
		CloneSwitch: cloneSwitch,
		RepoURL:     *repoURL,
		InitSwitch:  initSwitch,
		AuthOption:  authOption,
		UserName:    *userName,
		Password:    *password,
		SSHKeyPath:  *sshKeyPath,
	}
	return addRepoObject.AddRepo(), nil
}

func (r *mutationResolver) AddBranch(ctx context.Context, repoID string, branchName string) (string, error) {
	logger.Log("Initiating branch addition request", global.StatusInfo)
	repoChan := make(chan git.RepoDetails)

	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid. Branch addition failed", global.StatusError)
		return global.BranchAddError, nil
	}

	var addBranchObj git.AddBranchInterface
	addBranchObj = git.AddBranchInput{
		Repo:       repo.GitRepo,
		BranchName: branchName,
	}
	return addBranchObj.AddBranch(), nil
}

func (r *mutationResolver) CheckoutBranch(ctx context.Context, repoID string, branchName string) (string, error) {
	logger.Log("Initiating branch checkout request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)

	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)

	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repository is invalid or HEAD is null", global.StatusError)
		return global.BranchCheckoutError, nil
	}

	var checkoutObject git.BranchCheckoutInterface
	checkoutObject = git.BranchCheckoutInputs{
		Repo:       repo.GitRepo,
		BranchName: branchName,
	}
	return checkoutObject.CheckoutBranch(), nil
}

func (r *mutationResolver) DeleteBranch(ctx context.Context, repoID string, branchName string, forceFlag bool) (*model.BranchDeleteStatus, error) {
	logger.Log("Initiating branch deletion request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is null", global.StatusError)
		return &model.BranchDeleteStatus{
			Status: global.BranchDeleteError,
		}, nil
	}

	var deleteBranchObject git.DeleteBranchInterface
	deleteBranchObject = git.DeleteBranchInputs{
		Repo:       repo.GitRepo,
		BranchName: branchName,
	}
	return deleteBranchObject.DeleteBranch(), nil
}

func (r *mutationResolver) FetchFromRemote(ctx context.Context, repoID string, remoteURL *string, remoteBranch *string) (*model.FetchResult, error) {
	logger.Log("Initiating fetch from remote request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)

		return &model.FetchResult{
			Status:       global.FetchFromRemoteError,
			FetchedItems: nil,
		}, nil
	}

	var remoteDataObject git.RemoteDataInterface
	remoteDataObject = git.RemoteDataStruct{
		Repo:      repo.GitRepo,
		RemoteURL: *remoteURL,
	}

	var fetchObject git.FetchInterface
	fetchObject = git.FetchStruct{
		Repo:         repo.GitRepo,
		RemoteName:   remoteDataObject.GetRemoteName(),
		RepoPath:     repo.RepoPath,
		RemoteURL:    *remoteURL,
		RemoteBranch: *remoteBranch,
		RepoName:     repo.RepoName,
		AuthOption:   repo.AuthOption,
		UserName:     repo.UserName,
		Password:     repo.Password,
		SSHKeyPath:   repo.SSHKeyPath,
	}
	return fetchObject.FetchFromRemote(), nil
}

func (r *mutationResolver) PullFromRemote(ctx context.Context, repoID string, remoteURL *string, remoteBranch *string) (*model.PullResult, error) {
	logger.Log("Initiating pull from remote request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return &model.PullResult{
			Status:      global.PullFromRemoteError,
			PulledItems: nil,
		}, nil
	}

	var pullObject git.PullInterface
	pullObject = git.PullStruct{
		Repo:         repo.GitRepo,
		RemoteURL:    *remoteURL,
		RemoteBranch: *remoteBranch,
		RepoPath:     repo.RepoPath,
		RepoName:     repo.RepoName,
		AuthOption:   repo.AuthOption,
		UserName:     repo.UserName,
		Password:     repo.Password,
		SSHKeyPath:   repo.SSHKeyPath,
	}
	return pullObject.PullFromRemote(), nil
}

func (r *mutationResolver) StageItem(ctx context.Context, repoID string, item string) (string, error) {
	logger.Log("Initiating stage item request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return global.StageItemError, nil
	}

	var stageItemObject git.StageItemInterface
	stageItemObject = git.StageItemStruct{
		Repo:     repo.GitRepo,
		FileItem: item,
	}
	return stageItemObject.StageItem(), nil
}

func (r *mutationResolver) RemoveStagedItem(ctx context.Context, repoID string, item string) (string, error) {
	logger.Log("Initiating remove item request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return global.RemoveItemError, nil
	}

	var resetObject git.ResetInterface
	resetObject = git.ResetStruct{
		Repo:     repo.GitRepo,
		RepoPath: repo.RepoPath,
		FileItem: item,
	}
	return resetObject.RemoveItem(), nil
}

func (r *mutationResolver) RemoveAllStagedItem(ctx context.Context, repoID string) (string, error) {
	logger.Log("Initiating remove all items request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return global.RemoveAllItemsError, nil
	}

	var resetAllObject git.ResetAllInterface
	resetAllObject = git.ResetAllStruct{Repo: repo.GitRepo}
	return resetAllObject.ResetAllItems(), nil
}

func (r *mutationResolver) StageAllItems(ctx context.Context, repoID string) (string, error) {
	logger.Log("Initiating stage all items request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return global.StageAllItemsError, nil
	}

	var stageAllObject git.StageAllInterface
	stageAllObject = git.StageAllStruct{Repo: repo.GitRepo}
	return stageAllObject.StageAllItems(), nil
}

func (r *mutationResolver) CommitChanges(ctx context.Context, repoID string, commitMessage string) (string, error) {
	logger.Log("Initiating commit changes request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	var commitObject git.CommitInterface
	commitObject = git.CommitStruct{
		Repo:          repo.GitRepo,
		CommitMessage: commitMessage,
		RepoPath:      repo.RepoPath,
	}

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or worktree is null", global.StatusError)
		return global.CommitChangeError, nil
	}
	return commitObject.CommitChanges(), nil
}

func (r *mutationResolver) PushToRemote(ctx context.Context, repoID string, remoteHost string, branch string) (string, error) {
	logger.Log("Initiating push to remote request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	var remoteDataObject git.RemoteDataInterface
	remoteDataObject = git.RemoteDataStruct{
		Repo:      repo.GitRepo,
		RemoteURL: remoteHost,
	}

	remoteName := remoteDataObject.GetRemoteName()

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil || remoteName == "" {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return global.PushToRemoteError, nil
	}

	var pushObject git.PushInterface
	pushObject = git.PushStruct{
		Repo:         repo.GitRepo,
		RepoName:     repo.RepoName,
		AuthOption:   repo.AuthOption,
		UserName:     repo.UserName,
		Password:     repo.Password,
		SSHKeyPath:   repo.SSHKeyPath,
		RemoteName:   remoteName,
		RemoteBranch: branch,
		RepoPath:     repo.RepoPath,
	}
	return pushObject.PushToRemote(), nil
}

func (r *mutationResolver) SettingsEditPort(ctx context.Context, newPort string) (string, error) {
	logger.Log("Initiating update port number request", global.StatusInfo)
	return api.UpdatePortNumber(newPort), nil
}

func (r *mutationResolver) UpdateRepoDataFile(ctx context.Context, newDbFile string) (string, error) {
	logger.Log("Initiating update data file request", global.StatusInfo)
	return api.UpdateDBFilePath(newDbFile), nil
}

func (r *mutationResolver) DeleteRepo(ctx context.Context, repoID string) (*model.DeleteStatus, error) {
	logger.Log("Initiating delete repo request", global.StatusInfo)
	return api.DeleteRepo(repoID), nil
}

func (r *mutationResolver) UpdateRepoName(ctx context.Context, repoID string, repoName string) (string, error) {
	logger.Log("Initiating repo name update request", global.StatusInfo)
	return api.UpdateRepoName(repoID, repoName)
}

func (r *mutationResolver) AddRemote(ctx context.Context, repoID string, remoteName string, remoteURL string) (*model.RemoteMutationResult, error) {
	logger.Log("Initiating remote addition request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return &model.RemoteMutationResult{Status: global.RemoteAddError}, nil
	}

	var addRemoteObject git.AddRemoteInterface
	addRemoteObject = git.AddRemoteStruct{
		Repo:       repo.GitRepo,
		RemoteName: remoteName,
		RemoteURL:  remoteURL,
	}
	return addRemoteObject.AddRemote(), nil
}

func (r *mutationResolver) DeleteRemote(ctx context.Context, repoID string, remoteName string) (*model.RemoteMutationResult, error) {
	logger.Log("Initiating remote deletion request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return &model.RemoteMutationResult{Status: global.RemoteAddError}, nil
	}

	var addRemoteObject git.DeleteRemoteInterface
	addRemoteObject = &git.DeleteRemoteStruct{
		Repo:       repo.GitRepo,
		RemoteName: remoteName,
	}
	return addRemoteObject.DeleteRemote(), nil
}

func (r *mutationResolver) EditRemote(ctx context.Context, repoID string, remoteName string, remoteURL string) (*model.RemoteMutationResult, error) {
	logger.Log("Initiating remote edit request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return &model.RemoteMutationResult{Status: global.RemoteEditError}, nil
	}

	var editRemoteObject git.RemoteEditInterface
	editRemoteObject = git.RemoteEditStruct{
		Repo:       repo.GitRepo,
		RemoteName: remoteName,
		RemoteUrl:  remoteURL,
	}
	return editRemoteObject.EditRemoteUrl(), nil
}

func (r *queryResolver) HealthCheck(ctx context.Context) (*model.HealthCheckParams, error) {
	logger.Log("Initiating health check request", global.StatusInfo)
	return api.HealthCheckApi(), nil
}

func (r *queryResolver) FetchRepo(ctx context.Context) (*model.FetchRepoParams, error) {
	logger.Log("Initiating fetch repo request", global.StatusInfo)
	return api.FetchRepo(), nil
}

func (r *queryResolver) GitRepoStatus(ctx context.Context, repoID string) (*model.GitRepoStatusResults, error) {
	logger.Log("Initiating repo status request", global.StatusInfo)
	return api.RepoStatus(repoID), nil
}

func (r *queryResolver) GitFolderContent(ctx context.Context, repoID string, directoryName *string) (*model.GitFolderContentResults, error) {
	logger.Log("Initiating get folder content request", global.StatusInfo)
	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return &model.GitFolderContentResults{
			TrackedFiles:     nil,
			FileBasedCommits: nil,
		}, nil
	}

	tmp := &model.GitFolderContentResults{}
	tmp.FileBasedCommits = nil
	tmp.TrackedFiles = nil

	var listFileObject git.ListFilesInterface
	listFileObject = git.ListFilesStruct{
		Repo:          repo.GitRepo,
		RepoPath:      repo.RepoPath,
		DirectoryName: *directoryName,
		FileName:      nil,
	}
	return listFileObject.ListFiles(), nil
}

func (r *queryResolver) GitCommitLogs(ctx context.Context, repoID string, referenceCommit string) (*model.GitCommitLogResults, error) {
	logger.Log("Initiating get commit logs request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	var commitLogObject git.CommitLogInterface
	commitLogObject = git.CommitLogStruct{
		Repo:            repo.GitRepo,
		ReferenceCommit: referenceCommit,
	}

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return &model.GitCommitLogResults{
			TotalCommits: nil,
			Commits:      nil,
		}, nil
	}
	return commitLogObject.CommitLogs(), nil
}

func (r *queryResolver) GitCommitFiles(ctx context.Context, repoID string, commitHash string) ([]*model.GitCommitFileResult, error) {
	logger.Log("Initiating get commit files request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	var commitFileListObject git.CommitFileListInterface
	commitFileListObject = git.CommitFileListStruct{
		Repo:       repo.GitRepo,
		CommitHash: commitHash,
	}
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return []*model.GitCommitFileResult{
			{
				Type:     "",
				FileName: "",
			},
		}, nil
	}
	return commitFileListObject.CommitFileList(), nil
}

func (r *queryResolver) SearchCommitLogs(ctx context.Context, repoID string, searchType string, searchKey string) ([]*model.GitCommits, error) {
	logger.Log("Initiating search commit logs request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	var searchCommitObject git.SearchCommitInterface

	searchCommitObject = git.SearchCommitStruct{
		Repo:       repo.GitRepo,
		SearchType: searchType,
		SearchKey:  searchKey,
	}

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return []*model.GitCommits{
			{
				Hash:             nil,
				Author:           nil,
				CommitTime:       nil,
				CommitMessage:    nil,
				CommitFilesCount: nil,
			},
		}, nil
	}
	return searchCommitObject.SearchCommitLogs(), nil
}

func (r *queryResolver) CodeFileDetails(ctx context.Context, repoID string, fileName string) (*model.CodeFileType, error) {
	logger.Log("Initiating code file view request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return &model.CodeFileType{
			FileData: nil,
		}, nil
	}
	var codeFileViewObject api.CodeViewInterface
	codeFileViewObject = api.CodeViewInputs{
		RepoPath: repo.RepoPath,
		FileName: fileName,
	}
	return codeFileViewObject.CodeFileView(), nil
}

func (r *queryResolver) GitChanges(ctx context.Context, repoID string) (*model.GitChangeResults, error) {
	logger.Log("Initiating repo changes request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	var gitChangeObject git.ChangedItemsInterface

	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)

		return &model.GitChangeResults{
			GitUntrackedFiles: nil,
			GitChangedFiles:   nil,
			GitStagedFiles:    nil,
		}, nil
	}

	gitChangeObject = git.ChangedItemStruct{
		Repo:     repo.GitRepo,
		RepoPath: repo.RepoPath,
	}
	return gitChangeObject.ChangedFiles(), nil
}

func (r *queryResolver) GitUnPushedCommits(ctx context.Context, repoID string, remoteURL string, remoteBranch string) ([]*model.GitCommits, error) {
	logger.Log("Initiating get un-pushed commits request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return nil, nil
	}
	var remoteDataObject git.RemoteDataInterface
	remoteDataObject = git.RemoteDataStruct{
		Repo:      repo.GitRepo,
		RemoteURL: remoteURL,
	}
	remoteName := remoteDataObject.GetRemoteName()
	remoteRef := remoteName + "/" + remoteBranch

	var unPushedObject git.UnPushedCommitInterface
	unPushedObject = git.UnPushedCommitStruct{
		Repo:      repo.GitRepo,
		RemoteRef: remoteRef,
	}
	return unPushedObject.UnPushedCommits(), nil
}

func (r *queryResolver) GitFileLineChanges(ctx context.Context, repoID string, fileName string) (*model.FileLineChangeResult, error) {
	logger.Log("Initiating get line by line diff request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return &model.FileLineChangeResult{
			DiffStat: "",
			FileDiff: nil,
		}, nil
	}

	var codeFileViewObject api.CodeViewInterface
	codeFileViewObject = api.CodeViewInputs{
		RepoPath: repo.RepoPath,
		FileName: fileName,
	}
	fileContent := codeFileViewObject.CodeFileView()

	var fileLineDiffObject git.FileLineDiffInterface
	fileLineDiffObject = git.FileLineDiffStruct{
		Repo:     repo.GitRepo,
		FileName: fileName,
		Data:     fileContent.FileData,
	}
	return fileLineDiffObject.FileLineDiff(), nil
}

func (r *queryResolver) SettingsData(ctx context.Context) (*model.SettingsDataResults, error) {
	logger.Log("Initiating get settings data request", global.StatusInfo)
	return api.GetSettingsData(), nil
}

func (r *queryResolver) CommitCompare(ctx context.Context, repoID string, baseCommit string, compareCommit string) ([]*model.GitCommitFileResult, error) {
	logger.Log("Initiating compare commit request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return []*model.GitCommitFileResult{
			{
				Type:     "",
				FileName: "",
			},
		}, nil
	}

	var compareCommitObject git.CompareCommitInterface
	compareCommitObject = git.CompareCommitStruct{
		Repo:                repo.GitRepo,
		BaseCommitString:    baseCommit,
		CompareCommitString: compareCommit,
	}

	return compareCommitObject.CompareCommit(), nil
}

func (r *queryResolver) BranchCompare(ctx context.Context, repoID string, baseBranch string, compareBranch string) ([]*model.BranchCompareResults, error) {
	logger.Log("Initiating compare branch request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return []*model.BranchCompareResults{
			{
				Date:    "",
				Commits: nil,
			},
		}, nil
	}

	var branchCompareObject git.BranchCompareInterface
	branchCompareObject = git.BranchCompareInputs{
		Repo:       repo.GitRepo,
		BaseBranch: baseBranch,
		DiffBranch: compareBranch,
	}
	return branchCompareObject.CompareBranch(), nil
}

func (r *queryResolver) GetRemote(ctx context.Context, repoID string) ([]*model.RemoteDetails, error) {
	logger.Log("Initating remote data fetching", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	var remoteObject git.RemoteDataInterface
	remoteObject = git.RemoteDataStruct{
		Repo: repo.GitRepo,
	}
	allRemoteData := remoteObject.GetAllRemotes()
	return allRemoteData, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var logger global.Logger
