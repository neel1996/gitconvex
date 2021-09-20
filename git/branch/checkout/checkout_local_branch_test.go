package checkout

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/neel1996/gitconvex/test_utils"
	"github.com/neel1996/gitconvex/validator"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type CheckOutLocalBranchTestSuite struct {
	suite.Suite
	repo                middleware.Repository
	mockController      *gomock.Controller
	mockRepo            *mocks.MockRepository
	mockBranch          *mocks.MockBranch
	mockCommit          *mocks.MockCommit
	branchName          string
	checkOutLocalBranch Checkout
}

func TestCheckOutLocalBranchTestSuite(t *testing.T) {
	suite.Run(t, new(CheckOutLocalBranchTestSuite))
}

func (suite *CheckOutLocalBranchTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)
	suite.branchName = "test_checkout"

	test_utils.AddNewTestLocalBranch(suite.repo, suite.branchName)
}

func (suite *CheckOutLocalBranchTestSuite) TearDownSuite() {
	_ = branch.NewDeleteBranch(suite.repo, validator.NewBranchValidator()).DeleteBranch(suite.branchName)
}

func (suite *CheckOutLocalBranchTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockBranch = mocks.NewMockBranch(suite.mockController)
	suite.mockCommit = mocks.NewMockCommit(suite.mockController)

	suite.checkOutLocalBranch = NewCheckOutLocalBranch(suite.mockRepo, suite.branchName)
}

func (suite *CheckOutLocalBranchTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *CheckOutLocalBranchTestSuite) TestCheckoutBranch_WhenBranchIsValid_ShouldCheckoutLocalBranch() {
	suite.checkOutLocalBranch = NewCheckOutLocalBranch(suite.repo, suite.branchName)

	err := suite.checkOutLocalBranch.CheckoutBranch()

	suite.Nil(err)
}

func (suite *CheckOutLocalBranchTestSuite) TestCheckoutBranch_WhenBranchLookupFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().LookupBranch(suite.branchName, git2go.BranchLocal).Return(nil, errors.New("LOOKUP_ERROR"))

	err := suite.checkOutLocalBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutLocalBranchTestSuite) TestCheckoutBranch_WhenCommitLookupFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().LookupBranch(suite.branchName, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(nil, errors.New("LOOKUP_ERROR"))

	err := suite.checkOutLocalBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutLocalBranchTestSuite) TestCheckoutBranch_WhenFetchCommitTreeFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().LookupBranch(suite.branchName, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().Tree().Return(nil, errors.New("TREE_ERROR"))

	err := suite.checkOutLocalBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutLocalBranchTestSuite) TestCheckoutBranch_WhenCheckoutTreeFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().LookupBranch(suite.branchName, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().Tree().Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CheckoutTree(gomock.Any(), gomock.Any()).Return(errors.New("TREE_CHECKOUT_ERROR"))

	err := suite.checkOutLocalBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutLocalBranchTestSuite) TestCheckoutBranch_WhenSetHeadFails_ShouldReturnError() {
	referenceName := "refs/heads/" + suite.branchName

	suite.mockRepo.EXPECT().LookupBranch(suite.branchName, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().Tree().Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CheckoutTree(gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().SetHead(referenceName).Return(errors.New("SET_HEAD_ERROR"))

	err := suite.checkOutLocalBranch.CheckoutBranch()

	suite.NotNil(err)
}
