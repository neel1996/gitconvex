package validator

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ValidationErrorTestSuite struct {
	suite.Suite
	validationError Error
}

func TestValidationErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationErrorTestSuite))
}

func (suite *ValidationErrorTestSuite) TestError_ShouldReturnErrorString() {
	suite.validationError = NilRepoError

	suite.Equal("Repo is nil", suite.validationError.Error())
}
