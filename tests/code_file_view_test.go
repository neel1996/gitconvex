package tests

import (
	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCodeFileView(t *testing.T) {
	r, _ := git.OpenRepository(TestRepo)

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
		}{repo: r, repoPath: r.Workdir(), fileName: "README.md"}},
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
