const fetchRepopath = require("../global/fetchGitRepoPath");
const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

/**
 * @param  {String} repoId
 * @param  {String} commitHash
 * @returns {Array[type: String, fileName: String]} - fetches the list of files that are changed as part of a commit and returns the file with change type
 */

const gitCommitFileApi = async (repoId, commitHash) => {
  try {
    if (commitHash.match(/[^a-zA-Z0-9]/gi)) {
      throw new Error("Invliad commit hash string!");
    }

    const repoPath = fetchRepopath.getRepoPath(repoId);
    return await execPromisified(
      `git diff-tree --no-commit-id --name-status -r "${commitHash}"`,
      { cwd: repoPath, windowsHide: true }
    )
      .then(({ stdout, stderr }) => {
        if (stdout) {
          const commitedFiles = stdout.trim().split("\n");
          return commitedFiles.map((entry) => {
            if (entry) {
              const splitEntry = entry.split(/\s/gi);
              return {
                type: splitEntry[0],
                fileName: splitEntry.slice(1, splitEntry.length).join(" "),
              };
            } else {
              return {
                type: "",
                fileName: "",
              };
            }
          });
        } else {
          console.log(stderr);
          return [
            {
              type: "",
              fileName: "",
            },
          ];
        }
      })
      .catch((err) => {
        console.log(err);
        return [
          {
            type: "",
            fileName: "",
          },
        ];
      });
  } catch (err) {
    console.log(err);
    return [
      {
        type: "",
        fileName: "",
      },
    ];
  }
};

module.exports.gitCommitFileApi = gitCommitFileApi;
