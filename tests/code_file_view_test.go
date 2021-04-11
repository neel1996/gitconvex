package tests

import (
	"fmt"
	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/api"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestCodeFileView(t *testing.T) {
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

	type args struct {
		repo     *git.Repository
		repoPath string
		fileName string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Code view API test case", args: struct {
			repo     *git.Repository
			repoPath string
			fileName string
		}{repo: r, repoPath: repoPath, fileName: "README.md"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj api.CodeViewInterface
			obj = api.CodeViewInputs{
				RepoPath: tt.args.repoPath,
				FileName: tt.args.fileName,
			}
			assert.NotNil(t, obj.CodeFileView(), "")
		})
	}
}
