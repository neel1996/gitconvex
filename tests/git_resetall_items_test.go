package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestResetAllItems(t *testing.T) {
	var repoPath string
	var r *git.Repository

	cwd, _ := os.Getwd()
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = path.Join(cwd, "..")
		r, _ = git.OpenRepository(repoPath)
	} else {
		repoPath = path.Join(cwd, "../..")
		r, _ = git.OpenRepository(repoPath)
	}

	mockFile := "untracked.txt"
	_ = ioutil.WriteFile(mockFile, []byte{byte(63)}, 0755)

	var stageObject git2.StageItemInterface
	stageObject = git2.StageItemStruct{
		Repo:     r,
		FileItem: mockFile,
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
			got := testObj.ResetAllItems()

			assert.Equal(t, tt.want, got)

			_ = os.Remove(mockFile)
		})
	}
}
