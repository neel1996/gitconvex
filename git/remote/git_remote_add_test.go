package remote

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type RemoteAddTestSuite struct {
	suite.Suite
	addRemote Add
}

func TestRemoteAddTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteAddTestSuite))
}

func (suite *RemoteAddTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	suite.addRemote = NewAddRemote(r, "new_origin", "https://github.com/neel1996/gitconvex-test.git")
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenNewRemoteIsAdded_ShouldReturnNoError() {
	err := suite.addRemote.NewRemote()

	suite.Nil(err)
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenRequiredFieldsAreEmpty_ShouldReturnError() {
	suite.addRemote = NewAddRemote(nil, "", "")

	err := suite.addRemote.NewRemote()

	suite.NotNil(err)
}

func (suite *RemoteAddTestSuite) TestAddNewRemote_WhenRemoteCreationFails_ShouldReturnError() {
	r, _ := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))

	suite.addRemote = NewAddRemote(r, "new_origin", "https://github.com/neel1996/gitconvex-test.git")

	err := suite.addRemote.NewRemote()

	suite.NotNil(err)
}
