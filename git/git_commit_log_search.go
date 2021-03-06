package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"go/types"
	"regexp"
)

type SearchCommitInterface interface {
	SearchCommitLogs() []*model.GitCommits
}

type SearchCommitStruct struct {
	Repo       *git.Repository
	SearchType string
	SearchKey  string
}

// SearchCommitLogs searches for the required commits matching the searchKey with the respective searchType
//
// The UI uses 'message' as the default search type to look up based on commit messages
func (s SearchCommitStruct) SearchCommitLogs() []*model.GitCommits {
	var searchResult []*model.GitCommits
	logger := global.Logger{}

	repo := s.Repo
	searchKey := s.SearchKey
	searchType := s.SearchType

	commitLogs, _ := repo.Log(&git.LogOptions{
		Order: git.LogOrderDFSPost,
		All:   true,
	})

	logger.Log(fmt.Sprintf("Searching commit logs with %s for -> %s", searchType, searchKey), global.StatusInfo)

	_ = commitLogs.ForEach(func(commit *object.Commit) error {
		if len(searchResult) == 10 {
			return types.Error{Msg: "Commit limit exceeded"}
		}

		switch searchType {
		case "message":
			if isMatch, _ := regexp.MatchString(searchKey, commit.Message); isMatch {
				commitLog := commitOrganizer([]object.Commit{*commit})
				searchResult = append(searchResult, commitLog...)
			}
			break
		case "hash":
			if isMatch, _ := regexp.MatchString(searchKey, commit.Hash.String()); isMatch {
				commitLog := commitOrganizer([]object.Commit{*commit})
				searchResult = append(searchResult, commitLog...)
			}
			break
		case "user":
			if isMatch, _ := regexp.MatchString(searchKey, commit.Author.Name); isMatch {
				commitLog := commitOrganizer([]object.Commit{*commit})
				searchResult = append(searchResult, commitLog...)
			}
			break
		}
		return nil
	})

	commitLogs.Close()
	return searchResult
}
