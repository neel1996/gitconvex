const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

async function gitCommitCompare(repoPath, baseCommit, compareCommit) {
  const diffFileItems = await execPromisified(
    `git diff ${baseCommit}..${compareCommit} --raw`,
    { windowsHide: true, maxBuffer: 1024 * 1024, cwd: repoPath }
  )
    .then((res) => {
      const { stdout, stderr } = res;

      if (stderr) {
        console.log(stderr);
        return {
          message: stderr,
        };
      }

      if (stdout) {
        const fileChanges = stdout.trim().split("\n");
        let diffArray = [];

        if (fileChanges) {
          diffArray = fileChanges.map((change) => {
            const splitItem = change.split(/\s/gi);

            let status = splitItem[4];
            let fileName = splitItem[splitItem.length - 1];

            return {
              status,
              fileName,
            };
          });

          return { difference: diffArray };
        } else {
          return {
            message: "No difference found",
          };
        }
      }
    })
    .catch((err) => {
      console.log("ERROR", err);
      return {
        message: "Error occurred while comparing branches",
      };
    });

  if (diffFileItems && diffFileItems.length > 0) {
    console.log(diffFileItems);
    return diffFileItems;
  } else {
    return {
      message:
        "warning: There are no file changes in the base commit. Make sure that it is not a PR merge commit!",
    };
  }
}

module.exports.gitCommitCompare = gitCommitCompare;
