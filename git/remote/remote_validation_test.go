package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RemoteValidationTestSuite struct {
	suite.Suite
	repo                 *git2go.Repository
	validateRemoteFields Validation
}

func TestRemoteValidationTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteValidationTestSuite))
}

func (suite *RemoteValidationTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = r
	suite.validateRemoteFields = NewRemoteValidation()
}

func (suite *RemoteValidationTestSuite) TestValidateRemoteFields_WhenAllFieldsAreValid_ShouldReturnNil() {
	wantErr := suite.validateRemoteFields.ValidateRemoteFields(suite.repo)

	fmt.Println(wantErr)

	suite.Nil(wantErr)
}

func (suite *RemoteValidationTestSuite) TestValidateRemoteFields_WhenRepoIsNil_ShouldReturnError() {
	wantErr := suite.validateRemoteFields.ValidateRemoteFields(nil)
	wantErrorText := "repo is nil"

	suite.NotNil(wantErr)
	suite.Equal(wantErrorText, wantErr.Error())
}

func (suite *RemoteValidationTestSuite) TestValidateRemoteFields_WhenRemoteCollectionIsNil_ShouldReturnError() {
	wantErr := suite.validateRemoteFields.ValidateRemoteFields(&git2go.Repository{
		Remotes: git2go.RemoteCollection{},
	})
	wantErrorText := "remote collection is nil"

	suite.NotNil(wantErr)
	suite.Equal(wantErrorText, wantErr.Error())
}
