package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
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
		r, _ = git.OpenRepository(repoPath)
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
			FileName: ".github/workflows/codeql-analysis.yml",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObject git2.CommitFileListInterface
			testObject = git2.CommitFileListStruct{
				Repo:       tt.args.repo,
				CommitHash: tt.args.commitHash,
			}

			if got := testObject.CommitFileList(); !reflect.DeepEqual(got, tt.want) {
				for _, fileItem := range got {
					fmt.Println(fileItem.FileName)
					if fileItem.FileName != tt.want[0].FileName {
						t.Errorf("CommitFileList() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}
