package tests

import (
	"github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListFiles(t *testing.T) {
	lsFileChan := make(chan *git2.LsFileInfo)
	r, _ := git.OpenRepository(TestRepo)
	repoPath := r.Workdir()

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
