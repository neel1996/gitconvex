package search

import (
	git2go "github.com/libgit2/git2go/v31"
)

type Type int

const (
	CommitHash = iota + 1
	CommitMessage
	CommitAuthor
)

var typeMap = map[string]Type{
	"hash":    CommitHash,
	"message": CommitMessage,
	"author":  CommitAuthor,
}

func GetSearchAction(searchType string, commits []git2go.Commit) Search {
	mappedSearchType := typeMap[searchType]

	switch mappedSearchType {
	case CommitHash:
		return commitHashSearch{
			commits: commits,
		}
	case CommitMessage:
		return commitMessageSearch{
			commits: commits,
		}
	case CommitAuthor:
		return commitAuthorSearch{
			commits: commits,
		}
	default:
		return nil
	}
}
