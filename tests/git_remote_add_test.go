package tests

import (
	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/remote"
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
			var testObj remote.Add

			testObj = remote.NewAddRemote(r, tt.args.remoteName, tt.args.remoteURL)
			got := testObj.NewRemote()

			assert.Equal(t, tt.want, got.Status)
			remote.NewDeleteRemoteInterface(r, tt.args.remoteName).DeleteRemote()
		})
	}
}
