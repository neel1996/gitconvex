package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/graph/model"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestChangedFiles(t *testing.T) {
	var repoPath string
	var r *git.Repository
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex-server/starfleet"
		r, _ = git.PlainOpen(repoPath)
	}

	untrackedResult := "untracked.txt"
	changedResult := "README.md"
	stagedResult := "README.md"

	uErr := ioutil.WriteFile(repoPath+"/"+untrackedResult, []byte{byte(63)}, 0755)
	cErr := ioutil.WriteFile(repoPath+"/"+changedResult, []byte{byte(83)}, 0755)

	var stageObject git2.StageItemInterface
	stageObject = git2.StageItemStruct{
		Repo:     r,
		FileItem: repoPath + "/" + changedResult,
	}
	stageObject.StageItem()

	sErr := ioutil.WriteFile(repoPath+"/"+changedResult, []byte{byte(70)}, 0755)
	fmt.Println(uErr, cErr, sErr)

	expectedResults := &model.GitChangeResults{
		GitUntrackedFiles: []*string{&untrackedResult},
		GitChangedFiles:   []*string{&changedResult},
		GitStagedFiles:    []*string{&stagedResult},
	}

	type args struct {
		repo *git.Repository
	}
	tests := []struct {
		name string
		args args
		want *model.GitChangeResults
	}{
		{name: "Git changed files test case", args: struct{ repo *git.Repository }{repo: r}, want: expectedResults},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.ChangeInterface
			testObj = git2.ChangedStruct{
				Repo:     tt.args.repo,
				RepoPath: "",
			}
			got := testObj.ChangedFiles()

			stagedFile := *got.GitStagedFiles[0]
			untrackedFile := *got.GitUntrackedFiles[0]
			changedFile := *got.GitChangedFiles[0]
			changedFile = strings.Split(changedFile, ",")[1]

			fmt.Println(stagedFile)
			fmt.Println(untrackedFile)
			fmt.Println(changedFile)

			if stagedFile == *tt.want.GitStagedFiles[0] && untrackedFile == *tt.want.GitUntrackedFiles[0] && changedFile == *tt.want.GitChangedFiles[0] {
				fmt.Println("Test Passed")
			} else {
				t.Errorf("ChangedFiles() = %v, want %v", *got.GitStagedFiles[0], *tt.want.GitStagedFiles[0])
			}
		})
	}
}
