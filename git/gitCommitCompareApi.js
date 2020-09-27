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
          message: "Error occurred while fetching commit difference",
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

  console.log(diffFileItems);

  return diffFileItems;
}

module.exports.gitCommitCompare = gitCommitCompare;
