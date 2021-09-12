package checkout

import (
	"fmt"
	"github.com/neel1996/gitconvex/git/middleware"
	"github.com/neel1996/gitconvex/global"
	"github.com/neel1996/gitconvex/validator"
	"strings"
)

type Factory interface {
	GetCheckoutAction(branchName string) Checkout
}

type factory struct {
	repo            middleware.Repository
	repoValidator   validator.Validator
	branchValidator validator.ValidatorWithStringFields
}

func (f factory) GetCheckoutAction(branchName string) Checkout {
	logger.Log(fmt.Sprintf("Received branch %s", branchName), global.StatusInfo)
	if repoValidationErr := f.repoValidator.Validate(f.repo); repoValidationErr != nil {
		logger.Log(repoValidationErr.Error(), global.StatusError)
		return nil
	}

	if validationErr := f.validateBranchFields(branchName); validationErr != nil {
		return nil
	}

	if strings.Contains(branchName, "remotes/") {
		return NewCheckoutRemoteBranch(f.repo, branchName, nil)
	} else {
		return NewCheckOutLocalBranch(f.repo, branchName)
	}
}

func (f factory) validateBranchFields(branchName string) error {
	logger.Log("Validating branch fields before checkout", global.StatusInfo)

	if err := f.branchValidator.ValidateWithFields(branchName); err != nil {
		logger.Log(err.Error(), global.StatusError)
		return err
	}
	return nil
}

func NewCheckoutFactory(repo middleware.Repository, repoValidator validator.Validator, branchValidator validator.ValidatorWithStringFields) Factory {
	return factory{
		repo:            repo,
		repoValidator:   repoValidator,
		branchValidator: branchValidator,
	}
}
