package tests

import (
	"os"
	"path"
	"testing"

	git "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/api"
	"github.com/neel1996/gitconvex-server/graph/model"
)

func TestCodeFileView(t *testing.T) {
	cwd, _ := os.Getwd()
	r, _ := git.OpenRepository(path.Join(cwd, ".."))

	repoPath := r.Path()
	expectedLine := "# gitconvex GoLang project"

	type args struct {
		repo     *git.Repository
		repoPath string
		fileName string
	}
	tests := []struct {
		name string
		args args
		want *model.CodeFileType
	}{
		{name: "Code view API test case", args: struct {
			repo     *git.Repository
			repoPath string
			fileName string
		}{repo: r, repoPath: path.Join(repoPath, ".."), fileName: "README.md"}, want: &model.CodeFileType{FileData: []*string{&expectedLine}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var obj api.CodeViewInterface
			obj = api.CodeViewInputs{
				RepoPath: tt.args.repoPath,
				FileName: tt.args.fileName,
			}
			if got := obj.CodeFileView(); *got.FileData[0] != *tt.want.FileData[0] {
				t.Errorf("CodeFileView() = %v, want %v", *got.FileData[0], *tt.want.FileData[0])
			}
		})
	}
}
