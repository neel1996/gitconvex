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
	remoteFields         []string
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

	suite.remoteFields = []string{"origin", "origin_1"}
	suite.repo = r
	suite.validateRemoteFields = NewRemoteValidation(suite.repo, suite.remoteFields[0], suite.remoteFields[1])
}

func (suite *RemoteValidationTestSuite) TestValidateRemoteFields_WhenAllFieldsAreValid_ShouldReturnNil() {
	wantErr := suite.validateRemoteFields.ValidateRemoteFields()

	fmt.Println(wantErr)

	suite.Nil(wantErr)
}

func (suite *RemoteValidationTestSuite) TestValidateRemoteFields_WhenRepoIsNil_ShouldReturnError() {
	suite.validateRemoteFields = NewRemoteValidation(nil)
	wantErr := suite.validateRemoteFields.ValidateRemoteFields()
	wantErrorText := "repo is nil"

	suite.NotNil(wantErr)
	suite.Equal(wantErrorText, wantErr.Error())
}

func (suite *RemoteValidationTestSuite) TestValidateRemoteFields_WhenRemoteCollectionIsNil_ShouldReturnError() {
	suite.validateRemoteFields = NewRemoteValidation(&git2go.Repository{
		Remotes: git2go.RemoteCollection{},
	})
	wantErr := suite.validateRemoteFields.ValidateRemoteFields()
	wantErrorText := "remote collection is nil"

	suite.NotNil(wantErr)
	suite.Equal(wantErrorText, wantErr.Error())
}

func (suite *RemoteValidationTestSuite) TestValidateRemoteFields_WhenRemoteFieldsAreEmpty_ShouldReturnError() {
	suite.validateRemoteFields = NewRemoteValidation(suite.repo, "", "")
	wantErr := suite.validateRemoteFields.ValidateRemoteFields()
	wantErrorText := "one or more remote fields are empty"

	suite.NotNil(wantErr)
	suite.Equal(wantErrorText, wantErr.Error())
}
