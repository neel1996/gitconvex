package tests

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestCommitChanges(t *testing.T) {
	var repoPath string
	var r *git2go.Repository
	cwd, _ := os.Getwd()
	mockRepoPath := path.Join(cwd, "../..") + "/starfleet"
	currentEnv := os.Getenv("GOTESTENV")
	fmt.Println("Environment : " + currentEnv)

	if currentEnv == "ci" {
		repoPath = mockRepoPath
		r, _ = git2go.OpenRepository(repoPath)
	}

	sampleFile := "untracked.txt"
	err := ioutil.WriteFile(repoPath+"/"+sampleFile, []byte{byte(63)}, 0755)

	var stageAllObjects git2.StageAllInterface
	stageAllObjects = git2.StageAllStruct{Repo: r}

	fmt.Println(err)
	fmt.Println(stageAllObjects.StageAllItems())

	type args struct {
		repo          *git2go.Repository
		commitMessage string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Git commit changes test case", args: struct {
			repo          *git2go.Repository
			commitMessage string
		}{repo: r, commitMessage: "Test commit"}, want: "COMMIT_DONE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testObj git2.CommitInterface
			testObj = git2.CommitStruct{
				Repo:          tt.args.repo,
				CommitMessage: tt.args.commitMessage,
				RepoPath:      repoPath,
			}

			got := testObj.CommitChanges()
			assert.Equal(t, tt.want, got)
		})
	}
}
