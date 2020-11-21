package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/neel1996/gitconvex-server/api"
	"github.com/neel1996/gitconvex-server/git"
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
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		return "BRANCH_ADD_FAILED", nil
	}

	return git.AddBranch(repo.GitRepo, branchName), nil
}

func (r *mutationResolver) CheckoutBranch(ctx context.Context, repoID string, branchName string) (string, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		return "Failed to checkout branch", nil
	}

	return git.CheckoutBranch(repo.GitRepo, branchName), nil
}

func (r *mutationResolver) DeleteBranch(ctx context.Context, repoID string, branchName string, forceFlag bool) (*model.BranchDeleteStatus, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		return &model.BranchDeleteStatus{
			Status: "BRANCH_DELETE_FAILED",
		}, nil
	}
	return git.DeleteBranch(repo.GitRepo, branchName, forceFlag), nil
}

func (r *mutationResolver) AddRemote(ctx context.Context, repoID string, remoteName string, remoteURL string) (string, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		return "REMOTE_ADD_FAILED", nil
	}
	return git.AddRemote(repo.GitRepo, remoteName, remoteURL), nil
}

func (r *mutationResolver) FetchFromRemote(ctx context.Context, repoID string, remoteURL *string, remoteBranch *string) (*model.FetchResult, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		return &model.FetchResult{
			Status:       "FETCH ERROR",
			FetchedItems: nil,
		}, nil
	}
	return git.FetchFromRemote(repo.GitRepo, *remoteURL, *remoteBranch), nil
}

func (r *mutationResolver) PullFromRemote(ctx context.Context, repoID string, remoteURL *string, remoteBranch *string) (*model.PullResult, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		return &model.PullResult{
			Status:      "PULL ERROR",
			PulledItems: nil,
		}, nil
	}

	return git.PullFromRemote(repo.GitRepo, *remoteURL, *remoteBranch), nil
}

func (r *mutationResolver) StageItem(ctx context.Context, repoID string, item string) (string, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		return "ADD_ITEM_FAILED", nil
	}

	return git.StageItem(repo.GitRepo, item), nil
}

func (r *mutationResolver) RemoveStagedItem(ctx context.Context, repoID string, item string) (string, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		return "STAGE_REMOVE_FAILED", nil
	}
	return git.RemoveItem(repo.RepoPath, item), nil
}

func (r *mutationResolver) RemoveAllStagedItem(ctx context.Context, repoID string) (string, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		return "STAGE_ALL_REMOVE_FAILED", nil
	}
	return git.ResetAllItems(repo.GitRepo), nil
}

func (r *mutationResolver) StageAllItems(ctx context.Context, repoID string) (string, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		return "ALL_STAGE_FAILED", nil
	}
	return git.StageAllItems(repo.GitRepo), nil
}

func (r *mutationResolver) CommitChanges(ctx context.Context, repoID string, commitMessage string) (string, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		w, _ := repo.GitRepo.Worktree()
		if w != nil {
			return git.CommitChanges(repo.GitRepo, commitMessage), nil
		} else {
			return "COMMIT_FAILED", nil
		}
	}
	return git.CommitChanges(repo.GitRepo, commitMessage), nil
}

func (r *mutationResolver) PushToRemote(ctx context.Context, repoID string, remoteHost string, branch string) (string, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	remoteName := git.GetRemoteName(repo.GitRepo, remoteHost)

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil || remoteName == "" {
		return "PUSH_FAILED", nil
	}

	return git.PushToRemote(repo.GitRepo, remoteName, branch), nil
}

func (r *mutationResolver) SettingsEditPort(ctx context.Context, newPort string) (string, error) {
	return api.UpdatePortNumber(newPort), nil
}

func (r *mutationResolver) UpdateRepoDataFile(ctx context.Context, newDbFile string) (string, error) {
	return api.UpdateDBFilePath(newDbFile), nil
}

func (r *mutationResolver) DeleteRepo(ctx context.Context, repoID string) (*model.DeleteStatus, error) {
	return api.DeleteRepo(repoID), nil
}

func (r *queryResolver) HealthCheck(ctx context.Context) (*model.HealthCheckParams, error) {
	return api.HealthCheckApi(), nil
}

func (r *queryResolver) FetchRepo(ctx context.Context) (*model.FetchRepoParams, error) {
	return api.FetchRepo(), nil
}

func (r *queryResolver) GitRepoStatus(ctx context.Context, repoID string) (*model.GitRepoStatusResults, error) {
	return api.RepoStatus(repoID), nil
}

func (r *queryResolver) GitFolderContent(ctx context.Context, repoID string, directoryName *string) (*model.GitFolderContentResults, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
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
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		return &model.GitCommitLogResults{
			TotalCommits: nil,
			Commits:      nil,
		}, nil
	}
	return git.CommitLogs(repo.GitRepo, referenceCommit), nil
}

func (r *queryResolver) GitCommitFiles(ctx context.Context, repoID string, commitHash string) ([]*model.GitCommitFileResult, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
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
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		return []*model.GitCommits{
			{
				Hash:               nil,
				Author:             nil,
				CommitTime:         nil,
				CommitMessage:      nil,
				CommitRelativeTime: nil,
				CommitFilesCount:   nil,
			},
		}, nil
	}
	return git.SearchCommitLogs(repo.GitRepo, searchType, searchKey), nil
}

func (r *queryResolver) CodeFileDetails(ctx context.Context, repoID string, fileName string) (*model.CodeFileType, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		return &model.CodeFileType{
			FileData: nil,
		}, nil
	}
	return api.CodeFileView(repo.RepoPath, fileName), nil
}

func (r *queryResolver) GitChanges(ctx context.Context, repoID string) (*model.GitChangeResults, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if repo.GitRepo == nil {
		return &model.GitChangeResults{
			GitUntrackedFiles: nil,
			GitChangedFiles:   nil,
			GitStagedFiles:    nil,
		}, nil
	}
	return git.ChangedFiles(repo.GitRepo), nil
}

func (r *queryResolver) GitUnPushedCommits(ctx context.Context, repoID string, remoteURL string, remoteBranch string) ([]*string, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		return nil, nil
	}

	remoteName := git.GetRemoteName(repo.GitRepo, remoteURL)
	remoteRef := remoteName + "/" + remoteBranch

	return git.UnPushedCommits(repo.GitRepo, remoteRef), nil
}

func (r *queryResolver) GitFileLineChanges(ctx context.Context, repoID string, fileName string) (*model.FileLineChangeResult, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan

	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
		return &model.FileLineChangeResult{
			DiffStat: "",
			FileDiff: nil,
		}, nil
	}

	fileContent := api.CodeFileView(repo.RepoPath, fileName)
	return git.FileLineDiff(repo.GitRepo, fileName, fileContent.FileData), nil
}

func (r *queryResolver) SettingsData(ctx context.Context) (*model.SettingsDataResults, error) {
	return api.GetSettingsData(), nil
}

func (r *queryResolver) CommitCompare(ctx context.Context, repoID string, baseCommit string, compareCommit string) ([]*model.GitCommitFileResult, error) {
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
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
	repoChan := make(chan git.RepoDetails)
	go git.Repo(repoID, repoChan)
	repo := <-repoChan
	if head, _ := repo.GitRepo.Head(); repo.GitRepo == nil || head == nil {
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
