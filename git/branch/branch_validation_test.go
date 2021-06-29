package branch

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type BranchFieldValidationTestSuite struct {
	suite.Suite
	repo       *git2go.Repository
	branchName string
	validation Validation
}

func TestBranchFieldValidationTestSuite(t *testing.T) {
	suite.Run(t, new(BranchFieldValidationTestSuite))
}

func (suite *BranchFieldValidationTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = r
	suite.branchName = "test_branch"
	suite.validation = NewBranchFieldsValidation(suite.repo, suite.branchName)
}

func (suite *BranchFieldValidationTestSuite) TestValidateBranchFields_WhenAllFieldsAreValid_ShouldReturnNil() {
	err := suite.validation.ValidateBranchFields()

	suite.Nil(err)
}

func (suite *BranchFieldValidationTestSuite) TestValidateBranchFields_WhenRepoIsNil_ShouldReturnError() {
	suite.validation = NewBranchFieldsValidation(nil, "test_branch")
	err := suite.validation.ValidateBranchFields()

	suite.NotNil(err)
	suite.Equal("repo is nil", err.Error())
}

func (suite *BranchFieldValidationTestSuite) TestValidateBranchFields_WhenBranchNameIsEmpty_ShouldReturnError() {
	suite.validation = NewBranchFieldsValidation(suite.repo, "")
	err := suite.validation.ValidateBranchFields()

	suite.NotNil(err)
	suite.Equal("branch name is empty", err.Error())
}
