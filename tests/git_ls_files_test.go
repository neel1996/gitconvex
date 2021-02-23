package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	assert2 "github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestListFiles(t *testing.T) {
	assert := assert2.New(t)
	lsFileChan := make(chan *git2.LsFileInfo)
	var repoPath string

	cwd, _ := os.Getwd()
	repoPath = path.Join(cwd, "..")

	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = "/home/runner/work/gitconvex-server/starfleet"
	}
	r, _ := git.PlainOpen(repoPath)

	type args struct {
		repo       *git.Repository
		repoPath   string
		lsFileChan chan *git2.LsFileInfo
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Git ls files test case", args: struct {
			repo       *git.Repository
			repoPath   string
			lsFileChan chan *git2.LsFileInfo
		}{repo: r, repoPath: repoPath, lsFileChan: lsFileChan}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.ListFilesInterface
			testObj = git2.ListFilesStruct{
				Repo:          tt.args.repo,
				RepoPath:      tt.args.repoPath,
				DirectoryName: "",
			}
			repoContent := testObj.ListFiles()

			trackedFiles := repoContent.TrackedFiles
			commits := repoContent.FileBasedCommits

			assert.Greater(len(trackedFiles), 0, "Repo has no files")
			assert.Greater(len(commits), 0, "File based commit list is empty")
		})
	}
}
