const { buildSchema, GraphQLString } = require("graphql");

const globalAPISchema = new buildSchema(
  `
        type GitConvexAPI{
            gitConvexApi(route: String!, payload: String) : GitConvexResults!
        }

        type GitConvexResults{
            healthCheck: healthCheckResults!
            fetchRepo: fetchRepoResults!
            gitRepoStatus: gitRepoStatusResults!
            gitChanges: gitChangeResults!
            gitFileLineChanges: gitFileLineChangeResults!
            gitCommitLogs: gitCommitLogResults!
            gitCommitFiles: [commitFileType]!
            gitStagedFiles: gitStagedFileResults!
            gitUnpushedCommits: gitUnpushedCommitResults!
            settingsDatabasePath: String!
            settingsPortDetails: Int!
            settingsRepoDetails: [settingsFetchRepoResults]!
        }

        type healthCheckResults{
            osCheck: String!
            gitCheck: String!
            nodeCheck: String!
        }

        type fetchRepoResults{
            repoId: [String]
            repoPath: [String]
            repoName: [String]
        }
        
        type gitCommits{
            hash: String,
            author: String,
            commitTime: String,
            commitMessage: String,
            commitRelativeTime: String
            commitFilesCount: Int
        }

        type commitFileType{
            type: String!
            fileName: String!
        }

        type gitUnpushedCommitResults{
            commits: [String]
        }

        type gitCommitLogResults{
            commits: [gitCommits]
        }

        type gitRepoStatusResults{
            gitRemoteData: String
            gitRepoName: String
            gitBranchList: [String]
            gitAllBranchList: [String]
            gitCurrentBranch: String
            gitRemoteHost: String
            gitTotalCommits: Int
            gitLatestCommit: String
            gitTrackedFiles: [String]
            gitFileBasedCommit: [String]
            gitTotalTrackedFiles: Int
        }
        
        type gitChangeResults{
            gitUntrackedFiles: [String]
            gitChangedFiles: [String]
        }

        type gitFileLineChangeResults{
            diffStat: [String]
            fileDiff: [String]
        }

        type gitStagedFileResults{
            stagedFiles: [String]
        }

        type gitUnpushedCommitResult{
            gitUnpushedCommits: [String]
        }

        type settingsFetchRepoResults{
            id: String!
            repoName: String!
            timeStamp: String!
            repoPath: String!
        }

        type GitConvexMutation{
            addRepo(repoName: String!, repoPath: String!, initSwitch: Boolean!, cloneSwitch: Boolean!, cloneUrl: String): addRepoStatus!
            setBranch(repoId: String!, branch: String!): String!
            stageItem(repoId: String!, item: String!): String!
            stageAllItems(repoId: String): String!
            commitChanges(repoId: String!, commitMessage: String!): String!
            pushToRemote(repoId: String!, remoteHost: String!, branch: String!): String!
            settingsEditDbPath(newPath: String!): String!
            settingsEditPort(newPort: Int!): String!
            settingsDeleteRepo(repoId: String!): String!
            removeStagedItem(repoId: String!, item: String!): String!
            removeAllStagedItem(repoId: String!): String!
            addBranch(repoId: String!, branchName: String!): String!
            fetchFromRemote(repoId: String!, remoteUrl: String, remoteBranch: String): gitFetchStatus!
            pullFromRemote(repoId: String!, remoteUrl: String!, remoteBranch: String!): gitPullStatus!
            updateRepoDataFile(newDbFile: String!): String!
            deleteRepo(repoId: String!): deleteRepoStatus!
            addRemoteRepo(repoId: String!, remoteName: String!, remoteUrl: String!): String!
            deleteBranch(repoId: String!, branchName: String!, forceFlag: Boolean!): deleteBranchStatus!
        }

        type addRepoStatus{
            message: String!,
            repoId: String
        }

        type gitFetchStatus{
            status: String!
            fetchedItems: [String]
        }

        type gitPullStatus{
            status: String!
            pulledItems: [String]
        }

        type deleteRepoStatus{
            status: String!
            repoId: String
        }

        type deleteBranchStatus{
            status: String!,
            deletionResponse: String
        }

        schema{
            query: GitConvexAPI
            mutation: GitConvexMutation
        }

    `
);

module.exports = globalAPISchema;
