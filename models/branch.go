package models

import (
	"github.com/neel1996/gitconvex/git/branch"
	"github.com/neel1996/gitconvex/git/branch/checkout"
)

type Branch struct {
	Add      branch.Add
	Checkout checkout.Factory
	Compare  branch.Compare
	Delete   branch.Delete
	List     branch.List
}
