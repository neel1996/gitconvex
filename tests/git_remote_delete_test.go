package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteRemoteStruct_DeleteRemote(t *testing.T) {
	r, _ := git.OpenRepository(TestRepo)

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
