package checkout

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/branch"
	branchMocks "github.com/neel1996/gitconvex/git/branch/mocks"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/neel1996/gitconvex/validator"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type CheckOutRemoteBranchTestSuite struct {
	suite.Suite
	repo                     middleware.Repository
	mockController           *gomock.Controller
	mockRepo                 *mocks.MockRepository
	mockBranch               *mocks.MockBranch
	mockCommit               *mocks.MockCommit
	branchName               string
	remoteBranchName         string
	localBranchName          string
	localBranchReferenceName string
	branchValidation         validator.ValidatorWithStringFields
	addBranch                branch.Add
	mockValidation           *branchMocks.MockValidation
	mockAddBranch            *branchMocks.MockAdd
	checkOutRemoteBranch     Checkout
}

func TestCheckOutRemoteBranchTestSuite(t *testing.T) {
	suite.Run(t, new(CheckOutRemoteBranchTestSuite))
}

func (suite *CheckOutRemoteBranchTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)
	suite.branchName = "remotes/origin/test_remote"
	suite.remoteBranchName = "origin/test_remote"
	suite.localBranchName = "test_remote"
	suite.localBranchReferenceName = "refs/heads/test_remote"
}

func (suite *CheckOutRemoteBranchTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockBranch = mocks.NewMockBranch(suite.mockController)
	suite.mockCommit = mocks.NewMockCommit(suite.mockController)
	suite.mockValidation = branchMocks.NewMockValidation(suite.mockController)
	suite.mockAddBranch = branchMocks.NewMockAdd(suite.mockController)
	suite.branchValidation = validator.NewBranchValidator()
	suite.addBranch = branch.NewAddBranch(suite.repo, suite.branchValidation)

	suite.checkOutRemoteBranch = NewCheckoutRemoteBranch(suite.mockRepo, suite.branchName, suite.mockAddBranch)
}

func (suite *CheckOutRemoteBranchTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *CheckOutRemoteBranchTestSuite) TestCheckoutBranch_WhenBranchIsValid_ShouldCheckoutLocalBranch() {
	suite.checkOutRemoteBranch = NewCheckoutRemoteBranch(suite.repo, suite.branchName, branch.NewAddBranch(suite.repo, validator.NewBranchValidator()))

	err := suite.checkOutRemoteBranch.CheckoutBranch()

	suite.Nil(err)
}

func (suite *CheckOutRemoteBranchTestSuite) TestCheckoutBranch_WhenBranchLookupFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().LookupBranch(suite.remoteBranchName, git2go.BranchRemote).Return(nil, errors.New("LOOKUP_ERROR"))

	err := suite.checkOutRemoteBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutRemoteBranchTestSuite) TestCheckoutBranch_WhenCommitLookupFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().LookupBranch(suite.remoteBranchName, git2go.BranchRemote).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(nil, errors.New("LOOKUP_ERROR"))

	err := suite.checkOutRemoteBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutRemoteBranchTestSuite) TestCheckoutBranch_WhenFetchCommitTreeFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().LookupBranch(suite.remoteBranchName, git2go.BranchRemote).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().Tree().Return(nil, errors.New("TREE_ERROR"))

	err := suite.checkOutRemoteBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutRemoteBranchTestSuite) TestCheckoutBranch_WhenCheckoutTreeFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().LookupBranch(suite.remoteBranchName, git2go.BranchRemote).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().Tree().Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CheckoutTree(gomock.Any(), gomock.Any()).Return(errors.New("TREE_CHECKOUT_ERROR"))

	err := suite.checkOutRemoteBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutRemoteBranchTestSuite) TestCheckoutBranch_WhenLookupLocalBranchFails_ShouldAddNewLocalBranch() {
	gitCommit := &git2go.Commit{}

	suite.mockRepo.EXPECT().LookupBranch(suite.remoteBranchName, git2go.BranchRemote).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().Tree().Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CheckoutTree(gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.localBranchName, git2go.BranchLocal).Return(nil, errors.New("LOOKUP_ERROR"))
	suite.mockCommit.EXPECT().GetGitCommit().Return(gitCommit)
	suite.mockAddBranch.EXPECT().AddBranch(suite.branchName, false, gitCommit).Return(nil)
	suite.mockRepo.EXPECT().SetHead(suite.localBranchReferenceName).Return(nil)

	err := suite.checkOutRemoteBranch.CheckoutBranch()

	suite.Nil(err)
}

func (suite *CheckOutRemoteBranchTestSuite) TestCheckoutBranch_WhenLocalAddBranchFails_ShouldReturnError() {
	gitCommit := &git2go.Commit{}

	suite.mockRepo.EXPECT().LookupBranch(suite.remoteBranchName, git2go.BranchRemote).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().Tree().Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CheckoutTree(gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.localBranchName, git2go.BranchLocal).Return(nil, errors.New("LOOKUP_ERROR"))
	suite.mockCommit.EXPECT().GetGitCommit().Return(gitCommit)
	suite.mockAddBranch.EXPECT().AddBranch(suite.branchName, false, gitCommit).Return(errors.New("ADD_ERROR"))

	err := suite.checkOutRemoteBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutRemoteBranchTestSuite) TestCheckoutBranch_WhenSetHeadFailsAfterAddingNewBranch_ShouldReturnError() {
	gitCommit := &git2go.Commit{}

	suite.mockRepo.EXPECT().LookupBranch(suite.remoteBranchName, git2go.BranchRemote).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().Tree().Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CheckoutTree(gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.localBranchName, git2go.BranchLocal).Return(nil, errors.New("LOOKUP_ERROR"))
	suite.mockCommit.EXPECT().GetGitCommit().Return(gitCommit)
	suite.mockAddBranch.EXPECT().AddBranch(suite.branchName, false, gitCommit).Return(nil)
	suite.mockRepo.EXPECT().SetHead(suite.localBranchReferenceName).Return(errors.New("SET_HEAD_ERROR"))

	err := suite.checkOutRemoteBranch.CheckoutBranch()

	suite.NotNil(err)
}

func (suite *CheckOutRemoteBranchTestSuite) TestCheckoutBranch_WhenSetHeadFails_ShouldReturnError() {
	suite.mockRepo.EXPECT().LookupBranch(suite.remoteBranchName, git2go.BranchRemote).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Target().Return(&git2go.Oid{})
	suite.mockRepo.EXPECT().LookupCommitV2(gomock.Any()).Return(suite.mockCommit, nil)
	suite.mockCommit.EXPECT().Tree().Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CheckoutTree(gomock.Any(), gomock.Any()).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.localBranchName, git2go.BranchLocal).Return(nil, nil)
	suite.mockRepo.EXPECT().SetHead(suite.localBranchReferenceName).Return(errors.New("SET_HEAD_ERROR"))

	err := suite.checkOutRemoteBranch.CheckoutBranch()

	suite.NotNil(err)
}
