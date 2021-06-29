package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type BranchAddTestSuite struct {
	suite.Suite
	repo      *git2go.Repository
	branchAdd Add
}

func TestBranchAddTestSuite(t *testing.T) {
	suite.Run(t, new(BranchAddTestSuite))
}

func (suite *BranchAddTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = r
	suite.branchAdd = NewAddBranch(suite.repo, "test_1", false, nil)
}

func (suite *BranchAddTestSuite) TearDownTest() {
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	NewDeleteBranch(r, "test_1").DeleteBranch()
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
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	suite.branchAdd = NewAddBranch(r, "", false, nil)
	branchAddError := suite.branchAdd.AddBranch()

	suite.NotNil(branchAddError)
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenBranchAdditionFails_ShouldReturnError() {
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	suite.branchAdd = NewAddBranch(r, "test_1", false, nil)

	//Adding duplicate branches
	_ = suite.branchAdd.AddBranch()
	branchAddError := suite.branchAdd.AddBranch()

	suite.NotNil(branchAddError)
}
