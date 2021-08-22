package commit

type Error struct {
	error
	ErrorCode   string
	ErrorString string
}

var (
	ChangeError                = Error{ErrorCode: "COMMIT_FAILED", ErrorString: "Failed to commit the new changes"}
	LogsError                  = Error{ErrorCode: "LIST_COMMIT_LOGS_FAILED", ErrorString: "Failed to list the commit logs for the repo"}
	FileHistoryError           = Error{ErrorCode: "LISTING_COMMIT_FILE_HISTORY_FAILED", ErrorString: "Failed to get the list of file changes for the commit"}
	FileHistoryNoParentError   = Error{ErrorCode: "LISTING_COMMIT_FILE_HISTORY_NO_PARENT", ErrorString: "The HEAD commit is the only commit in the repo and has no previous histories"}
	FileHistoryTreeError       = Error{ErrorCode: "LISTING_COMMIT_FILE_HISTORY_INVALID_TREE", ErrorString: "The commit tree is invalid"}
	InvalidSearchCategoryError = Error{ErrorCode: "COMMIT_SEARCH_INVALID_SEARCH_TYPE", ErrorString: "The selected search category is invalid"}
	EmptyCommitHashError       = Error{ErrorCode: "COMMIT_HASH_REQUIRED", ErrorString: "Commit hash cannot be empty"}
	OidConversionError         = Error{ErrorCode: "INVALID_HASH_FOR_OID_CONVERSION", ErrorString: "Commit hash cannot be converted into an OID for commit lookup"}
	CommitLookupError          = Error{ErrorCode: "COMMIT_LOOKUP_FAILED", ErrorString: "Unable to find commit with the OID"}
)

func (e Error) Error() string {
	return e.ErrorString
}
