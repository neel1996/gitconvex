const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @param  {String} remoteHost
 * @param  {String} branch
 * @returns {String} - status of the git push
 */

const gitPushToRemoteApi = async (repoId, remoteHost, branch) => {
  return await execPromisified(`git remote -v`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  })
    .then(async ({ stdout, stderr }) => {
      if (!stderr) {
        const remoteVerbose = stdout.trim().split("\n");

        const filteredRemote = remoteVerbose
          .filter((remote) => {
            if (remote.includes(remoteHost) && remote.includes("push")) {
              return true;
            }
            return false;
          })
          .join("");

        const remoteName = filteredRemote.trim().split(/\s/gi)[0];
        branch = branch.trim();

        if (branch.match(/[^a-zA-Z0-9-_.:~@$^/\\s\\n\\r]/gi)) {
          console.log(
            new Error("Invalid remote branch string!"),
            branch,
            branch.match(/[^a-zA-Z0-9-_.:~@$^/\\s\\r\\n]/gi)
          );
          return "PUSH_FAILED";
        }

        const pushCommand = `git push -u "${remoteName}" "${branch}"`;

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
