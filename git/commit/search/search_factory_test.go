package search

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SearchFactoryTestSuite struct {
	suite.Suite
}

func TestSearchFactoryTestSuite(t *testing.T) {
	suite.Run(t, new(SearchFactoryTestSuite))
}

func (suite *SearchFactoryTestSuite) TestGetSearchAction_ShouldReturnCommitHashSearchAction() {
	action := GetSearchAction("hash", []git2go.Commit{})

	suite.IsType(commitHashSearch{}, action)
}

func (suite *SearchFactoryTestSuite) TestGetSearchAction_ShouldReturnCommitAuthorSearchAction() {
	action := GetSearchAction("author", []git2go.Commit{})

	suite.IsType(commitAuthorSearch{}, action)
}

func (suite *SearchFactoryTestSuite) TestGetSearchAction_ShouldReturnCommitMessageSearchAction() {
	action := GetSearchAction("message", []git2go.Commit{})

	suite.IsType(commitMessageSearch{}, action)
}

func (suite *SearchFactoryTestSuite) TestGetSearchAction_WhenTypeIsInvalid_ShouldReturnNil() {
	action := GetSearchAction("invalid", []git2go.Commit{})

	suite.Nil(action)
}
