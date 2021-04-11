package tests

import (
	"fmt"
	"github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestListFiles(t *testing.T) {
	lsFileChan := make(chan *git2.LsFileInfo)
	var repoPath string
	var r *git.Repository
	cwd, _ := os.Getwd()
	mockRepoPath := path.Join(cwd, "../..") + "/starfleet"
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = mockRepoPath
		r, _ = git.OpenRepository(repoPath)
	} else {
		repoPath = path.Join(cwd, "../..")
		r, _ = git.OpenRepository(repoPath)
	}

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

			assert.NotZerof(t, len(trackedFiles), "Repo has no files")
			assert.NotZerof(t, len(commits), "File based commit list is empty")
		})
	}
}
