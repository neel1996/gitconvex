package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/graph/model"
	"os"
	"testing"
)

func TestSearchCommitLogs(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex/starfleet"
		r, _ = git.PlainOpen(repoPath)
	}

	type args struct {
		repo       *git.Repository
		searchType string
		searchKey  string
	}

	hash := "46aa56e78f2a26d23f604f8e9bbdc240a0a5dbbe"
	author := "Neel"

	expectedResult := &model.GitCommits{
		Hash:   &hash,
		Author: &author,
	}

	tests := []struct {
		name string
		args args
		want []*model.GitCommits
	}{
		{name: "Git commit log search test case", args: struct {
			repo       *git.Repository
			searchType string
			searchKey  string
		}{repo: r, searchType: "hash", searchKey: "46aa56e"}, want: []*model.GitCommits{expectedResult}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := git2.SearchCommitLogs(tt.args.repo, tt.args.searchType, tt.args.searchKey); *got[0].Hash != *tt.want[0].Hash {
				t.Errorf("SearchCommitLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}
