package validator

type Error struct {
	error
	ErrorCode   string
	ErrorString string
}

var (
	NilRepoError    = Error{ErrorCode: "REPO_NIL_ERROR", ErrorString: "Repo is nil"}
	EmptyBranchName = Error{ErrorCode: "EMPTY_BRANCH_NAME", ErrorString: "Branch name is empty"}
)

func (e Error) Error() string {
	return e.ErrorString
}
