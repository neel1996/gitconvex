package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/neel1996/gitconvex/api"
	"github.com/neel1996/gitconvex/constants"
	"github.com/neel1996/gitconvex/controller"
	"github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/git/commit"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/git/remote"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/generated"
	"github.com/neel1996/gitconvex/graph/model"
	initialize "github.com/neel1996/gitconvex/init"
	"github.com/neel1996/gitconvex/utils"
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

	branchController := r.BranchController(ctx, middleware.NewRepository(repo.GitRepo))
	if branchController == nil {
		return global.BranchAddError, nil
	}

	return branchController.GitAddBranch(branchName, false, nil)
}

func (r *mutationResolver) CheckoutBranch(ctx context.Context, repoID string, branchName string) (string, error) {
	logger.Log("Initiating branch checkout request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)

	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)

	repo := <-repoChan

	branchController := r.BranchController(ctx, middleware.NewRepository(repo.GitRepo))
	if branchController == nil {
		return global.BranchCheckoutError, nil
	}

	return branchController.GitCheckoutBranch(branchName)
}

func (r *mutationResolver) DeleteBranch(ctx context.Context, repoID string, branchName string, forceFlag bool) (*model.BranchDeleteStatus, error) {
	logger.Log("Initiating branch deletion request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	branchController := r.BranchController(ctx, middleware.NewRepository(repo.GitRepo))
	if branchController == nil {
		return &model.BranchDeleteStatus{
			Status: global.BranchDeleteError,
		}, nil
	}

	return branchController.GitDeleteBranch(branchName)
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

	remoteName := initialize.RemoteObjects(ctx).RemoteName

	var fetchObject git.FetchInterface
	fetchObject = git.FetchStruct{
		Repo:         repo.GitRepo,
		RemoteName:   remoteName.GetRemoteNameWithUrl(),
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

func (r *mutationResolver) CommitChanges(ctx context.Context, repoID string, commitMessage []*string) (string, error) {
	logger.Log("Initiating commit changes request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	commitChanges := commit.NewCommitChanges(middleware.NewRepository(repo.GitRepo), utils.GenerateNonPointerArrayFrom(commitMessage))
	commitOperation := commit.Operation{Changes: commitChanges}

	return commitOperation.GitCommitChange()
}

func (r *mutationResolver) PushToRemote(ctx context.Context, repoID string, remoteHost string, branch string) (string, error) {
	logger.Log("Initiating push to remote request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	repo := <-repoChan

	remoteName := initialize.RemoteObjects(ctx).RemoteName.GetRemoteNameWithUrl()

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

	addRemote := initialize.RemoteObjects(ctx).AddRemote
	remoteObject := remote.Operation{Add: addRemote}

	return remoteObject.GitAddRemote()
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

	deleteRemote := initialize.RemoteObjects(ctx).DeleteRemote

	remoteObject := remote.Operation{
		Delete: deleteRemote,
	}

	return remoteObject.GitDeleteRemote()
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
		return &model.RemoteMutationResult{Status: constants.RemoteEditError.Error()}, nil
	}

	editRemote := initialize.RemoteObjects(ctx).EditRemote
	editRemoteObject := remote.Operation{
		Edit: editRemote,
	}

	return editRemoteObject.GitEditRemote()
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

	return controller.NewCommitLogController(repoID, referenceCommit).GetCommitLogs()
}

func (r *queryResolver) GitCommitFiles(ctx context.Context, repoID string, commitHash string) ([]*model.GitCommitFileResult, error) {
	logger.Log("Initiating get commit files request", global.StatusInfo)

	return controller.NewCommitFileHistoryController(repoID, commitHash).GetCommitFileHistory()
}

func (r *queryResolver) SearchCommitLogs(ctx context.Context, repoID string, searchType string, searchKey string) ([]*model.GitCommits, error) {
	logger.Log("Initiating search commit logs request", global.StatusInfo)

	return controller.NewCommitLogSearchController(repoID, searchType, searchKey).GetMatchingCommits()
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

func (r *queryResolver) GitUnPushedCommits(ctx context.Context, repoID string, remoteURL string, remoteBranch string) (*model.UnPushedCommitResult, error) {
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

	remoteName := initialize.RemoteObjects(ctx).RemoteName.GetRemoteNameWithUrl()

	remoteRef := remoteName + "/" + remoteBranch

	var unPushedObject git.UnPushedCommitInterface
	unPushedObject = git.UnPushedCommitStruct{
		Repo:      repo.GitRepo,
		RemoteRef: remoteRef,
		LocalRef:  remoteBranch,
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

	branchController := r.BranchController(ctx, middleware.NewRepository(repo.GitRepo))
	if branchController == nil {
		return []*model.BranchCompareResults{
			{
				Date:    "",
				Commits: nil,
			},
		}, nil
	}

	return branchController.GitCompareBranches(baseBranch, compareBranch)

}

func (r *queryResolver) GetRemote(ctx context.Context, repoID string) ([]*model.RemoteDetails, error) {
	logger.Log("Initating remote data fetching", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	var repoObject git.RepoInterface
	repoObject = git.RepoStruct{RepoId: repoID}
	go repoObject.Repo(repoChan)
	//repo := <-repoChan

	listRemote := initialize.RemoteObjects(ctx).ListRemote
	remoteUrlList := remote.Operation{List: listRemote}

	return remoteUrlList.GitGetAllRemote()
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
