package commit

import (
	"fmt"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
)

type ListAllLogs interface {
	Get() ([]git2go.Commit, error)
}

type listAllLogs struct {
	repo      middleware.Repository
	limit     *uint
	reference *git2go.Oid
}

type commitType struct {
	commits []git2go.Commit
}

func (l listAllLogs) Get() ([]git2go.Commit, error) {
	repo := l.repo

	logItr, itrErr := repo.Walk()
	if itrErr != nil {
		logger.Log(fmt.Sprintf("Repo has no logs -> %s", itrErr.Error()), global.StatusError)
		return nil, itrErr
	}

	commits, err := l.allCommitLogs(logItr)
	if err != nil {
		logger.Log(fmt.Sprintf("Unable to obtain commits for the repo"), global.StatusError)
		return nil, err
	}

	return commits, nil
}

func (l listAllLogs) allCommitLogs(logItr middleware.RevWalk) ([]git2go.Commit, error) {
	var (
		c commitType
	)

	pushErr := l.pushCommitReference(logItr)
	if pushErr != nil {
		logger.Log(pushErr.Error(), global.StatusError)
		return nil, pushErr
	}

	err := l.iterate(logItr, &c)

	return c.commits, err
}

func (l listAllLogs) iterate(logItr middleware.RevWalk, c *commitType) error {
	if l.limit == nil {
		logger.Log("Iterating commit logs without limit", global.StatusInfo)
		return logItr.Iterate(revIterator(c))
	}

	logger.Log(fmt.Sprintf("Iterating commit logs with limit %v", *l.limit), global.StatusInfo)
	return logItr.Iterate(revIteratorWithLimit(c, *l.limit))
}

func (l listAllLogs) pushCommitReference(logItr middleware.RevWalk) error {
	if l.reference != nil {
		logger.Log(fmt.Sprintf("Iterating commits from reference %s", l.reference), global.StatusInfo)
		return logItr.Push(l.reference)
	}

	logger.Log(fmt.Sprintf("Iterating commits from HEAD"), global.StatusInfo)
	return logItr.PushHead()
}

func revIterator(c *commitType) git2go.RevWalkIterator {
	return func(commit *git2go.Commit) bool {
		if commit != nil {
			c.commits = append(c.commits, *commit)
			return true
		}

		return false
	}
}

func revIteratorWithLimit(c *commitType, limit uint) git2go.RevWalkIterator {
	var commitCounter uint

	return func(commit *git2go.Commit) bool {
		if commit != nil && commitCounter < limit {
			c.commits = append(c.commits, *commit)
			commitCounter++
			return true
		}

		return false
	}
}

func NewListAllLogs(repo middleware.Repository, limit *uint, reference *git2go.Oid) ListAllLogs {
	return listAllLogs{
		repo:      repo,
		limit:     limit,
		reference: reference,
	}
}
