package tests

import (
	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/remote"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteRemoteStruct_DeleteRemote(t *testing.T) {
	r, _ := git.OpenRepository(TestRepo)

	testRemoteName := "test_remote"
	addRemote := remote.NewAddRemote(r, testRemoteName, "git@github.com:neel1996/gitconvex-server.git")
	_ = addRemote.NewRemote()

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
			remoteObj := remote.NewDeleteRemoteInterface(&tt.fields.Repo, tt.fields.RemoteName)
			got := remoteObj.DeleteRemote()

			assert.Equal(t, tt.want.Status, got.Status)
		})
	}
}
