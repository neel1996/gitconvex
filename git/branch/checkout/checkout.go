package checkout

import "github.com/neel1996/gitconvex/global"

var logger global.Logger

type BranchDetails struct {
	BranchName          string
	ReferenceBranchName string
	RemoteBranchName    string
}

type Checkout interface {
	CheckoutBranch() error
	GenerateBranchFields() BranchDetails
	LogAndReturnError(err error) error
}
