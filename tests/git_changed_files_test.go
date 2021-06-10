package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestChangedFiles(t *testing.T) {
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

	mockFile := "mockFile.txt"
	_ = ioutil.WriteFile(repoPath+"/"+mockFile, []byte{byte(63)}, 0755)
	var stageObject git2.StageItemInterface
	stageObject = git2.StageItemStruct{
		Repo:     r,
		FileItem: mockFile,
	}

	stageObject.StageItem()
	_ = ioutil.WriteFile(repoPath+"/"+mockFile, []byte{byte(83)}, 0755)

	untrackedMockFile := "mockFileTwo.txt"
	_ = ioutil.WriteFile(repoPath+"/"+untrackedMockFile, []byte{byte(63)}, 0755)

	expectedResults := &model.GitChangeResults{
		GitUntrackedFiles: []*string{&untrackedMockFile},
		GitChangedFiles:   []*string{&mockFile},
		GitStagedFiles:    []*string{&mockFile},
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
			var testObj git2.ChangedItemsInterface
			testObj = git2.ChangedItemStruct{
				Repo:     tt.args.repo,
				RepoPath: "",
			}

			got := testObj.ChangedFiles()

			stagedFile := *got.GitStagedFiles[0]
			untrackedFile := *got.GitUntrackedFiles[0]
			changedFile := *got.GitChangedFiles[0]
			changedFile = strings.Split(changedFile, ",")[1]

			assert.Equal(t, *tt.want.GitStagedFiles[0], stagedFile)
			assert.Equal(t, *tt.want.GitChangedFiles[0], changedFile)
			assert.Equal(t, *tt.want.GitUntrackedFiles[0], untrackedFile)

			reset := git2.ResetAllStruct{
				Repo: r,
			}
			reset.ResetAllItems()

			fmt.Println(os.Remove(repoPath + "/" + mockFile))
			fmt.Println(os.Remove(repoPath + "/" + untrackedMockFile))
		})
	}
}
