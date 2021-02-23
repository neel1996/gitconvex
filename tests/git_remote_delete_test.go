package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestDeleteRemoteStruct_DeleteRemote(t *testing.T) {
	var r *git.Repository
	cwd, _ := os.Getwd()
	var testRemoteName string

	if os.Getenv("GOTESTENV") == "ci" {
		r, _ = git.OpenRepository(path.Join(cwd, ".."))
		testRemoteName = "origin"
	} else {
		r, _ = git.OpenRepository(os.Getenv("REPODIR"))
		testRemoteName = "github"
	}

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
			if got := d.DeleteRemote(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteRemote() = %v, want %v", got, tt.want)
			}
		})
	}
}
