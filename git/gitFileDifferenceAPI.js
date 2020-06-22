const fetchRepopath = require("../global/fetchGitRepoPath");
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
  return await execPromisified(
    `cd ${repoPath}; git diff --stat ${fileName} && echo "SPLIT___LINE" && git diff -U$(wc -l ${fileName} | xargs)`
  ).then((res) => {
    const { stdout, stderr } = res;

    if (stdout && !stderr) {
      var splitLines = stdout.split("SPLIT___LINE");
      var diffStat = splitLines[0].trim().split("\n");
      var fileDiff = splitLines[1].trim().split("\n");

      return {
        diffStat,
        fileDiff,
      };
    }
  });
}

module.exports.gitFileDifferenceHandler = gitFileDifferenceHandler;
