package git

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex-server/global"
	"github.com/neel1996/gitconvex-server/graph/model"
	"regexp"
)

type SearchCommitInterface interface {
	SearchCommitLogs() []*model.GitCommits
}

type SearchCommitStruct struct {
	Repo       *git2go.Repository
	SearchType string
	SearchKey  string
}

// SearchCommitLogs searches for the required commits matching the searchKey with the respective searchType
//
// The UI uses 'message' as the default search type to look up based on commit messages
func (s SearchCommitStruct) SearchCommitLogs() []*model.GitCommits {
	var searchResult []*model.GitCommits

	repo := s.Repo

	searchKey := s.SearchKey
	searchType := s.SearchType

	walker, walkerErr := repo.Walk()

	if walkerErr == nil {
		logger.Log(fmt.Sprintf("Searching commit logs with %s for -> %s", searchType, searchKey), global.StatusInfo)
		_ = walker.PushHead()
		_ = walker.Iterate(func(commit *git2go.Commit) bool {
			if len(searchResult) == 10 {
				return false
			}

			switch searchType {
			case "message":
				if isMatch, _ := regexp.MatchString(searchKey, commit.Message()); isMatch {
					commitLog := commitOrganizer(repo, []git2go.Commit{*commit})
					searchResult = append(searchResult, commitLog...)
				}
				break
			case "hash":
				if isMatch, _ := regexp.MatchString(searchKey, commit.Id().String()); isMatch {
					commitLog := commitOrganizer(repo, []git2go.Commit{*commit})
					searchResult = append(searchResult, commitLog...)
				}
				break
			case "user":
				if isMatch, _ := regexp.MatchString(searchKey, commit.Author().Name); isMatch {
					commitLog := commitOrganizer(repo, []git2go.Commit{*commit})
					searchResult = append(searchResult, commitLog...)
				}
				break
			}
			return true
		})
	} else {
		logger.Log(walkerErr.Error(), global.StatusError)
		return searchResult
	}
	return searchResult
}
