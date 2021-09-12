package graph

import (
	"context"
	"github.com/neel1996/gitconvex/controller"
	"github.com/neel1996/gitconvex/git/middleware"
	initialize "github.com/neel1996/gitconvex/init"
	"github.com/neel1996/gitconvex/models"
	"github.com/neel1996/gitconvex/validator"
)

type Resolver struct{}

func (r *Resolver) BranchController(ctx context.Context, repo middleware.Repository) controller.BranchController {
	repoValidator := validator.NewRepoValidator()
	if repoValidator.Validate(repo) != nil {
		ctx.Done()
		return nil
	}
	branchObject := initialize.NewInitBranch(repo)

	branchModel := models.Branch{
		Add:      branchObject.BranchAdd(),
		Checkout: branchObject.BranchCheckout(),
		Delete:   branchObject.BranchDelete(),
	}

	return controller.NewBranchController(branchModel)
}
