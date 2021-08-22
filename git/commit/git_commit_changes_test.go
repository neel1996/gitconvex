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
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type CommitChangesTestSuite struct {
	suite.Suite
	mockController *gomock.Controller
	repo           middleware.Repository
	noHeadRepo     middleware.Repository
	commitChanges  Changes
	commitMessage  []string
	mockRepo       *mocks.MockRepository
	mockReference  *mocks.MockReference
	mockIndex      *mocks.MockIndex
}

func TestCommitChangesTestSuite(t *testing.T) {
	suite.Run(t, new(CommitChangesTestSuite))
}

func (suite *CommitChangesTestSuite) SetupTest() {
	r, err := git2go.OpenRepository(os.Getenv("GITCONVEX_TEST_REPO"))
	if err != nil {
		fmt.Println(err)
	}

	noHeadPath := os.Getenv("GITCONVEX_TEST_REPO") + string(filepath.Separator) + "no_head_for_commit"
	noHeadRepo, _ := git2go.OpenRepository(noHeadPath)

	suite.repo = middleware.NewRepository(r)
	suite.noHeadRepo = middleware.NewRepository(noHeadRepo)

	suite.mockController = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRepository(suite.mockController)
	suite.mockReference = mocks.NewMockReference(suite.mockController)
	suite.mockIndex = mocks.NewMockIndex(suite.mockController)

	suite.commitMessage = []string{"Test commit message " + time.Now().String()}
	suite.commitChanges = NewCommitChanges(suite.mockRepo, suite.commitMessage)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenChangesAreCommitted_ShouldReturnNil() {
	suite.commitChanges = NewCommitChanges(suite.repo, suite.commitMessage)

	err := suite.commitChanges.Add()

	suite.Nil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenCommitMessageIsMultiLine_ShouldCommitToRepo() {
	suite.commitChanges = NewCommitChanges(suite.repo, []string{"Multi line commit", "Test 1"})

	err := suite.commitChanges.Add()

	suite.Nil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenRepoHasNoSignature_ShouldReturnError() {
	suite.mockRepo.EXPECT().DefaultSignature().Return(nil, errors.New("SIGNATURE_ERROR"))

	err := suite.commitChanges.Add()

	suite.NotNil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenRepoHasNoHead_ShouldCommitToHead() {
	suite.commitMessage = []string{""}
	suite.commitChanges = NewCommitChanges(suite.repo, suite.commitMessage)

	err := suite.commitChanges.Add()

	suite.Nil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenRepoHeadIsInvalid_ShouldReturnError() {
	oid, _ := git2go.NewOid("test")

	suite.mockRepo.EXPECT().DefaultSignature().Return(&git2go.Signature{}, nil)
	suite.mockRepo.EXPECT().Head().Return(suite.mockReference, nil)
	suite.mockReference.EXPECT().Target().Return(oid)
	suite.mockRepo.EXPECT().LookupCommit(oid).Return(nil, errors.New("LOOKUP_ERROR"))

	err := suite.commitChanges.Add()

	suite.NotNil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenRepoIndexIsInvalid_ShouldReturnError() {
	suite.mockRepo.EXPECT().DefaultSignature().Return(&git2go.Signature{}, nil)
	suite.mockRepo.EXPECT().Head().Return(nil, errors.New("NO_HEAD"))
	suite.mockRepo.EXPECT().Index().Return(nil, errors.New("INDEX_ERROR"))

	err := suite.commitChanges.Add()

	suite.NotNil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenRepoTreeLookupFails_ShouldReturnError() {
	oid, _ := git2go.NewOid("test")

	suite.mockRepo.EXPECT().DefaultSignature().Return(&git2go.Signature{}, nil)
	suite.mockRepo.EXPECT().Head().Return(nil, errors.New("NO_HEAD"))
	suite.mockRepo.EXPECT().Index().Return(suite.mockIndex, nil)
	suite.mockIndex.EXPECT().WriteTree().Return(oid, nil)
	suite.mockRepo.EXPECT().LookupTree(oid).Return(nil, errors.New("TREE_LOOKUP_ERROR"))

	err := suite.commitChanges.Add()

	suite.NotNil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenCommitCreation_ShouldReturnError() {
	oid, _ := git2go.NewOid("test")

	suite.mockRepo.EXPECT().DefaultSignature().Return(&git2go.Signature{}, nil)
	suite.mockRepo.EXPECT().Head().Return(nil, errors.New("NO_HEAD"))
	suite.mockRepo.EXPECT().Index().Return(suite.mockIndex, nil)
	suite.mockIndex.EXPECT().WriteTree().Return(oid, nil)
	suite.mockRepo.EXPECT().LookupTree(oid).Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CreateCommit(
		"HEAD",
		&git2go.Signature{},
		&git2go.Signature{},
		strings.Join(suite.commitMessage, ""),
		&git2go.Tree{},
	).Return(nil, errors.New("CREATE_COMMIT_ERROR"))

	err := suite.commitChanges.Add()

	suite.NotNil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenNewCommitLookupFails_ShouldReturnError() {
	oid, _ := git2go.NewOid("test")
	commitMessage := strings.Join(suite.commitMessage, "")

	suite.mockRepo.EXPECT().DefaultSignature().Return(&git2go.Signature{}, nil)
	suite.mockRepo.EXPECT().Head().Return(nil, errors.New("NO_HEAD"))
	suite.mockRepo.EXPECT().Index().Return(suite.mockIndex, nil)
	suite.mockIndex.EXPECT().WriteTree().Return(oid, nil)
	suite.mockRepo.EXPECT().LookupTree(oid).Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CreateCommit(
		"HEAD",
		&git2go.Signature{},
		&git2go.Signature{},
		commitMessage,
		&git2go.Tree{},
	).Return(oid, nil)
	suite.mockRepo.EXPECT().LookupCommit(oid).Return(nil, errors.New("LOOKUP_ERROR"))

	err := suite.commitChanges.Add()

	suite.NotNil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenNewHeadIsInvalid_ShouldReturnError() {
	oid, _ := git2go.NewOid("test")
	commitMessage := strings.Join(suite.commitMessage, "")

	suite.mockRepo.EXPECT().DefaultSignature().Return(&git2go.Signature{}, nil)
	suite.mockRepo.EXPECT().Head().Return(nil, errors.New("NO_HEAD"))
	suite.mockRepo.EXPECT().Index().Return(suite.mockIndex, nil)
	suite.mockIndex.EXPECT().WriteTree().Return(oid, nil)
	suite.mockRepo.EXPECT().LookupTree(oid).Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CreateCommit(
		"HEAD",
		&git2go.Signature{},
		&git2go.Signature{},
		commitMessage,
		&git2go.Tree{},
	).Return(oid, nil)
	suite.mockRepo.EXPECT().LookupCommit(oid).Return(nil, nil)
	suite.mockRepo.EXPECT().Head().Return(nil, errors.New("HEAD_ERR"))

	err := suite.commitChanges.Add()

	suite.NotNil(err)
}

func (suite *CommitChangesTestSuite) TestAdd_WhenSetNewHeadFails_ShouldReturnError() {
	oid, _ := git2go.NewOid("test")
	commitMessage := strings.Join(suite.commitMessage, "")

	suite.mockRepo.EXPECT().DefaultSignature().Return(&git2go.Signature{}, nil)
	suite.mockRepo.EXPECT().Head().Return(nil, errors.New("NO_HEAD"))
	suite.mockRepo.EXPECT().Index().Return(suite.mockIndex, nil)
	suite.mockIndex.EXPECT().WriteTree().Return(oid, nil)
	suite.mockRepo.EXPECT().LookupTree(oid).Return(&git2go.Tree{}, nil)
	suite.mockRepo.EXPECT().CreateCommit(
		"HEAD",
		&git2go.Signature{},
		&git2go.Signature{},
		commitMessage,
		&git2go.Tree{},
	).Return(oid, nil)
	suite.mockRepo.EXPECT().LookupCommit(oid).Return(nil, nil)
	suite.mockRepo.EXPECT().Head().Return(suite.mockReference, nil)
	suite.mockReference.EXPECT().SetTarget(oid, commitMessage).Return(nil, errors.New("SET_TARGET_ERROR"))

	err := suite.commitChanges.Add()

	suite.NotNil(err)
}
