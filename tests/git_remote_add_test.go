package tests

import (
	git "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex/git"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddRemote(t *testing.T) {
	r, _ := git.OpenRepository(TestRepo)

	type args struct {
		repo       *git.Repository
		remoteName string
		remoteURL  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Git remote add test case", args: struct {
			repo       *git.Repository
			remoteName string
			remoteURL  string
		}{repo: r, remoteName: "test_remote", remoteURL: "https://github.com/neel1996/starfleet.git"}, want: "REMOTE_ADD_SUCCESS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.AddRemoteInterface
			testObj = git2.AddRemoteStruct{
				Repo:       tt.args.repo,
				RemoteName: tt.args.remoteName,
				RemoteURL:  tt.args.remoteURL,
			}
			got := testObj.AddRemote()

			assert.Equal(t, tt.want, got.Status)

			obj := git2.DeleteRemoteStruct{
				Repo:       r,
				RemoteName: tt.args.remoteName,
			}
			obj.DeleteRemote()
		})
	}
}
