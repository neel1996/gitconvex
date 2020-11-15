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

func TestCompareCommit(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex-server/starfleet"
		r, _ = git.PlainOpen(repoPath)
	}

	type args struct {
		repo                *git.Repository
		baseCommitString    string
		compareCommitString string
	}
	tests := []struct {
		name string
		args args
		want []*model.GitCommitFileResult
	}{
		{name: "Git commit compare test case", args: struct {
			repo                *git.Repository
			baseCommitString    string
			compareCommitString string
		}{repo: r, baseCommitString: "46aa56e78f2a26d23f604f8e9bbdc240a0a5dbbe", compareCommitString: "bc87f72ab7206afa091d648fc5a001761b6b890c"}, want: []*model.GitCommitFileResult{&model.GitCommitFileResult{
			Type:     "D",
			FileName: ".github/workflows/codeql-analysis.yml",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := git2.CompareCommit(tt.args.repo, tt.args.baseCommitString, tt.args.compareCommitString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompareCommit() = %v, want %v", got, tt.want)
			}
		})
	}
}
