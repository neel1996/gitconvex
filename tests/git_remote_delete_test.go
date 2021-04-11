package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestDeleteRemoteStruct_DeleteRemote(t *testing.T) {
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

	testRemoteName := "test_remote"
	var testObj git2.AddRemoteInterface
	testObj = git2.AddRemoteStruct{
		Repo:       r,
		RemoteName: testRemoteName,
		RemoteURL:  "git@github.com:neel1996/gitconvex-server.git",
	}
	_ = testObj.AddRemote()

	type fields struct {
		Repo       git.Repository
		RemoteName string
	}
	tests := []struct {
		name   string
		fields fields
		want   *model.RemoteMutationResult
	}{
		{name: "Test for remote deletion", fields: struct {
			Repo       git.Repository
			RemoteName string
		}{Repo: *r, RemoteName: testRemoteName}, want: &model.RemoteMutationResult{Status: global.RemoteDeleteSuccess}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &git2.DeleteRemoteStruct{
				Repo:       &tt.fields.Repo,
				RemoteName: tt.fields.RemoteName,
			}
			got := d.DeleteRemote()

			assert.Equal(t, tt.want.Status, got.Status)
		})
	}
}
