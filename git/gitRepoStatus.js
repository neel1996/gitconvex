const { exec } = require("child_process");

const fs = require("fs");

const util = require("util");
const execPromised = util.promisify(exec);

const getGitStatus = async (repoPath) => {
  console.log("Repo Path : " + repoPath);

  let gitRemoteData = "";
  let gitBranchList = [];
  let gitCurrentBranch = "No Active Branch";
  let gitRemoteHost = "";
  let gitRepoName = "";
  let gitTotalCommits = "";
  let gitLatestCommit = "";
  let gitTrackedFiles = "";
  let gitTotalTrackedFiles = 0;

  const gitRemoteReference = [
    "github",
    "gitlab",
    "bitbucket",
    "azure",
    "codecommit",
  ];

  const currentDir = `cd ${repoPath};`;

  let isGitLogAvailable = fs.promises
    .access(`${repoPath}/.git/logs`)
    .then(() => {
      isGitLogAvailable = true;
      return isGitLogAvailable;
    })
    .catch((err) => {
      // console.log(err);
      isGitLogAvailable = false;
      return isGitLogAvailable;
    });

  // Module to get git remote repo URL

  let gitRemotePromise =
    isGitLogAvailable &&
    (await execPromised(`${currentDir} git remote`).then(
      ({ stdout, stderr }) => {
        if (stdout && !stderr) {
          const localRemote = stdout.trim().split("\n");

          const multiPromise = Promise.all(
            localRemote &&
              localRemote.map(async (remote) => {
                console.log("LOOP ::", remote);
                return await execPromised(
                  `${currentDir} git remote get-url ${remote}`
                ).then(({ stdout, stderr }) => {
                  if (stdout && !stderr) {
                    console.log("REMOTE :: ", stdout);
                    return stdout.trim();
                  } else {
                    console.log(stderr);
                  }
                });
              })
          );
          return multiPromise;
        } else {
          console.log(stderr);
          return null;
        }
      }
    ));

  if (gitRemotePromise) {
    gitRemoteData = gitRemotePromise.join("||");
  } else {
    gitRemoteData = "NO_REMOTE";
  }

  // Module to get Git actual repo name
  if (gitRemoteData && gitRemoteData !== "NO_REMOTE") {
    let tempSplitLength = gitRemoteData.split("/").length;
    gitRepoName = gitRemoteData
      .split("/")
      [tempSplitLength - 1].split(".git")[0];

    gitRemoteReference.forEach((entry) => {
      if (gitRemoteData.includes(entry)) {
        gitRemoteHost = entry;
      }
    });
  } else if (gitRemoteData === "NO_REMOTE") {
    gitRepoName = repoPath.split("/")[currentDir.split("/").length - 1];
    gitRemoteHost = "No Remote Host Set";
  }

  // Module to get all available branches
  gitBranchList =
    isGitLogAvailable &&
    (await execPromised(`${currentDir} git branch`)
      .then((res) => {
        if (!res.stderr) {
          return res.stdout;
        } else {
          console.log(res.stderr);
        }
      })
      .catch((err) => {
        console.log(err);
      }));

  gitBranchList =
    gitBranchList.length > 0 &&
    gitBranchList
      .split("\n")
      .map((entry) => {
        if (entry.includes("*")) {
          gitCurrentBranch = entry.trim().replace("*", "");
          return null;
        }
        return entry.trim();
      })
      .filter((entry) => (entry !== "" ? entry : null));

  if (gitCurrentBranch.length > 0 && gitCurrentBranch !== "No Active Branch") {
    gitBranchList = [gitCurrentBranch, ...gitBranchList];
  } else {
    gitBranchList = ["NO_BRANCH"];
  }

  // Module to get total number of commits to current branch
  isGitLogAvailable &&
    (await execPromised(`${currentDir} git log --oneline`)
      .then((res) => {
        const { stdout, stderr } = res;
        if (stderr) {
          console.log(stderr);
        }
        if (res && !res.stderr) {
          const gitLocalTotal = res.stdout.trim().split("\n");
          if (gitLocalTotal && gitLocalTotal.length > 0) {
            gitTotalCommits = gitLocalTotal.length;
          } else if (gitLocalTotal.length === 1) {
            gitTotalCommits = 1;
          }
        } else {
          gitTotalCommits = 0;
          console.log(stderr);
        }
        return gitTotalCommits;
      })
      .catch((err) => {
        gitTotalCommits = 0;
        console.log(err);
      }));

  //Module to get latest git commit

  isGitLogAvailable &&
    (await execPromised(`${currentDir} git log -1 --oneline`).then((res) => {
      if (res && !res.stderr) {
        gitLatestCommit = res.stdout.trim();
      }
    }));

  //Module to get all git tracked files
  var gitTrackedFileDetails = [];

  gitTrackedFiles =
    isGitLogAvailable &&
    (await execPromised(`${currentDir} git ls-tree --name-status HEAD`).then(
      ({ stdout, stderr }) => {
        if (stdout && !stderr) {
          const fileList = stdout.trim().split("\n");

          const localFiles = Promise.all(
            fileList.map(async (item) => {
              gitTrackedFileDetails.push(item);

              return await fs.promises
                .stat(`${item}`)
                .then((fileType) => {
                  if (fileType.isFile()) {
                    return `${item}: File`;
                  } else if (fileType.isDirectory()) {
                    return `${item}: directory`;
                  } else {
                    return `${item}: File`;
                  }
                })
                .catch((err) => {
                  console.log(err);
                  return `${item}: File`;
                });
            })
          );
          return localFiles;
        } else {
          console.log(stderr);
          return [];
        }
      }
    ));

  //Module to fetch commit for each file and folder

  var gitFileBasedCommit = [];

  gitFileBasedCommit =
    isGitLogAvailable &&
    (await Promise.all(
      gitTrackedFileDetails.map(async (gitFile) => {
        return await execPromised(
          `${currentDir} git log -1 --oneline ${gitFile}`
        ).then(({ stdout, stderr }) => {
          if (stdout && !stderr) {
            return stdout.trim();
          } else {
            console.log(stderr);
            return "";
          }
        });
      })
    ));

  //Module to get totally tracked git artifacts

  isGitLogAvailable &&
    (await execPromised(`${currentDir} git ls-files`).then((res) => {
      const { stdout, stderr } = res;
      if (stdout && !stderr) {
        if (stdout.split("\n")) {
          gitTotalTrackedFiles = Number(stdout.trim().split("\n").length);
        } else {
          return 0;
        }
      } else {
        console.log(stderr);
      }
    }));

  if (!isGitLogAvailable) {
    console.log("Untracked Git Repo!");
    gitTotalCommits = 0;
    gitLatestCommit = "No Commits";
    gitTrackedFiles = ["NO_TRACKED_FILES"];
    gitFileBasedCommit = "No Changes";
    gitTotalTrackedFiles = 0;
  }

  const gitRepoDetails = {
    gitRemoteData,
    gitRepoName,
    gitBranchList,
    gitCurrentBranch,
    gitRemoteHost,
    gitTotalCommits,
    gitLatestCommit,
    gitTrackedFiles,
    gitFileBasedCommit,
    gitTotalTrackedFiles,
  };

  console.log(gitRepoDetails);

  return gitRepoDetails;
};

module.exports.getGitStatus = getGitStatus;
