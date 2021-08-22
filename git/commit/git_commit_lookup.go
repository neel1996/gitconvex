package commit

import (
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
)

type Lookup interface {
	WithReferenceId(string) (middleware.Commit, error)
}

type lookup struct {
	repo middleware.Repository
}

func (l lookup) WithReferenceId(commitHash string) (middleware.Commit, error) {
	oid, oidErr := git2go.NewOid(commitHash)
	if oidErr != nil {
		logger.Log(oidErr.Error(), global.StatusError)
		return nil, OidConversionError
	}

	commit, err := l.repo.LookupCommit(oid)
	if err != nil {
		logger.Log(err.Error(), global.StatusError)
		return nil, CommitLookupError
	}

	return middleware.NewCommit(commit), nil
}

func NewLookup(repo middleware.Repository) Lookup {
	return lookup{repo: repo}
}
