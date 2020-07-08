const fetchRepopath = require("../global/fetchGitRepoPath");
const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

async function gitCommitLogHandler(repoId) {
  const repoPath = fetchRepopath.getRepoPath(repoId);
  return await execPromisified(`git log --pretty=format:"%h||%an||%ad||%s"`, {
    cwd: repoPath,
    windowsHide: true,
  }).then((res) => {
    const { stdout, stderr } = res;

    if (stdout && !stderr) {
      let commits = stdout.trim();
      let commitArray = commits.split("\n").map((commit) => {
        return commitModel(commit);
      });
      return {
        commits: commitArray,
      };
    } else {
      return {
        commits: [],
      };
    }
  });
}

function commitModel(commit) {
  let commitObject = {
    hash: "",
    author: "",
    commitTime: "",
    commitMessage: "",
  };

  let commitSplit = commit.split("||");

  commitObject.hash = commitSplit[0];
  commitObject.author = commitSplit[1];
  commitObject.commitTime = commitSplit[2];
  commitObject.commitMessage = commitSplit[3];

  return commitObject;
}

module.exports.gitCommitLogHandler = gitCommitLogHandler;
