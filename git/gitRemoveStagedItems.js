const { exec } = require("child_process");

const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitRemoveStagedItemApi = async (repoId, item) => {
  return await execPromisified(`git reset ${item}`, {
    cwd: fetchRepopath.getRepoPath(repoId),
  }).then(({ stdout, stderr }) => {
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
  return await execPromisified(`git reset`, {
    cwd: fetchRepopath.getRepoPath(repoId),
  }).then(({ stdout, stderr }) => {
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
