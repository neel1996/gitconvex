const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitStageAllItemsApi = async (repoId) => {
  return await execPromisified(
    `cd ${fetchRepopath.getRepoPath(repoId)}; git add --all`
  )
    .then(({ stdout, stderr }) => {
      if (!stderr) {
        return "ALL_STAGED";
      }
    })
    .catch((err) => {
      return "ERR_STAGE_ALL";
    });
};

module.exports.gitStageAllItemsApi = gitStageAllItemsApi;
