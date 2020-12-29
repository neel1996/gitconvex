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

func (r *mutationResolver) AddRepo(ctx context.Context, repoName string, repoPath string, cloneSwitch bool, repoURL *string, initSwitch bool, authOption string, userName *string, password *string) (*model.AddRepoParams, error) {
	return api.AddRepo(model.NewRepoInputs{
		RepoName:    repoName,
		RepoPath:    repoPath,
		CloneSwitch: cloneSwitch,
		RepoURL:     repoURL,
		InitSwitch:  initSwitch,
		AuthOption:  authOption,
		UserName:    userName,
		Password:    password,
	}), nil
}

func (r *mutationResolver) AddBranch(ctx context.Context, repoID string, branchName string) (string, error) {
	logger.Log("Initiating branch addition request", global.StatusInfo)
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid. Branch addition failed", global.StatusError)
		return global.BranchAddError, nil
	}

	return git.AddBranch(repo.GitRepo, branchName), nil
}

func (r *mutationResolver) CheckoutBranch(ctx context.Context, repoID string, branchName string) (string, error) {
	logger.Log("Initiating branch checkout request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repository is invalid or HEAD is null", global.StatusError)
		return global.BranchCheckoutError, nil
	}

	return git.CheckoutBranch(repo.GitRepo, branchName), nil
}

func (r *mutationResolver) DeleteBranch(ctx context.Context, repoID string, branchName string, forceFlag bool) (*model.BranchDeleteStatus, error) {
	logger.Log("Initiating branch deletion request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is null", global.StatusError)
		return &model.BranchDeleteStatus{
			Status: global.BranchDeleteError,
		}, nil
	}
	return git.DeleteBranch(repo.GitRepo, branchName, forceFlag), nil
}

func (r *mutationResolver) AddRemote(ctx context.Context, repoID string, remoteName string, remoteURL string) (string, error) {
	logger.Log("Initiating remote addition request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return global.RemoteAddError, nil
	}
	return git.AddRemote(repo.GitRepo, remoteName, remoteURL), nil
}

func (r *mutationResolver) FetchFromRemote(ctx context.Context, repoID string, remoteURL *string, remoteBranch *string) (*model.FetchResult, error) {
	logger.Log("Initiating fetch from remote request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)

		return &model.FetchResult{
			Status:       global.FetchFromRemoteError,
			FetchedItems: nil,
		}, nil
	}
	return git.FetchFromRemote(repo.GitRepo, *remoteURL, *remoteBranch), nil
}

func (r *mutationResolver) PullFromRemote(ctx context.Context, repoID string, remoteURL *string, remoteBranch *string) (*model.PullResult, error) {
	logger.Log("Initiating pull from remote request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return &model.PullResult{
			Status:      global.PullFromRemoteError,
			PulledItems: nil,
		}, nil
	}

	return git.PullFromRemote(repo.GitRepo, *remoteURL, *remoteBranch), nil
}

func (r *mutationResolver) StageItem(ctx context.Context, repoID string, item string) (string, error) {
	logger.Log("Initiating stage item request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return global.StageItemError, nil
	}

	return git.StageItem(repo.GitRepo, item), nil
}

func (r *mutationResolver) RemoveStagedItem(ctx context.Context, repoID string, item string) (string, error) {
	logger.Log("Initiating remove item request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return global.RemoveItemError, nil
	}
	return git.RemoveItem(repo.RepoPath, item), nil
}

func (r *mutationResolver) RemoveAllStagedItem(ctx context.Context, repoID string) (string, error) {
	logger.Log("Initiating remove all items request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return global.RemoveAllItemsError, nil
	}
	return git.ResetAllItems(repo.GitRepo), nil
}

func (r *mutationResolver) StageAllItems(ctx context.Context, repoID string) (string, error) {
	logger.Log("Initiating stage all items request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return global.StageAllItemsError, nil
	}
	return git.StageAllItems(repo.GitRepo), nil
}

func (r *mutationResolver) CommitChanges(ctx context.Context, repoID string, commitMessage string) (string, error) {
	logger.Log("Initiating commit changes request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		w, _ := repo.GitRepo.Worktree()
		if w != nil {
			return git.CommitChanges(repo.GitRepo, commitMessage), nil
		} else {
			logger.Log("Repo is invalid or worktree is null", global.StatusError)
			return global.CommitChangeError, nil
		}
	}
	return git.CommitChanges(repo.GitRepo, commitMessage), nil
}

func (r *mutationResolver) PushToRemote(ctx context.Context, repoID string, remoteHost string, branch string) (string, error) {
	logger.Log("Initiating push to remote request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	remoteName := git.GetRemoteName(repo.GitRepo, remoteHost)

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil || remoteName == "" {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return global.PushToRemoteError, nil
	}

	return git.PushToRemote(repo.GitRepo, remoteName, branch), nil
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
	go git.Repo(repoID, repoChan)
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

	return git.ListFiles(repo.GitRepo, repo.RepoPath, *directoryName), nil
}

func (r *queryResolver) GitCommitLogs(ctx context.Context, repoID string, referenceCommit string) (*model.GitCommitLogResults, error) {
	logger.Log("Initiating get commit logs request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return &model.GitCommitLogResults{
			TotalCommits: nil,
			Commits:      nil,
		}, nil
	}
	return git.CommitLogs(repo.GitRepo, referenceCommit), nil
}

func (r *queryResolver) GitCommitFiles(ctx context.Context, repoID string, commitHash string) ([]*model.GitCommitFileResult, error) {
	logger.Log("Initiating get commit files request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
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
	return git.CommitFileList(repo.GitRepo, commitHash), nil
}

func (r *queryResolver) SearchCommitLogs(ctx context.Context, repoID string, searchType string, searchKey string) ([]*model.GitCommits, error) {
	logger.Log("Initiating search commit logs request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
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
	return git.SearchCommitLogs(repo.GitRepo, searchType, searchKey), nil
}

func (r *queryResolver) CodeFileDetails(ctx context.Context, repoID string, fileName string) (*model.CodeFileType, error) {
	logger.Log("Initiating code file view request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return &model.CodeFileType{
			FileData: nil,
		}, nil
	}
	return api.CodeFileView(repo.RepoPath, fileName), nil
}

func (r *queryResolver) GitChanges(ctx context.Context, repoID string) (*model.GitChangeResults, error) {
	logger.Log("Initiating repo changes request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		logger.Log("Repo is invalid", global.StatusError)
		return &model.GitChangeResults{
			GitUntrackedFiles: nil,
			GitChangedFiles:   nil,
			GitStagedFiles:    nil,
		}, nil
	}
	return git.ChangedFiles(repo.GitRepo), nil
}

func (r *queryResolver) GitUnPushedCommits(ctx context.Context, repoID string, remoteURL string, remoteBranch string) ([]*string, error) {
	logger.Log("Initiating get un-pushed commits request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return nil, nil
	}

	remoteName := git.GetRemoteName(repo.GitRepo, remoteURL)
	remoteRef := remoteName + "/" + remoteBranch

	return git.UnPushedCommits(repo.GitRepo, remoteRef), nil
}

func (r *queryResolver) GitFileLineChanges(ctx context.Context, repoID string, fileName string) (*model.FileLineChangeResult, error) {
	logger.Log("Initiating get line by line diff request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		logger.Log("Repo is invalid or HEAD is nil", global.StatusError)
		return &model.FileLineChangeResult{
			DiffStat: "",
			FileDiff: nil,
		}, nil
	}

	fileContent := api.CodeFileView(repo.RepoPath, fileName)
	return git.FileLineDiff(repo.GitRepo, fileName, fileContent.FileData), nil
}

func (r *queryResolver) SettingsData(ctx context.Context) (*model.SettingsDataResults, error) {
	logger.Log("Initiating get settings data request", global.StatusInfo)
	return api.GetSettingsData(), nil
}

func (r *queryResolver) CommitCompare(ctx context.Context, repoID string, baseCommit string, compareCommit string) ([]*model.GitCommitFileResult, error) {
	logger.Log("Initiating compare commit request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
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

	return git.CompareCommit(repo.GitRepo, baseCommit, compareCommit), nil
}

func (r *queryResolver) BranchCompare(ctx context.Context, repoID string, baseBranch string, compareBranch string) ([]*model.BranchCompareResults, error) {
	logger.Log("Initiating compare branch request", global.StatusInfo)

	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
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
	return git.CompareBranch(repo.GitRepo, baseBranch, compareBranch), nil
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
