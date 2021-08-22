package commit

import (
	"fmt"
	"github.com/neel1996/gitconvex/git/commit/search"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/graph/model"
)

type SearchLogs interface {
	GetMatchingLogs(searchType string, searchKey string) ([]*model.GitCommits, error)
}

type searchLogs struct {
	repo           middleware.Repository
	listCommitLogs ListAllLogs
	mapper         Mapper
}

// GetMatchingLogs searches for the required commits matching the searchKey with the respective searchType
//
// The UI uses 'message' as the default search type to look up based on commit messages
func (s searchLogs) GetMatchingLogs(searchType string, searchKey string) ([]*model.GitCommits, error) {
	logger.Log(fmt.Sprintf("Searching commit logs filtered with - %s for -> %s", searchType, searchKey), global.StatusInfo)

	commits, listErr := s.listCommitLogs.Get()
	if listErr != nil {
		logger.Log("Unable to fetch commit logs : "+listErr.Error(), global.StatusError)
		return nil, listErr
	}

	searchFactory := search.GetSearchAction(searchType, commits)
	if searchFactory == nil {
		logger.Log("Invalid search category", global.StatusError)
		return nil, InvalidSearchCategoryError
	}

	return s.mapper.Map(searchFactory.Search(searchKey)), nil
}

func NewSearchLogs(repo middleware.Repository, listCommitLogs ListAllLogs, mapper Mapper) SearchLogs {
	return searchLogs{
		repo:           repo,
		listCommitLogs: listCommitLogs,
		mapper:         mapper,
	}
}
