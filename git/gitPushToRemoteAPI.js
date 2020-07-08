const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitPushToRemoteApi = async (repoId, remoteHost, branch) => {
  return await execPromisified(
    `git remote -v | grep '${remoteHost}' | grep 'push'`,
    { cwd: fetchRepopath.getRepoPath(repoId), windowsHide: true }
  )
    .then(async ({ stdout, stderr }) => {
      if (!stderr) {
        const remoteName = stdout.trim().split(/\s/gi)[0];
        const pushCommand = `git push -u ${remoteName} ${branch}`;

        return await execPromisified(`${pushCommand}`, {
          cwd: fetchRepopath.getRepoPath(repoId),
        })
          .then((stdout, stderr) => {
            if (!stderr && stdout) {
              return "PUSH_DONE";
            } else {
              console.log(stderr);
            }
          })
          .catch((err) => {
            console.log(err);
            return "PUSH_FAILED";
          });
      } else {
        console.log(stderr);
        return "PUSH_FAILED";
      }
    })
    .catch((err) => {
      console.log(err);
      return "PUSH_FAILED";
    });
};

module.exports.gitPushToRemoteApi = gitPushToRemoteApi;
