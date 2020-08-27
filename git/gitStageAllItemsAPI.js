const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @description - stages all changed items
 */

const gitStageAllItemsApi = async (repoId) => {
  return await execPromisified(`git add --all`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  })
    .then(({ stdout, stderr }) => {
      if (!stderr) {
        return "ALL_STAGED";
      } else if (stderr.match(/warning: LF will be replaced by CRLF/gi)) {
        console.log("WARN:", stderr);
        return "ALL_STAGED";
      } else {
        console.log("ERROR:", stderr);
        return "ERR_STAGE_ALL";
      }
    })
    .catch((err) => {
      console.log("ERROR:", err);
      return "ERR_STAGE_ALL";
    });
};

module.exports.gitStageAllItemsApi = gitStageAllItemsApi;
