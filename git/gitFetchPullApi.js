const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitFetchApi = async (repoId) => {
  return await execPromisified(`git fetch`, {
    cwd: fetchRepopath.getRepoPath(repoId),
  })
    .then(({ stdout, stderr }) => {
      if (stdout && !stderr) {
        const fetchResponse = stdout.trim().split("\n");
        if (fetchResponse && fetchResponse.length > 0) {
          return {
            status: "FETCH_PRESENT",
            fetchedItems: fetchResponse,
          };
        } else {
          return {
            status: "FETCH_ABSENT",
          };
        }
      } else {
        return {
          status: "FETCH_ABSENT",
        };
      }
    })
    .catch((err) => {
      console.log(err);
      if (err) {
        return {
          status: "FETCH_ERROR",
        };
      }
    });
};

const gitPullApi = async (repoId) => {
  return await execPromisified(`git pull`, {
    cwd: fetchRepopath.getRepoPath(repoId),
  })
    .then(async ({ stdout, stderr }) => {
      if (stdout && !stderr) {
        const pullResponse = stdout.trim().split("\n");

        if (pullResponse && pullResponse.length > 0) {
          return {
            status: "PULL_SUCCESS",
            pulledItems: pullResponse,
          };
        } else {
          return {
            status: "PULL_EMPTY",
          };
        }
      } else {
        return {
          status: "PULL_FAILED",
        };
      }
    })
    .catch((err) => {
      console.log(err);
      return {
        status: "PULL_FAILED",
      };
    });
};

module.exports.gitFetchApi = gitFetchApi;
module.exports.gitPullApi = gitPullApi;
