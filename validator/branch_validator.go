package validator

type branchValidator struct {
}

func (v branchValidator) ValidateWithFields(branchNames ...string) error {
	if len(branchNames) == 0 {
		return EmptyBranchName
	}

	for _, branchName := range branchNames {
		if branchName == "" {
			return EmptyBranchName
		}
	}

	return nil
}

func NewBranchValidator() ValidatorWithStringFields {
	return branchValidator{}
}
