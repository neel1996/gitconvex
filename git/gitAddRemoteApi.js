const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitAddRemoteApi = async (repoId, remoteName, remoteUrl) => {
  return await execPromisified(`git remote add ${remoteName} ${remoteUrl}`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  })
    .then((res) => {
      console.log(res);
      return "REMOTE_ADD_SUCCESS";
    })
    .catch((err) => {
      console.log(err);
      return "REMOTE_ADD_FAILED";
    });
};

module.exports.gitAddRemoteApi = gitAddRemoteApi;
