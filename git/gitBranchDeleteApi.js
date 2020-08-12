const { exec } = require("child_process");
const fs = require("fs");
const util = require("util");
const execPromised = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

/**
 * @param  {String} repoId
 * @param  {String} branchName
 * @param  {Boolean} forceFlag
 * @returns {Object} - deletes a branch and sends the status as JSON {String}
 */

const gitDeleteBranchApi = async (repoId, branchName, forceFlag) => {
  try {
    branchName = branchName.trim();
    if (branchName.match(/[^a-zA-Z0-9-_.:~@$^/\\s\\r\\n]/gi)) {
      throw new Error("Invalid branchName string!");
    }

    let commandString = "";
    if (forceFlag) {
      commandString = `git branch -D "${branchName}"`;
    } else {
      commandString = `git branch -d "${branchName}"`;
    }

    return execPromised(commandString, {
      cwd: fetchRepopath.getRepoPath(repoId),
      windowsHide: true,
    })
      .then(({ stdout, stderr }) => {
        if (stdout && !stderr) {
          console.log(stdout);
          return {
            status: "BRANCH_DELETE_SUCCESS",
            deletionResponse: stdout.trim().split("\n"),
          };
        } else {
          console.log(stderr);
          if (stderr.includes("Error")) {
            return {
              status: "BRANCH_DELETE_FAILED",
            };
          } else {
            return {
              status: "BRANCH_DELETE_SUCCESS",
            };
          }
        }
      })
      .catch((err) => {
        console.log(err);
        return {
          status: "BRANCH_DELETE_FAILED",
        };
      });
  } catch (err) {
    return {
      status: "BRANCH_DELETE_FAILED",
    };
  }
};

module.exports.gitDeleteBranchApi = gitDeleteBranchApi;
