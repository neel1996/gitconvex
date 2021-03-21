package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"io/ioutil"
	"os"
	"testing"
)

func TestResetAllItems(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex-server/starfleet"
		r, _ = git.OpenRepository(repoPath)
	}

	untrackedResult := "untracked.txt"
	_ = ioutil.WriteFile(untrackedResult, []byte{byte(63)}, 0755)

	var stageObject git2.StageItemInterface
	stageObject = git2.StageItemStruct{
		Repo:     r,
		FileItem: untrackedResult,
	}
	stageObject.StageItem()

	type args struct {
		repo *git.Repository
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Git rest all test case", args: struct{ repo *git.Repository }{repo: r}, want: "STAGE_ALL_REMOVE_SUCCESS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.ResetAllInterface
			testObj = git2.ResetAllStruct{Repo: tt.args.repo}
			if got := testObj.ResetAllItems(); got != tt.want {
				t.Errorf("ResetAllItems() = %v, want %v", got, tt.want)
			}
		})
	}
}
