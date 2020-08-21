const { exec } = require("child_process");
const fs = require("fs");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @param  {String} item
 * @returns {String} - status of the unstaging action
 * @description - Removed the selected item from staged area
 */

const gitRemoveStagedItemApi = async (repoId, item) => {
  const repopath = fetchRepopath.getRepoPath(repoId);

  const fileItemValid = await fs.promises
    .stat(repopath + "/" + item)
    .then((res) => res.isFile());

  if (!fileItemValid) {
    console.log("Invalid item string");
    return "STAGE_REMOVE_FAILED";
  }

  return await execPromisified(`git reset "${item}"`, {
    cwd: repopath,
    windowsHide: true,
  })
    .then(({ stdout, stderr }) => {
      if (stderr) {
        console.log(stderr);
        return "STAGE_REMOVE_FAILED";
      } else {
        console.log(stdout);
        return "STAGE_REMOVE_SUCCESS";
      }
    })
    .catch((err) => {
      console.log(err);
      return "STAGE_REMOVE_FAILED";
    });
};

/**
 * @param  {String} repoId
 * @returns {String}
 * @description - removes all the the staged items
 */

const gitRemoveAllStagedItemApi = async (repoId) => {
  return await execPromisified(`git reset`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  })
    .then(({ stdout, stderr }) => {
      if (stderr) {
        console.log(stderr);
        return "STAGE_ALL_REMOVE_FAILED";
      } else {
        console.log(stdout);
        return "STAGE_ALL_REMOVE_SUCCESS";
      }
    })
    .catch((err) => {
      console.log(err);
      return "STAGE_ALL_REMOVE_FAILED";
    });
};

module.exports.gitRemoveStagedItemApi = gitRemoveStagedItemApi;
module.exports.gitRemoveAllStagedItemApi = gitRemoveAllStagedItemApi;
