const { exec } = require("child_process");
const util = require("util");

const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @param  {String} branch
 * @description - switches the current branch
 */

const gitSetBranchApi = async (repoId, branch) => {
  branch = branch.trim();
  
  if (branch.match(/[^a-zA-Z0-9-_.:~@$^/\\s\\r\\n]/gi)) {
    console.log("Invalid branch string!");
    return "BRANCH_SET_FAILED";
  }

  return await execPromisified(`git checkout ${branch}`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  }).then(({ stdout, stderr }) => {
    if (stderr) {
      return "BRANCH_SET_FAILED";
    } else {
      if (stdout) {
        console.log(stdout);
        return "BRANCH_SET_SUCCESS";
      }
    }
  });
};

module.exports.gitSetBranchApi = gitSetBranchApi;
