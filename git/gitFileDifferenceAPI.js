const fetchRepopath = require("../global/fetchGitRepoPath");
const fs = require("fs");
const { exec } = require("child_process");
const util = require("util");
const { LangLine } = require("@itassistors/langline");
const execPromisified = util.promisify(exec);

/**
 * @param  {String} repoId
 * @param  {String} fileName
 */

async function gitFileDifferenceHandler(repoId, fileName) {
  if (repoId && fileName) {
    var differencePayload = await getGitFileDifference(repoId, fileName);
    return differencePayload;
  }
}

/**
 * @param  {String} repoId
 * @param  {String} fileName
 * @returns {Object} - git diff status and the lines which of the files with the change indicator
 */

async function getGitFileDifference(repoId, fileName) {
  const repoPath = fetchRepopath.getRepoPath(repoId);

  const fileContentLength = await fs.promises
    .readFile(repoPath + "/" + fileName)
    .then((data) => {
      const interData = data.toString().split("\n");
      return interData.length;
    })
    .catch((err) => {
      console.log(err);
    });

  try {
    if (fileContentLength <= 0) {
      throw new Error("Invalid file selection!");
    }

    const diffStat = await execPromisified(`git diff --stat "${fileName}"`, {
      cwd: repoPath,
      windowsHide: true,
    })
      .then(({ stdout, stderr }) => {
        if (stdout) {
          return stdout.trim().split("\n");
        } else {
          console.log(stderr);
          return ["NO_STAT"];
        }
      })
      .catch((err) => {
        console.log(err);
        return ["NO_STAT"];
      });

    const fileDiff = await execPromisified(
      `git diff -U${fileContentLength} "${fileName}"`,
      {
        cwd: repoPath,
        windowsHide: true,
      }
    )
      .then(({ stdout, stderr }) => {
        if (stdout) {
          return stdout.trim().split("\n");
        } else {
          console.log(stderr);
          return ["NO_DIFF"];
        }
      })
      .catch((err) => {
        console.log(err);
        return ["NO_DIFF"];
      });

    const { prismIndicator } = new LangLine().withFileName(fileName);

    return {
      diffStat,
      fileDiff,
      language: prismIndicator ? prismIndicator : "markdown",
    };
  } catch (err) {
    console.log(err);
    return ["NO_DIFF"];
  }
}

module.exports.gitFileDifferenceHandler = gitFileDifferenceHandler;
