const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @param  {String} remoteName
 * @param  {String} remoteUrl
 * @returns {Object} - Adds a new remote and sends the status as JSON {String}
 */

const gitAddRemoteApi = async (repoId, remoteName, remoteUrl) => {
  try {
    if (
      remoteName.match(/[^a-zA-Z0-9-_.]/gi) ||
      remoteUrl.match(/[^a-zA-Z0-9-_.~@#$%:/]/gi)
    ) {
      throw new Error("Input string remoteName or remoteURL invalid");
    }

    return await execPromisified(
      `git remote add "${remoteName}" "${remoteUrl}"`,
      {
        cwd: fetchRepopath.getRepoPath(repoId),
        windowsHide: true,
      }
    )
      .then((res) => {
        console.log(res);
        return "REMOTE_ADD_SUCCESS";
      })
      .catch((err) => {
        console.log(err);
        return "REMOTE_ADD_FAILED";
      });
  } catch (err) {
    console.log(err);
    return "REMOTE_ADD_FAILED";
  }
};

module.exports.gitAddRemoteApi = gitAddRemoteApi;
