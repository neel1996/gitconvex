package commit

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/mocks"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type CommitLookupTestSuite struct {
	suite.Suite
	repo           middleware.Repository
	mockController *gomock.Controller
	mockRepo       *mocks.MockRepository
	commitHash     string
	lookup         Lookup
}

func TestCommitLookupTestSuite(t *testing.T) {
	suite.Run(t, new(CommitLookupTestSuite))
}

func (suite *CommitLookupTestSuite) SetupSuite() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}
	commits, _ := NewListAllLogs(middleware.NewRepository(r), nil, nil).Get()

	suite.commitHash = commits[0].Id().String()
}

func (suite *CommitLookupTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.repo = middleware.NewRepository(r)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.lookup = NewLookup(suite.repo)
}

func (suite *CommitLookupTestSuite) TestWithReference_WhenCommitHashIsValid_ShouldReturnCommit() {
	wantCommit, err := suite.lookup.WithReferenceId(suite.commitHash)

	suite.Nil(err)
	suite.Equal(suite.commitHash, wantCommit.Id().String())
}

func (suite *CommitLookupTestSuite) TestWithReference_WhenCommitHashIsInValid_ShouldReturnError() {
	_, err := suite.lookup.WithReferenceId("INVALID_HASH")

	suite.NotNil(err)
	suite.Equal(OidConversionError, err)
}

func (suite *CommitLookupTestSuite) TestWithReference_WhenCommitIsNotPresentInTheRepo_ShouldReturnError() {
	suite.lookup = NewLookup(suite.mockRepo)

	oid, _ := git2go.NewOid("a94a8fe5ccb19ba61c4c0873d391e987982fbbd3")
	suite.mockRepo.EXPECT().LookupCommit(oid).Return(nil, errors.New("LOOKUP_ERR"))

	_, err := suite.lookup.WithReferenceId("a94a8fe5ccb19ba61c4c0873d391e987982fbbd3")

	suite.NotNil(err)
	suite.Equal(CommitLookupError, err)
}
