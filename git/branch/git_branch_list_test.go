package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type BranchListTestSuite struct {
	suite.Suite
	repo       *git2go.Repository
	noHeadRepo *git2go.Repository
	branchList List
}

func TestBranchListTestSuite(t *testing.T) {
	suite.Run(t, new(BranchListTestSuite))
}

func (suite *BranchListTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.repo = r
	suite.noHeadRepo = noHeadRepo
	suite.branchList = NewBranchList(suite.repo)
}

func (suite *BranchListTestSuite) TestListBranches_WhenRepoHasBranches_ShouldReturnBranchList() {
	branchList, err := suite.branchList.ListBranches()

	suite.Nil(err)
	suite.Greater(len(branchList.BranchList), 2)
	suite.Greater(len(branchList.AllBranchList), 2)
	suite.Equal("master", branchList.CurrentBranch)
}

func (suite *BranchListTestSuite) TestListBranches_WhenRepoIsNil_ShouldReturnError() {
	suite.branchList = NewBranchList(nil)
	branchList, err := suite.branchList.ListBranches()

	suite.NotNil(err)
	suite.Empty(branchList)
}

func (suite *BranchListTestSuite) TestListBranches_WhenRepoHasNoHead_ShouldReturnError() {
	suite.branchList = NewBranchList(suite.noHeadRepo)
	branchList, err := suite.branchList.ListBranches()

	suite.NotNil(err)
	suite.Empty(branchList)
}
