package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type BranchAddTestSuite struct {
	suite.Suite
	repo       *git2go.Repository
	noHeadRepo *git2go.Repository
	branchName string
	branchAdd  Add
}

func TestBranchAddTestSuite(t *testing.T) {
	suite.Run(t, new(BranchAddTestSuite))
}

func (suite *BranchAddTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.repo = r
	suite.noHeadRepo = noHeadRepo
	suite.branchName = "test_1"
	suite.branchAdd = NewAddBranch(suite.repo, suite.branchName, false, nil)
}

func (suite *BranchAddTestSuite) TearDownSuite() {
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	_ = NewDeleteBranch(r, suite.branchName).DeleteBranch()
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenBranchAdditionSucceeds_ShouldReturnNil() {
	branchAddError := suite.branchAdd.AddBranch()

	suite.Nil(branchAddError)
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenRepoIsNil_ShouldReturnError() {
	suite.branchAdd = NewAddBranch(nil, "", false, nil)
	branchAddError := suite.branchAdd.AddBranch()

	suite.NotNil(branchAddError)
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenBranchNameIsEmpty_ShouldReturnError() {
	suite.branchAdd = NewAddBranch(suite.repo, "", false, nil)
	branchAddError := suite.branchAdd.AddBranch()

	suite.NotNil(branchAddError)
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenHeadIsNil_ShouldReturnError() {
	suite.branchAdd = NewAddBranch(suite.noHeadRepo, suite.branchName, false, nil)
	branchAddError := suite.branchAdd.AddBranch()

	suite.NotNil(branchAddError)
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenBranchAdditionFails_ShouldReturnError() {
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	suite.branchAdd = NewAddBranch(r, "master", false, nil)

	//Adding duplicate branches
	_ = suite.branchAdd.AddBranch()
	branchAddError := suite.branchAdd.AddBranch()

	suite.NotNil(branchAddError)
}
