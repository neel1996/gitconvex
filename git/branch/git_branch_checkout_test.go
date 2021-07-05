package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type BranchCheckoutTestSuite struct {
	suite.Suite
	branchName     string
	checkoutBranch Checkout
	repo           *git2go.Repository
	noHeadRepo     *git2go.Repository
}

func TestBranchCheckoutTestSuite(t *testing.T) {
	suite.Run(t, new(BranchCheckoutTestSuite))
}

func (suite *BranchCheckoutTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.repo = r
	suite.noHeadRepo = noHeadRepo
	suite.branchName = "test_checkout"
	addErr := NewAddBranch(suite.repo, suite.branchName, false, nil).AddBranch()
	if addErr != nil {
		fmt.Println(addErr)
	}
	suite.checkoutBranch = NewBranchCheckout(suite.repo, suite.branchName)
}

func (suite *BranchCheckoutTestSuite) TearDownSuite() {
	suite.branchName = "master"
	err := suite.checkoutBranch.CheckoutBranch()
	if err != nil {
		return
	}
	_ = NewDeleteBranch(suite.repo, suite.branchName).DeleteBranch()
}

func (suite *BranchCheckoutTestSuite) TestCheckoutBranch_WhenBranchIsCheckedOut_ShouldReturnNil() {
	err := suite.checkoutBranch.CheckoutBranch()

	suite.Nil(err)
}

func (suite *BranchCheckoutTestSuite) TestCheckoutBranch_WhenRemoteBranchIsCheckedOut_ShouldReturnNil() {
	suite.checkoutBranch = NewBranchCheckout(suite.repo, "remotes/origin/master")
	err := suite.checkoutBranch.CheckoutBranch()

	suite.Nil(err)
}

func (suite *BranchCheckoutTestSuite) TestCheckoutBranch_WhenRepoIsNil_ShouldReturnError() {
	suite.checkoutBranch = NewBranchCheckout(nil, "test_branch")
	err := suite.checkoutBranch.CheckoutBranch()

	suite.NotNil(err)
	suite.Equal("repo is nil", err.Error())
}

func (suite *BranchCheckoutTestSuite) TestCheckoutBranch_WhenBranchNameIsEmpty_ShouldReturnError() {
	suite.checkoutBranch = NewBranchCheckout(suite.repo, "")
	err := suite.checkoutBranch.CheckoutBranch()

	suite.NotNil(err)
	suite.Equal("branch name is empty", err.Error())
}

func (suite *BranchCheckoutTestSuite) TestCheckoutBranch_WhenRepoHasNoHead_ShouldReturnError() {
	suite.checkoutBranch = NewBranchCheckout(suite.noHeadRepo, suite.branchName)
	err := suite.checkoutBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *BranchCheckoutTestSuite) TestCheckoutBranch_WhenNonExistingBranchIsSelected_ShouldReturnError() {
	suite.checkoutBranch = NewBranchCheckout(suite.repo, "no_exists")
	err := suite.checkoutBranch.CheckoutBranch()

	suite.NotNil(err)
}
