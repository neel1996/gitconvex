const { gitBranchCompare } = require("../git/gitBranchCompare");
const { getRepoPath } = require("../global/fetchGitRepoPath");

async function branchCompareApi(repoId, baseBranch, compareBranch) {
  const repoPath = getRepoPath(repoId);

  if (baseBranch !== compareBranch) {
    return await gitBranchCompare(repoPath, baseBranch, compareBranch);
  } else {
    return {
      message: "Nothing to compare as the branches are the same",
    };
  }
}

module.exports.branchCompareApi = branchCompareApi;
