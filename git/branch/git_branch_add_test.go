package branch

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/neel1996/gitconvex/validator"
	validatorMock "github.com/neel1996/gitconvex/validator/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type BranchAddTestSuite struct {
	suite.Suite
	repo                 middleware.Repository
	mockController       *gomock.Controller
	mockRepo             *mocks.MockRepository
	mockReference        *mocks.MockReference
	mockBranchValidation *validatorMock.MockValidatorWithStringFields
	mockCommit           *mocks.MockCommit
	branchName           string
	branchValidation     validator.ValidatorWithStringFields
	branchAdd            Add
}

func TestBranchAddTestSuite(t *testing.T) {
	suite.Run(t, new(BranchAddTestSuite))
}

func (suite *BranchAddTestSuite) SetupSuite() {
	suite.branchName = uuid.New().String()
}

func (suite *BranchAddTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)
	suite.branchValidation = validator.NewBranchValidator()

	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockReference = mocks.NewMockReference(suite.mockController)
	suite.mockCommit = mocks.NewMockCommit(suite.mockController)
	suite.mockBranchValidation = validatorMock.NewMockValidatorWithStringFields(suite.mockController)

	suite.branchAdd = NewAddBranch(suite.mockRepo, suite.mockBranchValidation)
}

func (suite *BranchAddTestSuite) TearDownSuite() {
	_ = NewDeleteBranch(suite.repo, validator.NewBranchValidator()).DeleteBranch(suite.branchName)
	suite.mockController.Finish()
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenBranchAdditionSucceeds_ShouldReturnNil() {
	suite.branchAdd = NewAddBranch(suite.repo, suite.branchValidation)

	branchAddError := suite.branchAdd.AddBranch(suite.branchName, false, nil)

	suite.Nil(branchAddError)
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenBranchValidationFails_ShouldReturnError() {
	suite.mockBranchValidation.EXPECT().ValidateWithFields("").Return(EmptyBranchNameError)

	branchAddError := suite.branchAdd.AddBranch("", false, nil)

	suite.NotNil(branchAddError)
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenHeadIsNil_ShouldReturnError() {
	suite.mockBranchValidation.EXPECT().ValidateWithFields(suite.branchName).Return(nil)
	suite.mockRepo.EXPECT().Head().Return(nil, errors.New("HEAD_ERROR"))

	branchAddError := suite.branchAdd.AddBranch(suite.branchName, false, nil)

	suite.NotNil(branchAddError)
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenHeadCommitLookupFails_ShouldReturnError() {
	oid, _ := git2go.NewOid("0608f9c6c97f386d1fa3948f3b8a61ae1cdb5621")

	suite.mockBranchValidation.EXPECT().ValidateWithFields(suite.branchName).Return(nil)
	suite.mockRepo.EXPECT().Head().Return(suite.mockReference, nil)
	suite.mockReference.EXPECT().Target().Return(oid)
	suite.mockRepo.EXPECT().LookupCommit(oid).Return(nil, errors.New("LOOKUP_ERROR"))

	branchAddError := suite.branchAdd.AddBranch(suite.branchName, false, nil)

	suite.NotNil(branchAddError)
}

func (suite *BranchAddTestSuite) TestAddBranch_WhenCreateBranchFails_ShouldReturnError() {
	oid, _ := git2go.NewOid("0608f9c6c97f386d1fa3948f3b8a61ae1cdb5621")
	commit := &git2go.Commit{}

	suite.mockBranchValidation.EXPECT().ValidateWithFields(suite.branchName).Return(nil)
	suite.mockRepo.EXPECT().Head().Return(suite.mockReference, nil)
	suite.mockReference.EXPECT().Target().Return(oid)
	suite.mockRepo.EXPECT().LookupCommit(oid).Return(commit, nil)
	suite.mockRepo.EXPECT().CreateBranch(
		suite.branchName,
		commit,
		false,
	).Return(nil, errors.New("BRANCH_CREATE_ERROR"))

	branchAddError := suite.branchAdd.AddBranch(suite.branchName, false, nil)

	suite.NotNil(branchAddError)
}
