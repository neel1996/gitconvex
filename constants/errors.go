package constants

type ApiError struct {
	error
	ErrorCode   string
	ErrorString string
}

var (
	BranchAddError       = ApiError{ErrorCode: "BRANCH_ADD_FAILED", ErrorString: "Failed to add new branch"}
	BranchCheckoutError  = ApiError{ErrorCode: "CHECKOUT_FAILED", ErrorString: "Failed to checkout branch"}
	BranchDeleteError    = ApiError{ErrorCode: "BRANCH_DELETE_FAILED", ErrorString: "Failed to delete branch"}
	RemoteAddError       = ApiError{ErrorCode: "REMOTE_ADD_FAILED", ErrorString: "Failed to delete add new remote"}
	RemoteDeleteError    = ApiError{ErrorCode: "REMOTE_DELETE_FAILED", ErrorString: "Failed to delete add new remote"}
	RemoteEditError      = ApiError{ErrorCode: "REMOTE_EDIT_FAILED", ErrorString: "Failed to delete add new remote"}
	FetchFromRemoteError = ApiError{ErrorCode: "FETCH_ERROR", ErrorString: "Failed to delete add new remote"}
	PullFromRemoteError  = ApiError{ErrorCode: "PULL_ERROR", ErrorString: "Failed to delete add new remote"}
	StageItemError       = ApiError{ErrorCode: "ADD_ITEM_FAILED", ErrorString: "Failed to delete add new remote"}
	StageAllItemsError   = ApiError{ErrorCode: "ALL_STAGE_FAILED", ErrorString: "Failed to delete add new remote"}
	RemoveItemError      = ApiError{ErrorCode: "STAGE_REMOVE_FAILED", ErrorString: "Failed to delete add new remote"}
	RemoveAllItemsError  = ApiError{ErrorCode: "STAGE_ALL_REMOVE_FAILED", ErrorString: "Failed to delete add new remote"}
	PushToRemoteError    = ApiError{ErrorCode: "PUSH_FAILED", ErrorString: "Failed to delete add new remote"}
	PortUpdateError      = ApiError{ErrorCode: "PORT_UPDATE_FAILED", ErrorString: "Failed to delete add new remote"}
	DataFileUpdateError  = ApiError{ErrorCode: "DATAFILE_UPDATE_FAILED", ErrorString: "Failed to delete add new remote"}
	RepoDeleteError      = ApiError{ErrorCode: "DELETE_FAILED", ErrorString: "Failed to delete add new remote"}
)

func (e ApiError) Error() string {
	return e.ErrorString
}
