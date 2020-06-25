const { exec } = require("child_process");

const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitRemoveStagedItemApi = async (repoId, item) => {
  return await execPromisified(
    `cd ${fetchRepopath.getRepoPath(repoId)}; git reset ${item}`
  ).then(({ stdout, stderr }) => {
    if (stderr) {
      console.log(stderr);
      return "STAGE_REMOVE_FAILED";
    } else {
      console.log(stdout);
      return "STAGE_REMOVE_SUCCESS";
    }
  });
};

const gitRemoveAllStagedItemApi = async (repoId) => {
  return await execPromisified(
    `cd ${fetchRepopath.getRepoPath(repoId)}; git reset`
  ).then(({ stdout, stderr }) => {
    if (stderr) {
      console.log(stderr);
      return "STAGE_ALL_REMOVE_FAILED";
    } else {
      console.log(stdout);
      return "STAGE_ALL_REMOVE_SUCCESS";
    }
  });
};

module.exports.gitRemoveStagedItemApi = gitRemoveStagedItemApi;
module.exports.gitRemoveAllStagedItemApi = gitRemoveAllStagedItemApi;
