const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @param  {String} branchName
 * @returns {String} - adds a new branch and sends the status as JSON {String}
 */

const gitAddBranchApi = async (repoId, branchName) => {
  try {
    branchName = branchName.trim();
    if (branchName.match(/[^a-zA-Z0-9-_.:~@$^/\\s\\r\\n]/gi)) {
      throw new Error("Invalid branch name string");
    }

    return await execPromisified(`git checkout -b "${branchName}"`, {
      cwd: fetchRepopath.getRepoPath(repoId),
      windowsHide: true,
    })
      .then((res) => {
        console.log(res);
        return "BRANCH_CREATION_SUCCESS";
      })
      .catch((err) => {
        console.log(err);
        return "BRANCH_ADD_FAILED";
      });
  } catch (err) {
    console.log(err);
    return "BRANCH_ADD_FAILED";
  }
};

module.exports.gitAddBranchApi = gitAddBranchApi;
