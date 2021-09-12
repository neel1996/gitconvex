package branch

type Error struct {
	error
	ErrorCode   string
	ErrorString string
}

var (
	NilRepoError         = Error{ErrorCode: "REPO_NIL_ERROR", ErrorString: "Repo is nil"}
	EmptyBranchNameError = Error{ErrorCode: "BRANCH_EMPTY_ERROR", ErrorString: "Branch name(s) is empty"}
	LookupError          = Error{ErrorCode: "BRANCH_LOOKUP_ERROR", ErrorString: "Branch lookup Failed"}
	CompareError         = Error{ErrorCode: "BRANCH_COMPARE_ERROR", ErrorString: "Comparing the two branches returned no difference"}
	NilHeadError         = Error{ErrorCode: "REPO_HEAD_ERROR", ErrorString: "Can not fetch repo HEAD"}
)

func (e Error) Error() string {
	return e.ErrorString
}
