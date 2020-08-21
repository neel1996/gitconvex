const fetchRepopath = require("../global/fetchGitRepoPath");
const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

/**
 * @param  {String} repoId
 * @param  {Array.<Object></Object>} hashArray - Array of commit hashes
 * @returns {Array.<Object>} - commit details
 */

async function gitCommitLogSearchHandler(repoId, hashArray) {
  const repoPath = fetchRepopath.getRepoPath(repoId);

  if (!repoId || hashArray.length <= 0) {
    return {
      hash: "",
      author: "",
      commitTime: "",
      commitMessage: "",
      commitRelativeTime: "",
      commitFilesCount: 0,
    };
  }

  let commitSearchResults = hashArray.map(async (hash) => {
    return await execPromisified(
      `git log -1 ${hash} --pretty=format:"%h||%an||%ad||%s" --date=short`,
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

          const commitRelativeTime = await execPromisified(
            `git log -1 ${hash} --pretty=format:"%ad" --date=relative`,
            { cwd: repoPath, windowsHide: true }
          )
            .then(({ stdout, stderr }) => {
              if (stdout) {
                return stdout.trim();
              } else {
                console.log(stderr);
                return "";
              }
            })
            .catch((err) => {
              console.log(err);
              return "";
            });

          const commitFilesCount = await execPromisified(
            `git diff-tree --no-commit-id --name-only -r ${hash}`,
            { cwd: repoPath, windowsHide: true }
          )
            .then(({ stdout, stderr }) => {
              if (stdout) {
                return stdout.trim().split("\n").length;
              } else {
                return 0;
              }
            })
            .catch((err) => {
              return 0;
            });

          commits += "||" + commitRelativeTime;
          commits += "||" + commitFilesCount;
          return commitModel(commits);
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
  });

  commitSearchResults = commitSearchResults.map((item) => {
    return item.then((res) => res);
  });

  return {
    commits: [...commitSearchResults],
  };
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

module.exports.gitCommitLogSearchHandler = gitCommitLogSearchHandler;
