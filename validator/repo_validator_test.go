package validator

import (
	"github.com/golang/mock/gomock"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepoValidatorTestSuite struct {
	suite.Suite
	mockController *gomock.Controller
	mockRepo       *mocks.MockRepository
	repoValidator  Validator
}

func TestRepoValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(RepoValidatorTestSuite))
}

func (suite *RepoValidatorTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.repoValidator = NewRepoValidator()
}

func (suite *RepoValidatorTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *RepoValidatorTestSuite) TestValidate_WhenRepoIsNotNil_ShouldReturnNil() {
	err := suite.repoValidator.Validate(suite.mockRepo)

	suite.Nil(err)
}

func (suite *RepoValidatorTestSuite) TestValidate_WhenRepoIsNil_ShouldReturnRepoNilError() {
	suite.repoValidator = NewRepoValidator()

	err := suite.repoValidator.Validate(nil)

	suite.NotNil(err)
	suite.Equal(NilRepoError, err)
}
