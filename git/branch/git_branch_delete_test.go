package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type BranchDeleteTestSuite struct {
	suite.Suite
	branchDelete Delete
	branchName   string
	repo         *git2go.Repository
	noHeadRepo   *git2go.Repository
}

func TestBranchDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(BranchDeleteTestSuite))
}

func (suite *BranchDeleteTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.repo = r
	suite.noHeadRepo = noHeadRepo
	suite.branchName = "delete_branch"
	suite.branchDelete = NewDeleteBranch(suite.repo, suite.branchName)
}

func (suite *BranchDeleteTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = r
	suite.branchName = "delete_branch"
	addErr := NewAddBranch(suite.repo, suite.branchName, false, nil).AddBranch()
	if addErr != nil {
		fmt.Println(addErr)
	}
}

func (suite *BranchDeleteTestSuite) TearDownSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	_ = NewBranchCheckout(r, "master").CheckoutBranch()
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenBranchIsDeleted_ShouldReturnNil() {
	err := suite.branchDelete.DeleteBranch()

	suite.Nil(err)
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenRepoIsNil_ShouldReturnError() {
	suite.branchDelete = NewDeleteBranch(nil, suite.branchName)
	err := suite.branchDelete.DeleteBranch()

	suite.NotNil(err)
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenBranchNameIsEmpty_ShouldReturnError() {
	suite.branchDelete = NewDeleteBranch(suite.repo, "")
	err := suite.branchDelete.DeleteBranch()

	suite.NotNil(err)
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenRepoHasNoHead_ShouldReturnError() {
	suite.branchDelete = NewDeleteBranch(suite.noHeadRepo, suite.branchName)
	err := suite.branchDelete.DeleteBranch()

	suite.NotNil(err)
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenBranchDoesNotExists_ShouldReturnError() {
	suite.branchDelete = NewDeleteBranch(suite.repo, "no_exists")
	err := suite.branchDelete.DeleteBranch()

	suite.NotNil(err)
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenBranchIsCurrentBranch_ShouldReturnError() {
	head, _ := suite.repo.Head()
	name, _ := head.Branch().Name()

	suite.branchDelete = NewDeleteBranch(suite.repo, name)
	err := suite.branchDelete.DeleteBranch()

	suite.NotNil(err)
}
