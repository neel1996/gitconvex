package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestSearchCommitLogs(t *testing.T) {
	var repoPath string
	var r *git.Repository
	cwd, _ := os.Getwd()
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)
	mockRepoPath := path.Join(cwd, "../..") + "/starfleet"

	if currentEnv == "ci" {
		repoPath = mockRepoPath
		r, _ = git.OpenRepository(repoPath)
	} else {
		repoPath = path.Join(cwd, "../..")
		r, _ = git.OpenRepository(repoPath)
	}

	sampleCommits := git2.CommitLogStruct{
		Repo:            r,
		ReferenceCommit: "",
	}.CommitLogs()
	hash := *sampleCommits.Commits[0].Hash

	type args struct {
		repo       *git.Repository
		searchType string
		searchKey  string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Git commit log search test case", args: struct {
			repo       *git.Repository
			searchType string
			searchKey  string
		}{repo: r, searchType: "hash", searchKey: hash}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.SearchCommitInterface
			testObj = git2.SearchCommitStruct{
				Repo:       tt.args.repo,
				SearchType: tt.args.searchType,
				SearchKey:  tt.args.searchKey,
			}

			got := testObj.SearchCommitLogs()

			assert.NotZero(t, len(got))
		})
	}
}
