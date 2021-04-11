package tests

import (
	"fmt"
	"github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestRemoteEditStruct_EditRemoteUrl(t *testing.T) {
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
			got := e.EditRemoteUrl()

			assert.Equal(t, tt.want, got.Status)
		})
	}
}
