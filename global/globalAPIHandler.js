const express = require("express");
const app = express();
const cors = require("cors");
app.use(cors());

const {
  HEALTH_CHECK,
  FETCH_REPO,
  ADD_REPO,
  REPO_DETAILS,
  REPO_TRACKED_DIFF,
  REPO_FILE_DIFF,
  COMMIT_LOGS,
  GIT_STAGED_FILES,
  GIT_UNPUSHED_COMMITS,
  SETTINGS_DBPATH,
  SETTINGS_REPODETAILS,
} = require("./globalRouteStore");

const graphqlHTTP = require("express-graphql");

const globalGQLSchema = require("../global/gqlGlobalAPISchema");

const {
  addRepoFunction,
  fetchRepoFunction,
  gitChangeTrackerFunction,
  gitFileDiffFunction,
  healthCheckFunction,
  repoDetailsFunction,
  gitCommitLogsFunction,
  gitGetStagedFiles,
  gitUnpushedCommits,
  gitSetBranch,
  gitStageAllItems,
  gitCommitChanges,
  gitPushToRemote,
  gitStageItem,
  settingsFetchDbPath,
  settingsFetchRepoDetails,
  gitRemoveStagedItem,
  gitRemoveAllStagedItems,
  gitFetchFromRemote,
  gitPullFromRemote,
  deleteRepo,
} = require("./globalFunctionStore");

app.use(
  "/gitconvexapi",
  graphqlHTTP({
    schema: globalGQLSchema,
    graphiql: true,
    rootValue: {
      gitConvexApi: async (args) => {
        const { route, payload } = args;

        console.log("Api Route : " + route + "\nAPI Payload : " + payload);

        let parsedPayload = {};

        if (payload) {
          parsedPayload = JSON.parse(JSON.stringify(payload));
        }

        switch (route) {
          case HEALTH_CHECK:
            return healthCheckFunction();
          case FETCH_REPO:
            return fetchRepoFunction();
          case ADD_REPO:
            return addRepoFunction(parsedPayload);
          case REPO_DETAILS:
            return repoDetailsFunction(parsedPayload);
          case REPO_TRACKED_DIFF:
            return gitChangeTrackerFunction(parsedPayload);
          case REPO_FILE_DIFF:
            return gitFileDiffFunction(parsedPayload);
          case COMMIT_LOGS:
            return gitCommitLogsFunction(parsedPayload);
          case GIT_STAGED_FILES:
            return gitGetStagedFiles(parsedPayload);
          case GIT_UNPUSHED_COMMITS:
            return gitUnpushedCommits(parsedPayload);
          case SETTINGS_DBPATH:
            return settingsFetchDbPath();
          case SETTINGS_REPODETAILS:
            return settingsFetchRepoDetails();
          default:
            return { message: "Query Termination" };
        }
      },
      setBranch: async (args) => {
        const { repoId, branch } = args;
        console.log(args);
        return await gitSetBranch(repoId, branch);
      },
      stageAllItems: async (args) => {
        const { repoId } = args;
        return await gitStageAllItems(repoId);
      },
      commitChanges: async (args) => {
        const { repoId, commitMessage } = args;
        return await gitCommitChanges(repoId, commitMessage);
      },
      pushToRemote: async (args) => {
        const { repoId, remoteHost, branch } = args;
        return await gitPushToRemote(repoId, remoteHost, branch);
      },
      stageItem: async (args) => {
        const { repoId, item } = args;
        return await gitStageItem(repoId, item);
      },
      removeStagedItem: async (args) => {
        const { repoId, item } = args;
        return await gitRemoveStagedItem(repoId, item);
      },
      removeAllStagedItem: async (args) => {
        const { repoId } = args;
        return await gitRemoveAllStagedItems(repoId);
      },
      fetchFromRemote: async (args) => {
        const { repoId } = args;
        return await gitFetchFromRemote(repoId);
      },
      pullFromRemote: async (args) => {
        const { repoId } = args;
        return await gitPullFromRemote(repoId);
      },
      deleteRepo: async (args) => {
        const { repoId, name, pathName, time } = args;
        return await deleteRepo(repoId, name, pathName, time);
      },
    },
  })
);

module.exports = app;
