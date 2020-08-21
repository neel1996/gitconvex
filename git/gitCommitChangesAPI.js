const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @param  {String} commitMessage
 * @returns {Object} - commits the changes with a commit message and sends the status as JSON {String}
 */

const gitCommitChangesApi = async (repoId, commitMessage) => {
  commitMessage = commitMessage.split("||").join("\n");
  commitMessage = commitMessage.replace(/"/gi, '\\"');

  return await execPromisified(`git commit -m "${commitMessage.toString()}"`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  })
    .then(({ stdout, stderr }) => {
      console.log(stdout, stderr);
      if (!stderr) {
        return "COMMIT_DONE";
      } else {
        return "COMMIT_FAILED";
      }
    })
    .catch((err) => {
      console.log(err);
      return "COMMIT_FAILED";
    });
};

module.exports.gitCommitChangesApi = gitCommitChangesApi;
