const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @param  {String} remoteUrl
 * @returns {String} - name of the remote based on the remote URL
 */

const getRemoteName = async (repoId, remoteUrl) => {
  return await execPromisified(`git remote -v`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  })
    .then(({ stdout, stderr }) => {
      if (stdout && !stderr) {
        console.log(stdout);
        const localName = stdout.trim().split("\n");
        return localName
          .filter((item) => {
            if (item.includes(remoteUrl) && item.includes("fetch")) {
              return true;
            }
          })
          .join()
          .split(/\s/gi)[0];
      } else {
        console.log(stderr);
        return "";
      }
    })
    .catch((err) => {
      console.log(err);
      return "";
    });
};

/**
 * @param  {String} repoId
 * @param  {String} remoteUrl=""
 * @param  {String} remoteBranch=""
 * @returns {Object: {status: String, fetchedItems: Array[String]}} - performs a git fetch and returns the status along with the fetched changes
 */

const gitFetchApi = async (repoId, remoteUrl = "", remoteBranch = "") => {
  const remoteName = await getRemoteName(repoId, remoteUrl);
  console.log("Selected remote name : ", remoteName);

  if (!remoteName) {
    console.log("NO REMOTE MATCHING THE URL");

    return {
      status: "FETCH_ERROR",
    };
  }

  return await execPromisified(`git fetch ${remoteName} ${remoteBranch} -v`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  })
    .then(({ stdout, stderr }) => {
      if (stdout || stderr) {
        // Git fetch alone returns the result in the standard error stream
        let responseValue = "";
        if (stdout) {
          responseValue += stdout;
        }
        if (stderr) {
          responseValue += stderr;
        }

        const fetchResponse = responseValue.trim().split("\n");

        console.log("Fetch Response :" + fetchResponse);
        if (fetchResponse) {
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

/**
 * @param  {String} repoId
 * @param  {String} remoteUrl
 * @param  {String} remoteBranch
 * @returns {Object: {status: String, pulledItems: Array[String]}} - performs a git pull from the remote and returns the pulled changes
 */

const gitPullApi = async (repoId, remoteUrl, remoteBranch) => {
  const remoteName = await getRemoteName(repoId, remoteUrl);
  console.log("Selected remote name : ", remoteName);

  if (!remoteName) {
    console.log("NO REMOTE MATCHING THE URL");

    return {
      status: "PULL_ERROR",
    };
  }

  return await execPromisified(`git pull ${remoteName} ${remoteBranch} -v`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  })
    .then(async ({ stdout, stderr }) => {
      if (stdout || stderr) {
        let responseValue = "";
        if (stdout) {
          responseValue += stdout;
        }
        if (stderr) {
          responseValue += stderr;
        }
        const pullResponse = responseValue.trim().split("\n");

        if (pullResponse && pullResponse.length > 0) {
          return {
            status: "PULL_SUCCESS",
            pulledItems: pullResponse,
          };
        } else {
          return {
            status: "PULL_ABSENT",
          };
        }
      } else {
        return {
          status: "PULL_ERROR",
        };
      }
    })
    .catch((err) => {
      console.log(err);
      return {
        status: "PULL_ERROR",
      };
    });
};

module.exports.gitFetchApi = gitFetchApi;
module.exports.gitPullApi = gitPullApi;
