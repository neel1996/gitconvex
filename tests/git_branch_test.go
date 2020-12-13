package tests

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	git2 "github.com/neel1996/gitconvex-server/git"
	"github.com/neel1996/gitconvex-server/global"
	"os"
	"path"
	"testing"
)

func TestGetBranchList(t *testing.T) {
	b := make(chan git2.Branch)
	cwd, _ := os.Getwd()
	r, _ := git.PlainOpen(path.Join(cwd, ".."))

	type args struct {
		repo       *git.Repository
		branchChan chan git2.Branch
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "Git branch list test case", args: struct {
			repo       *git.Repository
			branchChan chan git2.Branch
		}{repo: r, branchChan: b}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go git2.GetBranchList(tt.args.repo, tt.args.branchChan)
			branchList := <-tt.args.branchChan

			cBranch := branchList.CurrentBranch
			aBranch := branchList.AllBranchList
			bList := branchList.BranchList

			logger := global.Logger{}
			logger.Log(fmt.Sprintf("%s - %+v - %+v", cBranch, aBranch, bList), global.StatusInfo)

			if cBranch == "" || len(aBranch) == 0 || len(bList) == 0 {
				t.Error("Required results are not available")
			}
		})
	}
}
