const { exec } = require("child_process");
const util = require("util");

const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @returns {Array} - list of all the staged items
 */

const getStagedFiles = async (repoId) => {
  return await execPromisified(`git diff --name-only --cached`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  })
    .then((res) => {
      const { stdout, stderr } = res;

      if (!stderr && stdout) {
        const stagedFiles = stdout.trim().split("\n");
        console.log(stagedFiles);
        return stagedFiles;
      } else {
        console.log(stderr);
        return [];
      }
    })
    .catch((err) => {
      console.log(err);
      return [];
    });
};

module.exports.getStagedFiles = getStagedFiles;
