const {
  exec
} = require("child_process");
const util = require("util");
const execPromosified = util.promisify(exec);
const fetchRepoPath = require("../global/fetchGitRepoPath");

async function gitTrackedDiff(repoId) {
  const repoPath = fetchRepoPath.getRepoPath(repoId);

  var responseObject = {
    gitChangedFiles: await getGitDiff(repoPath).then(res => res),
    gitUntrackedFiles: await getUntrackedFiles(repoPath).then(res => res),
  };

  console.log(responseObject)
  return responseObject;
}

async function getGitDiff(repoPath) {
  return await execPromosified(`cd ${repoPath}; git diff --raw`).then((res) => {
    const {
      stdout,
      stderr
    } = res;
    var parsedEntry = stdout.trim().split("\n");

    var gitDifference = parsedEntry.map((entry) => {
      if (entry.split(/\s+/)) {
        let splitEntry = entry.split(/\s+/);
        if (splitEntry[4] && splitEntry[5]) {
          return "" + splitEntry[4] + "," + splitEntry[5];
        }
      }
    });

    return gitDifference.filter((entry) => (entry ? entry : ""));
  });
}

async function getUntrackedFiles(repoPath) {
  return await execPromosified(
    `cd ${repoPath}; git ls-files --others --exclude-standard`
  ).then((res) => {
    const {
      stdout,
      stderr
    } = res;
    var parsedEntry = stdout
      .trim()
      .split("\n")
      .filter((item) => {
        if (item) {
          return item;
        }
      });

    var gitUntrackedFiles = parsedEntry.map((entry) => {
      let fileDirArray = [];
      if (entry.includes("/")) {
        let splitEntry = entry.split("/");
        let dirEntry = splitEntry.map((elm, index) => {
          if (index < entry.split("/").length - 1) {
            return elm;
          }
        });
        return `${dirEntry.join("/")},${splitEntry[splitEntry.length - 1]}`;
      } else {
        return `NO_DIR,${entry}`;
      }
    });
    return gitUntrackedFiles;
  });
}

module.exports.gitTrackedDiff = gitTrackedDiff;