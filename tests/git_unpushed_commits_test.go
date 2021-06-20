package tests

import (
	"github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestUnPushedCommits(t *testing.T) {
	r, _ := git.OpenRepository(TestRepo)
	repoPath := r.Workdir()

	untrackedResult := "untracked.txt"

	_ = ioutil.WriteFile(repoPath+untrackedResult, []byte{byte(63)}, 0755)

	var stageObject git2.StageItemInterface
	var commitObject git2.CommitInterface

	stageObject = git2.StageItemStruct{
		Repo:     r,
		FileItem: untrackedResult,
	}

	commitObject = git2.CommitStruct{
		Repo:          r,
		CommitMessage: "Test Commit",
	}

	stageObject.StageItem()
	commitObject.CommitChanges()

	type args struct {
		repo      *git.Repository
		remoteRef string
	}
	tests := []struct {
		name string
		args args
		want []*string
	}{
		{name: "Git unpushed commit list test case", args: struct {
			repo      *git.Repository
			remoteRef string
		}{repo: r, remoteRef: "origin/master"}, want: []*string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.UnPushedCommitInterface
			testObj = git2.UnPushedCommitStruct{
				Repo:      tt.args.repo,
				RemoteRef: tt.args.remoteRef,
			}

			got := testObj.UnPushedCommits()

			assert.NotNil(t, got)
			assert.NotEmpty(t, *got.GitCommits[0].CommitMessage)
		})
	}
}
