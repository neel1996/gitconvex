const fetchRepopath = require("../global/fetchGitRepoPath");
const fs = require("fs");
const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

async function gitFileDifferenceHandler(repoId, fileName) {
  if (repoId && fileName) {
    var differencePayload = await getGitFileDifference(repoId, fileName);

    return differencePayload;
  }
}

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

  const diffStat = await execPromisified(`git diff --stat ${fileName}`, {
    cwd: repoPath,
    windowsHide: true,
  })
    .then(({ stdout, stderr }) => {
      if (stdout && !stderr) {
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
    `git diff -U${fileContentLength} ${fileName}`,
    {
      cwd: repoPath,
      windowsHide: true,
    }
  )
    .then(({ stdout, stderr }) => {
      if (stdout && !stderr) {
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

  return {
    diffStat,
    fileDiff,
  };
}

module.exports.gitFileDifferenceHandler = gitFileDifferenceHandler;
