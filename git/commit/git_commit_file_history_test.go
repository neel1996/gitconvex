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

type FileHistoryTestSuite struct {
	suite.Suite
	mockController *gomock.Controller
	repo           middleware.Repository
	mockRepo       *mocks.MockRepository
	mockCommit     *mocks.MockCommit
	fileHistory    FileHistory
}

func TestFileHistoryTestSuite(t *testing.T) {
	suite.Run(t, new(FileHistoryTestSuite))
}

func (suite *FileHistoryTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	suite.mockController = gomock.NewController(suite.T())
	suite.repo = middleware.NewRepository(r)
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockCommit = mocks.NewMockCommit(suite.mockController)
	suite.fileHistory = NewFileHistory(suite.mockRepo)
}

func (suite *FileHistoryTestSuite) TestGet_WhenRepoHasMoreThanOneCommit_ShouldReturnDiffForChosenCommit() {
	suite.fileHistory = NewFileHistory(suite.repo)
	commitLogs := suite.commitLogs()

	commit, _ := suite.repo.LookupCommit(commitLogs[len(commitLogs)-2].Id())

	gotHistory, err := suite.fileHistory.Get(middleware.NewCommit(commit))

	suite.Nil(err)
	suite.NotZero(len(gotHistory))
}

func (suite *FileHistoryTestSuite) TestGet_WhenCommitIsTheFirstCommitWithNoParents_ShouldReturnNoParentError() {
	oid, _ := git2go.NewOid("5e0c30593a8819a05135f908e1729278d3ee086e")

	suite.mockCommit.EXPECT().Id().Return(oid)
	suite.mockCommit.EXPECT().ParentCount().Return(uint(0))

	_, err := suite.fileHistory.Get(suite.mockCommit)

	suite.NotNil(err)
	suite.Equal(FileHistoryNoParentError, err)
}

func (suite *FileHistoryTestSuite) TestGet_WhenCommitTreeIsInvalid_ShouldReturnTreeError() {
	oid, _ := git2go.NewOid("5e0c30593a8819a05135f908e1729278d3ee086e")

	suite.mockCommit.EXPECT().Id().Return(oid)
	suite.mockCommit.EXPECT().ParentCount().Return(uint(1))
	suite.mockCommit.EXPECT().Parent(uint(0)).Return(suite.mockCommit)
	suite.mockCommit.EXPECT().Tree().Return(nil, errors.New("TREE_ERROR")).Times(2)

	_, err := suite.fileHistory.Get(suite.mockCommit)

	suite.NotNil(err)
	suite.Equal(FileHistoryTreeError, err)
}

func (suite *FileHistoryTestSuite) TestGet_WhenDiffTreeFails_ShouldReturnError() {
	oid, _ := git2go.NewOid("5e0c30593a8819a05135f908e1729278d3ee086e")
	treePtr := &git2go.Tree{Object: git2go.Object{}}

	suite.mockCommit.EXPECT().Id().Return(oid)
	suite.mockCommit.EXPECT().ParentCount().Return(uint(1))
	suite.mockCommit.EXPECT().Parent(uint(0)).Return(suite.mockCommit)
	suite.mockCommit.EXPECT().Tree().Return(treePtr, nil).Times(2)

	suite.mockRepo.EXPECT().DiffTreeToTree(treePtr, treePtr, nil).Return(nil, errors.New("DIFF_ERROR"))

	_, err := suite.fileHistory.Get(suite.mockCommit)

	suite.NotNil(err)
	suite.Equal(FileHistoryError, err)
}

func (suite *FileHistoryTestSuite) TestGet_WhenDiffNumDeltaIsZero_ShouldReturnError() {
	oid, _ := git2go.NewOid("5e0c30593a8819a05135f908e1729278d3ee086e")
	treePtr := &git2go.Tree{Object: git2go.Object{}}

	suite.mockCommit.EXPECT().Id().Return(oid)
	suite.mockCommit.EXPECT().ParentCount().Return(uint(1))
	suite.mockCommit.EXPECT().Parent(uint(0)).Return(suite.mockCommit)
	suite.mockCommit.EXPECT().Tree().Return(treePtr, nil).Times(2)

	suite.mockRepo.EXPECT().DiffTreeToTree(treePtr, treePtr, nil).Return(&git2go.Diff{}, nil)

	_, err := suite.fileHistory.Get(suite.mockCommit)

	suite.NotNil(err)
	suite.Equal(FileHistoryError, err)
}

func (suite *FileHistoryTestSuite) commitLogs() []git2go.Commit {
	commits, err := NewListAllLogs(suite.repo, nil, nil).Get()
	if err != nil {
		return nil
	}

	return commits
}
