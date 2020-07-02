const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitCommitChangesApi = async (repoId, commitMessage) => {
  commitMessage = commitMessage.split("||").join("\n");

  return await execPromisified(`git commit -m "${commitMessage}"`, {
    cwd: fetchRepopath.getRepoPath(repoId),
  })
    .then(({ stdout, stderr }) => {
      if (!stderr) {
        return "COMMIT_DONE";
      } else {
        return "COMMIT_FAILED";
      }
    })
    .catch((err) => {
      return "COMMIT_FAILED";
    });
};

module.exports.gitCommitChangesApi = gitCommitChangesApi;
