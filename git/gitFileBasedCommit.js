const { exec } = require("child_process");
const fs = require("fs");
const util = require("util");
const execPromised = util.promisify(exec);
const fetchRepopath = require("../global/fetchGitRepoPath");
const path = require("path");

async function gitFileBasedCommit(repoPath, fileItem) {
  return await execPromised(`git log -1 --oneline "${fileItem}"`, {
    cwd: repoPath,
    windowsHide: true,
  })
    .then(({ stderr, stdout }) => {
      if (stderr) {
        console.log(stderr);
        return "";
      }

      const splitString = stdout.split(" ");
      return splitString.slice(1, splitString.length).join(" ");
    })
    .catch((err) => {
      console.log(err);
      return "";
    });
}

module.exports.gitFileBasedCommit = gitFileBasedCommit;
