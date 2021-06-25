package remote

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type RemoteHostTestSuite struct {
	suite.Suite
	host Host
}

func TestRemoteHostTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteHostTestSuite))
}

func (suite *RemoteHostTestSuite) SetupTest() {
	suite.host = NewRemoteHost("https://github.com/test/test.git")
}

func (suite *RemoteHostTestSuite) TestGetRemoteHost_WhenRemoteUrlIsPassed_ShouldReturnHostName() {
	expectedHost := "github"
	remoteHost := suite.host.GetRemoteHostForUrl()

	suite.Equal(expectedHost, remoteHost)
}

func (suite *RemoteHostTestSuite) TestGetRemoteHost_WhenRemoteUrlContainsUnknownHost_ShouldReturnEmptyString() {
	suite.host = NewRemoteHost("https://orca.com/test/test.git")
	remoteHost := suite.host.GetRemoteHostForUrl()

	suite.Equal("", remoteHost)
}
