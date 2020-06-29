const { healthCheckHandler } = require("../API/healthcheckApi");
const { fetchRepoHandler } = require("../API/fetchRepoApi");
const { addRepoHandler } = require("../API/addRepoApi");
const { getGitRepoStatus } = require("../git/gitRepoAPI");
const { gitTrackedDiff } = require("../git/gitTrackedDiff");
const { gitFileDifferenceHandler } = require("../git/gitFileDifferenceAPI");
const { gitCommitLogHandler } = require("../git/gitCommitLogsAPI");
const { getStagedFiles } = require("../git/gitGetStagedFilesAPI");
const { gitGetUnpushedCommits } = require("../git/gitGetUnpushedCommits");
const { gitSetBranchApi } = require("../git/gitSetBranch.js");
const { gitStageAllItemsApi } = require("../git/gitStageAllItemsAPI");
const { gitCommitChangesApi } = require("../git/gitCommitChangesAPI");
const { gitPushToRemoteApi } = require("../git/gitPushToRemoteAPI");
const { gitStageItem } = require("../git/gitStageItem");
const { fetchDatabaseFile, fetchRepoDetails } = require("../API/settingsApi");
const {
  gitRemoveAllStagedItemApi,
  gitRemoveStagedItemApi,
} = require("../git/gitRemoveStagedItems");
const { gitFetchApi, gitPullApi } = require("../git/gitFetchPullApi");
const { deleteRepoApi } = require("../API/deleteRepoApi");
const { gitAddBranchApi } = require("../git/gitAddBranchApi");

module.exports.healthCheckFunction = healthCheckFunction = async (payload) => {
  const hcPayload = await healthCheckHandler().then((res) => res);
  const { osCheck, gitCheck, nodeCheck } = JSON.parse(
    JSON.stringify(hcPayload)
  );
  return {
    healthCheck: {
      osCheck,
      gitCheck,
      nodeCheck,
    },
  };
};

module.exports.fetchRepoFunction = fetchRepoFunction = async (payload) => {
  const repoFetchPayload = await fetchRepoHandler().then((res) => res);

  const { repoId, repoName, repoPath } = repoFetchPayload;

  console.log(repoFetchPayload);

  return {
    fetchRepo: {
      repoId,
      repoName,
      repoPath,
    },
  };
};

module.exports.addRepoFunction = addRepoFunction = async (parsedPayload) => {
  const { repoName, repoPath, initCheck } = JSON.parse(parsedPayload);

  if (repoName && repoPath) {
    console.log(parsedPayload);

    return await addRepoHandler(repoName, repoPath, initCheck);
  } else {
    return {
      message: "REPO_WRITE_FAILURE",
    };
  }
};

module.exports.gitCommitLogsFunction = gitCommitLogsFunction = async (
  parsedPayload
) => {
  const { repoId } = JSON.parse(parsedPayload);
  if (repoId) {
    console.log(await gitCommitLogHandler(repoId));
    return {
      gitCommitLogs: {
        ...(await gitCommitLogHandler(repoId)),
      },
    };
  } else {
    return {
      gitCommitLogs: {
        commits: [],
      },
    };
  }
};

module.exports.repoDetailsFunction = repoDetailsFunction = async (
  parsedPayload
) => {
  const repoDetails = await getGitRepoStatus(JSON.parse(parsedPayload).repoId);
  return {
    gitRepoStatus: {
      ...repoDetails,
    },
  };
};

module.exports.gitChangeTrackerFunction = gitChangeTrackerFunction = async (
  parsedPayload
) => {
  let { repoId } = JSON.parse(parsedPayload);
  const gitChangeResults = await gitTrackedDiff(repoId);
  return {
    gitChanges: {
      ...gitChangeResults,
    },
  };
};

module.exports.gitFileDiffFunction = gitFileDiffFunction = async (
  parsedPayload
) => {
  let fileDiffArgs = JSON.parse(parsedPayload);
  const gitFileLineChanges = await gitFileDifferenceHandler(
    fileDiffArgs.repoId,
    fileDiffArgs.fileName
  ).then((res) => res);
  console.log(gitFileLineChanges);
  return {
    gitFileLineChanges: {
      ...gitFileLineChanges,
    },
  };
};

module.exports.gitGetStagedFiles = gitGetStagedFiles = async (payload) => {
  const { repoId } = JSON.parse(payload);
  const stagedFiles = await getStagedFiles(repoId);

  console.log(stagedFiles);

  if (stagedFiles) {
    return {
      gitStagedFiles: {
        stagedFiles: [...stagedFiles],
      },
    };
  } else {
    return {
      gitStagedFiles: {
        stagedFiles: [],
      },
    };
  }
};

module.exports.gitUnpushedCommits = gitUnpushedCommits = async (payload) => {
  const { repoId, remoteName } = JSON.parse(payload);
  const unPushedCommits = await gitGetUnpushedCommits(repoId, remoteName);

  if (unPushedCommits) {
    return {
      gitUnpushedCommits: {
        commits: unPushedCommits,
      },
    };
  } else {
    return {
      gitUnpushedCommits: {
        commits: [],
      },
    };
  }
};

module.exports.gitSetBranch = gitSetBranch = async (repoId, branch) => {
  return await gitSetBranchApi(repoId, branch);
};

module.exports.gitStageAllItems = gitStageAllItems = async (repoId) => {
  return await gitStageAllItemsApi(repoId);
};

module.exports.gitCommitChanges = gitCommitChanges = async (
  repoId,
  commitMessage
) => {
  return await gitCommitChangesApi(repoId, commitMessage);
};

module.exports.gitPushToRemote = gitPushToRemote = async (
  repoId,
  remoteHost,
  branch
) => {
  return await gitPushToRemoteApi(repoId, remoteHost, branch);
};

module.exports.gitStageItem = gitStageItemApi = async (repoId, item) => {
  return await gitStageItem(repoId, item);
};

module.exports.settingsFetchDbPath = settingsFetchDbPath = async () => {
  return await fetchDatabaseFile();
};

module.exports.settingsFetchRepoDetails = settingsFetchRepoDetails = async () => {
  return await fetchRepoDetails();
};

module.exports.gitRemoveStagedItem = gitRemoveStagedItem = async (
  repoId,
  item
) => {
  return await gitRemoveStagedItemApi(repoId, item);
};

module.exports.gitRemoveAllStagedItems = gitRemoveAllStagedItems = async (
  repoId
) => {
  return await gitRemoveAllStagedItemApi(repoId);
};

module.exports.gitFetchFromRemote = gitFetchFromRemote = async (repoId) => {
  return await gitFetchApi(repoId);
};

module.exports.gitPullFromRemote = gitPullFromRemote = async (repoId) => {
  return await gitPullApi(repoId);
};

module.exports.deleteRepo = deleteRepo = async (
  repoId,
  name,
  pathName,
  time
) => {
  return await deleteRepoApi(repoId, name, pathName, time);
};

module.exports.addBranch = addBranch = async (repoId, branchName) => {
  return await gitAddBranchApi(repoId, branchName);
};
