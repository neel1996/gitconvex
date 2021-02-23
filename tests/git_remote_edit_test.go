package tests

import (
	"github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/global"
	"os"
	"path"
	"testing"
)

func TestRemoteEditStruct_EditRemoteUrl(t *testing.T) {
	var r *git.Repository
	cwd, _ := os.Getwd()

	if os.Getenv("GOTESTENV") == "ci" {
		r, _ = git.OpenRepository(path.Join(cwd, ".."))
	} else {
		r, _ = git.OpenRepository(os.Getenv("REPODIR"))
	}

	type fields struct {
		Repo       *git.Repository
		RemoteName string
		RemoteUrl  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Test edit remote api", fields: struct {
			Repo       *git.Repository
			RemoteName string
			RemoteUrl  string
		}{Repo: r, RemoteName: "origin", RemoteUrl: "git@github.com:neel1996/gitconvex-server.git"}, want: global.RemoteEditSuccess},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &git2.RemoteEditStruct{
				Repo:       tt.fields.Repo,
				RemoteName: tt.fields.RemoteName,
				RemoteUrl:  tt.fields.RemoteUrl,
			}
			if got := e.EditRemoteUrl(); got.Status != tt.want {
				t.Errorf("EditRemote() = %v, want %v", got, tt.want)
			}
		})
	}
}
