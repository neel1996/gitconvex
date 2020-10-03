const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

async function gitBranchCompare(repoPath, baseBranch, compareBranch) {
  let formattedCommitLogs = [];

  const commitLogs = await execPromisified(
    `git log --oneline ${compareBranch}..${baseBranch} --pretty=format:"%h||%an||%ad||%s" --date=short`,
    { windowsHide: true, maxBuffer: 1024 * 1024, cwd: repoPath }
  )
    .then((res) => {
      const { stdout, stderr } = res;

      if (stderr) {
        console.log("ERROR", stderr);
        return {
          error: "Error occurred while comparing branches",
        };
      }

      if (stdout) {
        console.log(stdout.trim().split("\n"));
        formattedCommitLogs = [...commitLogModel(stdout.trim().split("\n"))];

        let dateGroupModel = [
          {
            date: "",
            commits: [
              {
                hash: "",
                author: "",
                commitMessage: "",
              },
            ],
          },
        ];

        let groupedCommits = [];

        for (let i = 0; i < formattedCommitLogs.length; i++) {
          let commitArray = [];
          let dateGroup = {
            date: formattedCommitLogs[i].commitDate,
          };

          for (var j = i; j < formattedCommitLogs.length; j++) {
            let item = formattedCommitLogs[j];
            if (formattedCommitLogs[i].commitDate === item.commitDate) {
              commitArray.push({
                hash: item.hash,
                author: item.author,
                commitMessage: item.commitMessage,
              });
            } else {
              break;
            }
          }

          i = j - 1;
          dateGroup.commits = commitArray;
          groupedCommits.push(dateGroup);
        }

        console.log(groupedCommits);

        return groupedCommits;
      }
    })
    .catch((err) => {
      console.log("ERROR", err);
      return {
        error: "Error occurred while comparing branches",
      };
    });

  return commitLogs;
}

function commitLogModel(commit) {
  const formattedCommits = commit.map((item) => {
    const splitCommit = item.split("||");

    let commitObject = {
      hash: "",
      author: "",
      commitDate: "",
      commitMessage: "",
    };

    commitObject.hash = splitCommit[0];
    commitObject.author = splitCommit[1];
    commitObject.commitDate = splitCommit[2];
    commitObject.commitMessage = splitCommit[3];

    return commitObject;
  });

  return formattedCommits;
}

module.exports.gitBranchCompare = gitBranchCompare;
