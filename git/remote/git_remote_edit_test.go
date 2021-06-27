package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RemoteEditTestSuite struct {
	suite.Suite
	repo       *git2go.Repository
	validation Validation
	editRemote Edit
}

func TestRemoteEditTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteEditTestSuite))
}

func (suite *RemoteEditTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.repo = r
	suite.validation = NewRemoteValidation()
	suite.editRemote = NewEditRemote(r, "origin", "https://github.com/neel1996/gitconvex-test.git", suite.validation)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteIsEdited_ShouldReturnNil() {
	wantErr := suite.editRemote.EditRemote()

	suite.Nil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRepoIsNil_ShouldReturnError() {
	suite.editRemote = NewEditRemote(nil, "origin", "https://github.com/neel1996/gitconvex-test.git", suite.validation)

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteCollectionIsNil_ShouldReturnError() {
	suite.editRemote = NewEditRemote(&git2go.Repository{
		Remotes: git2go.RemoteCollection{},
	}, "origin", "https://github.com/neel1996/gitconvex-test.git", suite.validation)

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}

func (suite *RemoteEditTestSuite) TestEditRemote_WhenRemoteEditFieldsAreEmpty_ShouldReturnError() {
	suite.editRemote = NewEditRemote(suite.repo, "", "", suite.validation)

	wantErr := suite.editRemote.EditRemote()

	suite.NotNil(wantErr)
}
