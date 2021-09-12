package validator

import "github.com/neel1996/gitconvex/git/middleware"

type repoValidator struct {
}

func (v repoValidator) Validate(r interface{}) error {
	if r == nil {
		return NilRepoError
	}

	if repo, ok := r.(middleware.Repository); !ok || repo == nil {
		return NilRepoError
	}
	return nil
}

func NewRepoValidator() Validator {
	return repoValidator{}
}
