package test_utils

import (
	"fmt"
	"github.com/neel1996/gitconvex/git"
	"github.com/neel1996/gitconvex/git/commit"
	"github.com/neel1996/gitconvex/git/middleware"
	"io/ioutil"
)

func StageAndCommitTestFile(repo middleware.Repository, branchName string, testFile string) {
	AddTestBranch(repo, branchName)
	CheckoutTestBranch(repo, branchName)

	for i := 1; i <= 2; i++ {
		err := ioutil.WriteFile(testFile, []byte{byte(i)}, 0644)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		git.StageItemStruct{
			Repo:     repo.GetGitRepository(),
			FileItem: testFile,
		}.StageItem()
		fmt.Println(commit.NewCommitChanges(repo, []string{fmt.Sprintf("Branch compare test commit - %v", i)}).Add())
	}
}
