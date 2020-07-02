const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitAddBranchApi = async (repoId, branchName) => {
  return await execPromisified(`git checkout -b ${branchName}`, {
    cwd: fetchRepopath.getRepoPath(repoId),
  })
    .then((res) => {
      console.log(res);
      return "BRANCH_CREATION_SUCCESS";
    })
    .catch((err) => {
      console.log(err);
      return "BRANCH_ADD_FAILED";
    });
};

module.exports.gitAddBranchApi = gitAddBranchApi;
