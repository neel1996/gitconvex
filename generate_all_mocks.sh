#!/bin/bash

# Middleware mocks
mockgen -source=git/middleware/repository.go -destination=mocks/mock_repository.go -package=mocks
mockgen -source=git/middleware/walk.go -destination=mocks/mock_walk.go -package=mocks
mockgen -source=git/middleware/reference.go -destination=mocks/mock_reference.go -package=mocks
mockgen -source=git/middleware/index.go -destination=mocks/mock_index.go -package=mocks
mockgen -source=git/middleware/branch.go -destination=mocks/mock_branch.go -package=mocks
mockgen -source=git/middleware/commit.go -destination=mocks/mock_commit.go -package=mocks

# Commit mocks
mockgen -source=git/commit/git_list_all_commit_logs.go -destination=git/commit/mocks/mock_git_list_all_commit_logs.go -package=mocks
mockgen -source=git/commit/git_commit_file_history.go -destination=git/commit/mocks/mock_git_commit_file_history.go -package=mocks

# Remote mocks
mockgen -source=git/remote/remote_validation.go -destination=git/remote/mocks/mock_remote_validation.go -package=mocks
mockgen -source=git/remote/git_remote_list.go -destination=git/remote/mocks/mock_git_remote_list.go -package=mocks
mockgen -source=git/middleware/remotes.go -destination=git/remote/mocks/mock_remotes.go -package=mocks

# Branch mocks
mockgen -source=git/branch/git_branch_add.go  -destination=git/branch/mocks/mock_git_branch_add.go -package=mocks

# Validator mocks
mockgen -source=validator/validator.go -destination=validator/mocks/mock_validator.go -package=mocks
