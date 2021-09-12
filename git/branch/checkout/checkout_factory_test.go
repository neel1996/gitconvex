package checkout

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/neel1996/gitconvex/validator"
	validatorMock "github.com/neel1996/gitconvex/validator/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CheckoutFactoryTestSuite struct {
	suite.Suite
	repo                middleware.Repository
	branchName          string
	remoteBranchName    string
	mockController      *gomock.Controller
	mockRepo            *mocks.MockRepository
	mockBranchValidator *validatorMock.MockValidatorWithStringFields
	mockRepoValidator   *validatorMock.MockValidator
	checkoutFactory     Factory
}

func TestCheckoutFactoryTestSuite(t *testing.T) {
	suite.Run(t, new(CheckoutFactoryTestSuite))
}

func (suite *CheckoutFactoryTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockBranchValidator = validatorMock.NewMockValidatorWithStringFields(suite.mockController)
	suite.mockRepoValidator = validatorMock.NewMockValidator(suite.mockController)
	suite.branchName = "test_branch"
	suite.remoteBranchName = "remotes/origin/test_branch"
}

func (suite *CheckoutFactoryTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *CheckoutFactoryTestSuite) TestGetCheckoutAction_WhenBranchIsLocal_ShouldReturnLocalCheckoutAction() {
	suite.checkoutFactory = NewCheckoutFactory(suite.mockRepo, suite.mockRepoValidator, suite.mockBranchValidator)

	suite.mockRepoValidator.EXPECT().Validate(suite.mockRepo).Return(nil)
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.branchName).Return(nil)

	wantAction := NewCheckOutLocalBranch(suite.mockRepo, suite.branchName)
	gotAction := suite.checkoutFactory.GetCheckoutAction(suite.branchName)

	suite.Equal(wantAction, gotAction)
}

func (suite *CheckoutFactoryTestSuite) TestGetCheckoutAction_WhenBranchIsRemote_ShouldReturnRemoteCheckoutAction() {
	suite.checkoutFactory = NewCheckoutFactory(suite.mockRepo, suite.mockRepoValidator, suite.mockBranchValidator)

	suite.mockRepoValidator.EXPECT().Validate(suite.mockRepo).Return(nil)
	suite.mockBranchValidator.EXPECT().ValidateWithFields(suite.remoteBranchName).Return(nil)

	wantAction := NewCheckoutRemoteBranch(suite.mockRepo, suite.remoteBranchName, nil)
	gotAction := suite.checkoutFactory.GetCheckoutAction(suite.remoteBranchName)

	suite.Equal(wantAction, gotAction)
}

func (suite *CheckoutFactoryTestSuite) TestGetCheckoutAction_WhenRepoValidationFails_ShouldReturnNil() {
	suite.checkoutFactory = NewCheckoutFactory(nil, suite.mockRepoValidator, suite.mockBranchValidator)

	suite.mockRepoValidator.EXPECT().Validate(nil).Return(validator.NilRepoError)
	gotAction := suite.checkoutFactory.GetCheckoutAction(suite.branchName)

	suite.Nil(gotAction)
}

func (suite *CheckoutFactoryTestSuite) TestGetCheckoutAction_WhenBranchValidationFails_ShouldReturnNil() {
	suite.checkoutFactory = NewCheckoutFactory(suite.mockRepo, suite.mockRepoValidator, suite.mockBranchValidator)

	suite.mockRepoValidator.EXPECT().Validate(suite.mockRepo).Return(nil)
	suite.mockBranchValidator.EXPECT().ValidateWithFields("").Return(errors.New("VALIDATE_ERROR"))

	gotAction := suite.checkoutFactory.GetCheckoutAction("")

	suite.Nil(gotAction)
}
