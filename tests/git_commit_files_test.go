package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/graph/model"
	"os"
	"reflect"
	"testing"
)

func TestCommitFileList(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex-server/starfleet"
		r, _ = git.PlainOpen(repoPath)
	}

	type args struct {
		repo       *git.Repository
		commitHash string
	}
	tests := []struct {
		name string
		args args
		want []*model.GitCommitFileResult
	}{
		{name: "Git commit file list test case", args: struct {
			repo       *git.Repository
			commitHash string
		}{repo: r, commitHash: "46aa56e78f2a26d23f604f8e9bbdc240a0a5dbbe"}, want: []*model.GitCommitFileResult{&model.GitCommitFileResult{
			Type:     "A",
			FileName: "codeql-analysis.yml",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := git2.CommitFileList(tt.args.repo, tt.args.commitHash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommitFileList() = %v, want %v", got, tt.want)
			}
		})
	}
}
