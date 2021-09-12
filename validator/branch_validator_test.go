package validator

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type BranchValidatorTestSuite struct {
	suite.Suite
	branchValidator ValidatorWithStringFields
}

func TestBranchValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(BranchValidatorTestSuite))
}

func (suite *BranchValidatorTestSuite) SetupTest() {
	suite.branchValidator = NewBranchValidator()
}

func (suite *BranchValidatorTestSuite) TestValidate_WhenPassedBranchIsValid_ShouldReturnNil() {
	err := suite.branchValidator.ValidateWithFields("test_1")

	suite.Nil(err)
}

func (suite *BranchValidatorTestSuite) TestValidate_WhenPassedBranchesAreValid_ShouldReturnNil() {
	err := suite.branchValidator.ValidateWithFields("test_1", "test_2")

	suite.Nil(err)
}

func (suite *BranchValidatorTestSuite) TestValidate_WhenPassedEmptyStringIsPassed_ShouldReturnEmptybranchError() {
	err := suite.branchValidator.ValidateWithFields()

	suite.NotNil(err)
	suite.Equal(EmptyBranchName, err)
}
