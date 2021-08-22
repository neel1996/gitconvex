package commit

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
	"strings"
)

type Mapper interface {
	Map([]git2go.Commit) []*model.GitCommits
}

type mapper struct {
	fileHistory FileHistory
}

type mapperCommitType struct {
	hash       *string
	author     *string
	commitTime *string
	message    *string
	filesCount *int
}

func (m mapper) Map(commits []git2go.Commit) []*model.GitCommits {
	var mappedCommits []*model.GitCommits

	if len(commits) == 0 {
		logger.Log("No commits received for mapping", global.StatusWarning)
		return nil
	}

	for _, commit := range commits {
		commitFields := m.commitFields(commit)

		mappedCommits = append(mappedCommits, &model.GitCommits{
			Hash:             commitFields.hash,
			Author:           commitFields.author,
			CommitTime:       commitFields.commitTime,
			CommitMessage:    commitFields.message,
			CommitFilesCount: commitFields.filesCount,
		})
	}

	return mappedCommits
}

func (m mapper) commitFields(commit git2go.Commit) mapperCommitType {
	commitHash := commit.Id().String()
	commitAuthor := strings.Split(commit.Author().Name, " ")[0]
	commitMessage := strings.Split(commit.Message(), "\n")[0]
	commitDate := commit.Committer().When.String()
	commitFileHistory := m.commitFileHistory(commit)

	commitFields := mapperCommitType{
		hash:       &commitHash,
		author:     &commitAuthor,
		commitTime: &commitDate,
		message:    &commitMessage,
		filesCount: &commitFileHistory,
	}
	return commitFields
}

func (m mapper) commitFileHistory(commit git2go.Commit) int {
	c := middleware.NewCommit(&commit)
	fileHistory, err := m.fileHistory.Get(c)
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return 0
	}

	return len(fileHistory)
}

func NewMapper(commitFile FileHistory) Mapper {
	return mapper{
		fileHistory: commitFile,
	}
}
