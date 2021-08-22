package remote

import (
	"errors"
	git2go "github.com/libgit2/git2go/v31"
	"github.com/neel1996/gitconvex/git/middleware"
)

type Validation interface {
	ValidateRemoteFields(remoteFields ...string) error
}

type validation struct {
	repo         middleware.Repository
	remoteFields []string
}

func (v validation) ValidateRemoteFields(remoteFields ...string) error {
	v.remoteFields = remoteFields
	validateRepoErr := v.validateRepo()
	if validateRepoErr != nil {
		return validateRepoErr
	}

	validateRemotesErr := v.validateRemoteCollection()
	if validateRemotesErr != nil {
		return validateRemotesErr
	}

	fieldValidationErr := v.validateRemoteFields()
	if fieldValidationErr != nil {
		return fieldValidationErr
	}

	return nil
}

func (v validation) validateRemoteFields() error {
	for _, field := range v.remoteFields {
		if field == "" {
			return errors.New("one or more remote fields are empty")
		}
	}
	return nil
}

func (v validation) validateRepo() error {
	if v.repo == nil {
		return errors.New("repo is nil")
	}

	return nil
}

func (v validation) validateRemoteCollection() error {
	if v.repo.Remotes().Get() == (git2go.RemoteCollection{}) {
		return errors.New("remote collection is nil")
	}

	return nil
}

func NewRemoteValidation(repo middleware.Repository) Validation {
	return validation{
		repo: repo,
	}
}
