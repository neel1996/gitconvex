const { healthCheckHandler } = require("../API/healthcheckApi");
const { fetchRepoHandler } = require("../API/fetchRepoApi");
const { addRepoHandler } = require("../API/addRepoApi");
const { gitCommitLogDbSerchApi } = require("../API/commitLogSearchApi");
const { getGitRepoStatus } = require("../git/gitRepoAPI");
const { gitTrackedDiff } = require("../git/gitTrackedDiff");
const { gitFileDifferenceHandler } = require("../git/gitFileDifferenceAPI");
const { gitCommitLogHandler } = require("../git/gitCommitLogsAPI");
const { gitCommitFileApi } = require("../git/gitCommitFilesApi");
const { getStagedFiles } = require("../git/gitGetStagedFilesAPI");
const { gitGetUnpushedCommits } = require("../git/gitGetUnpushedCommits");
const { gitSetBranchApi } = require("../git/gitSetBranch.js");
const { gitStageAllItemsApi } = require("../git/gitStageAllItemsAPI");
const { gitCommitChangesApi } = require("../git/gitCommitChangesAPI");
const { gitPushToRemoteApi } = require("../git/gitPushToRemoteAPI");
const { gitStageItem } = require("../git/gitStageItem");
const { gitAddRemoteApi } = require("../git/gitAddRemoteApi");
const { gitDeleteBranchApi } = require("../git/gitBranchDeleteApi");
const { gitFetchFolderContentApi } = require("../git/gitFolderDetailsApi");
const { codeFileViewApi } = require("../API/codeFileViewApi");
const { branchCompareApi } = require("../API/branchCompareApi");
const {
  fetchDatabaseFile,
  fetchRepoDetails,
  updateDbFile,
  updatePortDetails,
  getPortDetails,
} = require("../API/settingsApi");
const {
  gitRemoveAllStagedItemApi,
  gitRemoveStagedItemApi,
} = require("../git/gitRemoveStagedItems");
const { gitFetchApi, gitPullApi } = require("../git/gitFetchPullApi");
const { deleteRepoApi } = require("../API/deleteRepoApi");
const { gitAddBranchApi } = require("../git/gitAddBranchApi");

/**
 * @returns {Object} - platform and required software related info
 */

module.exports.healthCheckFunction = healthCheckFunction = async () => {
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

/**
 * @returns {Object} - details of all the stored repos
 */

module.exports.fetchRepoFunction = fetchRepoFunction = async () => {
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

/**
 * @param  {String} repoName - name of the repo to be stored
 * @param  {String} repoPath - path where the repo resides
 * @param  {Boolean} initCheck - boolean switch to enable git init if the repo is not a git repo
 * @param  {Boolean} cloneCheck - boolean switch to enable git clone if the repo needs to be pulled from a remote repo
 * @param  {String} cloneUrl - URL of the remote repo if the clone switch is true
 * @returns {Object} - returns object with status of the repo write operation
 */

module.exports.addRepoFunction = addRepoFunction = async (
  repoName,
  repoPath,
  initCheck,
  cloneCheck,
  cloneUrl
) => {
  if (repoName && repoPath) {
    return await addRepoHandler(
      repoName,
      repoPath,
      initCheck,
      cloneCheck,
      cloneUrl
    );
  } else {
    return {
      message: "REPO_WRITE_FAILURE",
    };
  }
};

module.exports.codeFileViewFunction = codeFileViewFunction = async (
  parsedPayload
) => {
  const { repoId, fileItem } = JSON.parse(parsedPayload);
  return await codeFileViewApi(repoId, fileItem);
};

/**
 * @param  {String} parsedPayload - strigified json with repoId: String and skipLimit: number
 * @returns {Object} - containing staggered commits and total commits available in the repo
 */

module.exports.gitCommitLogsFunction = gitCommitLogsFunction = async (
  parsedPayload
) => {
  const { repoId, skipLimit } = JSON.parse(parsedPayload);
  if (repoId) {
    return {
      gitCommitLogs: {
        ...(await gitCommitLogHandler(repoId, skipLimit)),
      },
    };
  } else {
    return {
      gitCommitLogs: {
        totalCommits: 0,
        commits: [],
      },
    };
  }
};

/**
 * @param  {String} repoId - Unique repo ID
 * @param  {String} searchType - Search category
 * @param  {String} searchKey - Key for searching
 */

module.exports.gitCommitLogSearchFunction = gitCommitLogSearchFunction = async (
  repoId,
  searchType,
  searchKey
) => {
  return await gitCommitLogDbSerchApi(repoId, searchType, searchKey).catch(
    (err) => {
      console.log(err);
      return [];
    }
  );
};

/**
 * @param  {String} parsedPayload - strigified json holding the repoId: String and commitHash: String
 * @returns {Object} - object holding all the files which are changed in the commit
 */

module.exports.gitCommitFileFunction = gitCommitFileFunction = async (
  parsedPayload
) => {
  const { repoId, commitHash } = JSON.parse(parsedPayload);
  console.log(await gitCommitFileApi(repoId, commitHash));
  return {
    gitCommitFiles: [...(await gitCommitFileApi(repoId, commitHash))],
  };
};

/**
 * @param  {String} parsedPayload - strigified json holding repoId
 * @returns {Object} - object holding the details of the repo
 */

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

/**
 * @param  {String} parsedPayload - strigified json holding repoId
 * @returns {Object} - files / folders within the selected directory and the latest commit messages of the content
 */

module.exports.gitFolderContentApi = gitFolderContentApi = async (
  parsedPayload
) => {
  const { repoId, directoryName } = JSON.parse(parsedPayload);
  return await gitFetchFolderContentApi(repoId, directoryName);
};

/**
 * @param  {String} parsedPayload - strigified json holding repoId
 * @returns {Object} - returns the files that have changed from HEAD
 */

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

/**
 * @param  {String} parsedPayload - strigified json holding repoId and the file name for which the git diff needs to be found
 * @returns {Object} - returns the git diff related data and the line changes within the file
 */

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

/**
 * @param  {String} parsedPayload - strigified json holding repoId
 * @returns {Object} - list of files that are staged
 */

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

/**
 * @param  {String} parsedPayload - strigified json holding repoId and the name of the user selected remote repo
 * @returns {Object} - list of commits that are waiting to be pushed
 */

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

/**
 * @param  {String} repoId
 * @param  {String} branch
 */

module.exports.gitSetBranch = gitSetBranch = async (repoId, branch) => {
  return await gitSetBranchApi(repoId, branch);
};

/**
 * @param  {String} repoId
 */

module.exports.gitStageAllItems = gitStageAllItems = async (repoId) => {
  return await gitStageAllItemsApi(repoId);
};

/**
 * @param  {String} repoId
 * @param  {String} commitMessage
 */

module.exports.gitCommitChanges = gitCommitChanges = async (
  repoId,
  commitMessage
) => {
  return await gitCommitChangesApi(repoId, commitMessage);
};

/**
 * @param  {String} repoId
 * @param  {String} remoteHost
 * @param  {String} branch
 */

module.exports.gitPushToRemote = gitPushToRemote = async (
  repoId,
  remoteHost,
  branch
) => {
  return await gitPushToRemoteApi(repoId, remoteHost, branch);
};

/**
 * @param  {String} repoId
 * @param  {String} item
 */

module.exports.gitStageItem = gitStageItemApi = async (repoId, item) => {
  return await gitStageItem(repoId, item);
};

module.exports.settingsFetchDbPath = settingsFetchDbPath = async () => {
  return await fetchDatabaseFile();
};

module.exports.settingsGetPortDetails = settingsGetPortDetails = async () => {
  return await getPortDetails();
};

/**
 * @param  {number} newPort
 */

module.exports.settingsUpdatePortDetail = settingsUpdatePortDetail = async (
  newPort
) => {
  return await updatePortDetails(newPort);
};

module.exports.settingsFetchRepoDetails = settingsFetchRepoDetails = async () => {
  return await fetchRepoDetails();
};

/**
 * @param  {String} repoId
 * @param  {String} item
 */

module.exports.gitRemoveStagedItem = gitRemoveStagedItem = async (
  repoId,
  item
) => {
  return await gitRemoveStagedItemApi(repoId, item);
};

/**
 * @param  {String} repoId
 */

module.exports.gitRemoveAllStagedItems = gitRemoveAllStagedItems = async (
  repoId
) => {
  return await gitRemoveAllStagedItemApi(repoId);
};

/**
 * @param  {String} repoId
 * @param  {String} remoteUrl=""
 * @param  {String} remoteBranch=""
 */

module.exports.gitFetchFromRemote = gitFetchFromRemote = async (
  repoId,
  remoteUrl = "",
  remoteBranch = ""
) => {
  return await gitFetchApi(repoId, remoteUrl, remoteBranch);
};

/**
 * @param  {String} repoId
 * @param  {String} remoteUrl
 * @param  {String} remoteBranch
 */

module.exports.gitPullFromRemote = gitPullFromRemote = async (
  repoId,
  remoteUrl,
  remoteBranch
) => {
  return await gitPullApi(repoId, remoteUrl, remoteBranch);
};

/**
 * @param  {String} repoId
 * @param  {String} branchName
 * @param  {Boolean} forceFlag
 */

module.exports.gitDeleteBranchApi = gitDeleteBranchFunction = async (
  repoId,
  branchName,
  forceFlag
) => {
  return gitDeleteBranchApi(repoId, branchName, forceFlag);
};

/**
 * @param  {String} repoId
 */

module.exports.deleteRepo = deleteRepo = async (repoId) => {
  return await deleteRepoApi(repoId);
};

/**
 * @param  {String} repoId
 * @param  {String} branchName
 */

module.exports.addBranch = addBranch = async (repoId, branchName) => {
  return await gitAddBranchApi(repoId, branchName);
};

/**
 * @param  {String} fileName
 */

module.exports.updateDbFileApi = updateDbFileApi = async (fileName) => {
  return await updateDbFile(fileName);
};

/**
 * @param  {String} repoId
 * @param  {String} remoteName
 * @param  {String} remoteUrl
 */

module.exports.gitAddRemoteRepoApi = gitAddRemoteRepoApi = async (
  repoId,
  remoteName,
  remoteUrl
) => {
  return await gitAddRemoteApi(repoId, remoteName, remoteUrl);
};

module.exports.branchCompareApi = branchCompareFunction = async (payload) => {
  const { repoId, baseBranch, compareBranch } = JSON.parse(payload);
  return await {
    branchCompare: branchCompareApi(repoId, baseBranch, compareBranch),
  };
};
