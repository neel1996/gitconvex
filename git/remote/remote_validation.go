package remote

import (
	"errors"
	git2go "github.com/libgit2/git2go/v31"
)

type Validation interface {
	ValidateRemoteFields() error
}

type validation struct {
	repo         *git2go.Repository
	remoteFields []string
}

func (v validation) ValidateRemoteFields() error {
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
	if v.repo.Remotes == (git2go.RemoteCollection{}) {
		return errors.New("remote collection is nil")
	}

	return nil
}

func NewRemoteValidation(repo *git2go.Repository, remoteFields ...string) Validation {
	return validation{
		repo:         repo,
		remoteFields: remoteFields,
	}
}
