const fetchRepopath = require("../global/fetchGitRepoPath");
const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

async function gitCommitLogHandler(repoId) {
  const repoPath = fetchRepopath.getRepoPath(repoId);
  return await execPromisified(
    `git log --pretty=format:"%h||%an||%ad||%s" --date=short`,
    {
      cwd: repoPath,
      windowsHide: true,
    }
  )
    .then(async (res) => {
      const { stdout, stderr } = res;

      if (stdout && !stderr) {
        let commits = stdout.trim();

        let commitRelativeTime = await execPromisified(
          `git log --pretty=format:"%ad" --date=relative`,
          { cwd: repoPath, windowsHide: true }
        )
          .then(({ stdout, stderr }) => {
            if (stdout) {
              return stdout.trim().split("\n");
            } else {
              console.log(stderr);
              return [];
            }
          })
          .catch((err) => {
            console.log(err);
            return [];
          });

        let commitArray = commits.split("\n").map(async (commit, index) => {
          commit += "||" + commitRelativeTime[index];
          const commitFilesCount = await execPromisified(
            `git diff-tree --no-commit-id --name-only -r ${
              commit.split("||")[0]
            }`,
            { cwd: repoPath, windowsHide: true }
          )
            .then(({ stdout, stderr }) => {
              if (stdout) {
                return stdout.trim().split("\n").length;
              } else {
                console.log("Error occurred!");
                return 0;
              }
            })
            .catch((err) => {
              console.log("Error occurred!");
              return 0;
            });
          commit += "||" + commitFilesCount;
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
    })
    .catch((err) => {
      console.log("ERROR : Commit log collection Error!");
      return {
        hash: "",
        author: "",
        commitTime: "",
        commitMessage: "",
        commitRelativeTime: "",
        commitFilesCount: 0,
      };
    });
}

function commitModel(commit) {
  let commitObject = {
    hash: "",
    author: "",
    commitTime: "",
    commitMessage: "",
    commitRelativeTime: "",
    commitFilesCount: 0,
  };

  let commitSplit = commit.split("||");

  commitObject.hash = commitSplit[0];
  commitObject.author = commitSplit[1];
  commitObject.commitTime = commitSplit[2];
  commitObject.commitMessage = commitSplit[3];
  commitObject.commitRelativeTime = commitSplit[4];
  commitObject.commitFilesCount = commitSplit[5];

  return commitObject;
}

module.exports.gitCommitLogHandler = gitCommitLogHandler;
