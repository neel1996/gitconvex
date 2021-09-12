package branch

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/neel1996/gitconvex/validator"
	mocks2 "github.com/neel1996/gitconvex/validator/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type BranchDeleteTestSuite struct {
	suite.Suite
	mockController      *gomock.Controller
	branchDelete        Delete
	branchName          string
	repo                middleware.Repository
	mockBranch          *mocks.MockBranch
	branchValidator     validator.ValidatorWithStringFields
	mockRepo            *mocks.MockRepository
	mockBranchValidator *mocks2.MockValidatorWithStringFields
}

func TestBranchDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(BranchDeleteTestSuite))
}

func (suite *BranchDeleteTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.branchName = "delete_branch"

	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockBranchValidator = mocks2.NewMockValidatorWithStringFields(suite.mockController)
	suite.mockBranch = mocks.NewMockBranch(suite.mockController)

	suite.branchDelete = NewDeleteBranch(suite.mockRepo, suite.mockBranchValidator)
}

func (suite *BranchDeleteTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)
	suite.branchValidator = validator.NewBranchValidator()
	suite.branchName = "delete_branch"
	addErr := NewAddBranch(suite.repo, validator.NewBranchValidator()).AddBranch(suite.branchName, false, nil)
	if addErr != nil {
		fmt.Println(addErr)
	}
}

func (suite *BranchDeleteTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenBranchIsDeleted_ShouldReturnNil() {
	suite.branchDelete = NewDeleteBranch(suite.repo, suite.branchValidator)

	err := suite.branchDelete.DeleteBranch(suite.branchName)

	suite.Nil(err)
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenValidationFails_ShouldReturnError() {
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.branchName).Return(errors.New("VALIDATION_ERROR"))

	err := suite.branchDelete.DeleteBranch(suite.branchName)

	suite.NotNil(err)
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenBranchLookupFails_ShouldReturnError() {
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.branchName).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.branchName, git2go.BranchLocal).Return(nil, errors.New("LOOKUP_ERR"))

	err := suite.branchDelete.DeleteBranch(suite.branchName)

	suite.NotNil(err)
}

func (suite *BranchDeleteTestSuite) TestDeleteBranch_WhenBranchDeleteFails_ShouldReturnError() {
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.branchName).Return(nil)
	suite.mockRepo.EXPECT().LookupBranch(suite.branchName, git2go.BranchLocal).Return(suite.mockBranch, nil)
	suite.mockBranch.EXPECT().Delete().Return(errors.New("DELETE_ERROR"))

	err := suite.branchDelete.DeleteBranch(suite.branchName)

	suite.NotNil(err)
}
