const { gitCommitCompare } = require("../git/gitCommitCompareApi");
const { getRepoPath } = require("../global/fetchGitRepoPath");

async function commitCompareApi(repoId, baseCommit, compareCommit) {
  const repoPath = getRepoPath(repoId);

  if (baseCommit !== compareCommit) {
    return await gitCommitCompare(repoPath, baseCommit, compareCommit);
  } else {
    return {
      message: "Nothing to compare as the commits are the same",
    };
  }
}

module.exports.commitCompareApi = commitCompareApi;
