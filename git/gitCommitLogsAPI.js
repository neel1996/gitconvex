const fetchRepopath = require("../global/fetchGitRepoPath");
const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

/**
 * @param  {String} repoId
 * @param  {String} skipLimit=0 - git commit --skip limit to skip previous commits
 * @returns {Object} - commit details -
 */

async function gitCommitLogHandler(repoId, skipLimit = 0) {
  const repoPath = fetchRepopath.getRepoPath(repoId);
  let commitLogLimit = 0;

  const totalCommits = await execPromisified(`git log --oneline`, {
    cwd: repoPath,
    maxBuffer: 1024 * 10240,
    windowsHide: true,
  })
    .then((res) => {
      const { stdout, stderr } = res;
      if (stdout && !stderr) {
        const gitLocalTotal = stdout.trim().split("\n").length;
        return gitLocalTotal;
      } else {
        console.log(stderr);
        return 0;
      }
    })
    .catch((err) => {
      console.log(err);
      return 0;
    });

  console.log("Total commits in the repo : ", totalCommits);

  commitLogLimit = totalCommits < 10 ? totalCommits : 10;

  if (!totalCommits) {
    return {
      hash: "",
      author: "",
      commitTime: "",
      commitMessage: "",
      commitRelativeTime: "",
      commitFilesCount: 0,
    };
  }

  return await execPromisified(
    `git log -n ${commitLogLimit} --skip ${skipLimit} --pretty=format:"%h||%an||%ad||%s" --date=short`,
    {
      cwd: repoPath,
      windowsHide: true,
      maxBuffer: 1024 * 10240,
    }
  )
    .then(async (res) => {
      const { stdout, stderr } = res;

      if (stdout && !stderr) {
        let commits = stdout.trim();

        let commitRelativeTime = await execPromisified(
          `git log -n ${commitLogLimit} --skip ${skipLimit} --pretty=format:"%ad" --date=relative`,
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
                console.log(stderr);
                return 0;
              }
            })
            .catch((err) => {
              console.log(err);
              return 0;
            });
          commit += "||" + commitFilesCount;
          return commitModel(commit);
        });
        return {
          totalCommits: totalCommits,
          commits: commitArray,
        };
      } else {
        return {
          totalCommits: 0,
          commits: [],
        };
      }
    })
    .catch((err) => {
      console.log("ERROR : Commit log collection Error!", err);
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

  const objKeys = Object.keys(commitObject);

  for (let i = 0; i < objKeys.length; i++) {
    commitObject[objKeys[i]] = commitSplit[i];
  }

  return commitObject;
}

module.exports.gitCommitLogHandler = gitCommitLogHandler;
