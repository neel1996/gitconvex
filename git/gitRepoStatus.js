const { exec } = require("child_process");
const fs = require("fs");
const util = require("util");
const execPromised = util.promisify(exec);

const getGitStatus = async (repoPath) => {
  console.log("Repo Path : " + repoPath);

  const errorStatus = {
    noRemote: "NO_REMOTE",
    noRemoteHost: "No Remote Host Set",
    noBranch: ["NO_BRANCH"],
    noActiveBranch: "No active branch",
    noCommits: "No Commits in the Current Branch",
    noFileCommits: ["NO_COMMITS"],
    noTrackedFiles: ["NO_TRACKED_FILES"],
  };

  let gitRemoteData = "";
  let gitBranchList = ["NO_BRANCHES"];
  let gitCurrentBranch = errorStatus.noActiveBranch;
  let gitRemoteHost = "";
  let gitRepoName = "";
  let gitTotalCommits = "";
  let gitLatestCommit = "";
  let gitTrackedFiles = "";
  let gitTotalTrackedFiles = 0;
  let gitAllBranchList = ["NO_BRANCHES"];
  let isGitLogAvailable = false;

  const gitRemoteReference = [
    "github",
    "gitlab",
    "bitbucket",
    "azure",
    "codecommit",
  ];

  const currentDir = `cd ${repoPath};`;

  isGitLogAvailable = await fs.promises
    .access(`${repoPath}/.git/logs`)
    .then(() => {
      isGitLogAvailable = true;
      return isGitLogAvailable;
    })
    .catch((err) => {
      console.log("Not a git repo or a new git repo with no commits!");
      // console.log(err);
      isGitLogAvailable = false;
      return isGitLogAvailable;
    });

  console.log("Git log available : ", isGitLogAvailable);

  // Module to get git remote repo URL

  let gitRemotePromise =
    isGitLogAvailable &&
    (await execPromised(`git remote`, {
      cwd: repoPath,
      windowsHide: true,
    }).then(({ stdout, stderr }) => {
      if (stdout && !stderr) {
        const localRemote = stdout.trim().split("\n");

        const multiPromise = Promise.all(
          localRemote &&
            localRemote.map(async (remote) => {
              console.log("Remote ::", remote);
              return await execPromised(`git remote get-url ${remote}`, {
                cwd: repoPath,
                windowsHide: true,
              }).then(({ stdout, stderr }) => {
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
    }));

  if (gitRemotePromise) {
    gitRemoteData = gitRemotePromise.join("||");
  } else {
    gitRemoteData = errorStatus.noRemote;
  }

  // Module to get Git actual repo name
  if (gitRemoteData && gitRemoteData !== errorStatus.noRemote) {
    let tempSplitLength = gitRemoteData.split("/").length;
    gitRepoName = gitRemoteData
      .split("/")
      [tempSplitLength - 1].split(".git")[0];

    gitRemoteReference.forEach((entry) => {
      if (gitRemoteData.includes(entry)) {
        gitRemoteHost = entry;
      }
    });
  } else if (gitRemoteData === errorStatus.noRemote) {
    gitRepoName = repoPath.split("/")[currentDir.split("/").length - 1];
    gitRemoteHost = errorStatus.noRemoteHost;
  }

  //Module to get all branch list
  gitAllBranchList =
    isGitLogAvailable &&
    (await execPromised(`git branch --all`, {
      cwd: repoPath,
      windowsHide: true,
    })
      .then((res) => {
        const { stdout, stderr } = res;
        if (stdout && !stderr) {
          let localBranchList = stdout.trim().split("\n");
          localBranchList = localBranchList.map((branch) => {
            return branch;
          });
          return localBranchList;
        } else {
          console.log(stderr);
          return [];
        }
      })
      .catch((err) => {
        console.log(err);
        return [];
      }));

  // Module to get all available branches
  gitBranchList =
    isGitLogAvailable &&
    (await execPromised(`git branch`, { cwd: repoPath, windowsHide: true })
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
    isGitLogAvailable &&
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

  if (
    isGitLogAvailable &&
    gitCurrentBranch.length > 0 &&
    gitCurrentBranch !== "No Active Branch"
  ) {
    gitBranchList = [gitCurrentBranch, ...gitBranchList];
  } else {
    gitBranchList = errorStatus.noBranch;
  }

  if (!gitBranchList && gitBranchList.length === 0) {
    gitBranchList = errorStatus.noBranch;
  }

  console.log("GIT BRANCH LIST", gitBranchList);

  // Module to get total number of commits to current branch
  isGitLogAvailable &&
    (await execPromised(`git log --oneline`, {
      cwd: repoPath,
      windowsHide: true,
    })
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
    (await execPromised(`git log -1 --oneline --pretty=format:"%s"`, {
      cwd: repoPath,
      windowsHide: true,
    })
      .then((res) => {
        if (res && !res.stderr) {
          gitLatestCommit = res.stdout.trim();
        } else {
          console.log(stderr);
          gitLatestCommit = errorStatus.noCommits;
        }
      })
      .catch((err) => {
        console.log(err);
        gitLatestCommit = errorStatus.noCommits;
      }));

  //Module to get all git tracked files
  var gitTrackedFileDetails = [];

  gitTrackedFiles =
    isGitLogAvailable &&
    (await execPromised(`git ls-tree --name-status HEAD`, {
      cwd: repoPath,
      windowsHide: true,
    })
      .then(({ stdout, stderr }) => {
        if (stdout && !stderr) {
          const fileList = stdout.trim().split("\n");

          const localFiles = Promise.all(
            fileList.map(async (item) => {
              gitTrackedFileDetails.push(item);

              return await fs.promises
                .stat(`${repoPath}/${item}`)
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
                  console.log("Tracked file has been removed!", err);
                  return `${item}: DEL`;
                });
            })
          );
          return localFiles;
        } else {
          console.log(stderr);
          return [];
        }
      })
      .catch((err) => {
        console.log(err);
      }));

  //Module to fetch commit for each file and folder

  var gitFileBasedCommit = errorStatus.noFileCommits;

  gitFileBasedCommit =
    isGitLogAvailable &&
    (await Promise.all(
      gitTrackedFileDetails.map(async (gitFile) => {
        return await execPromised(`git log -1 --oneline "${gitFile}"`, {
          cwd: repoPath,
          windowsHide: true,
        })
          .then(({ stdout, stderr }) => {
            if (stdout && !stderr) {
              return stdout.trim();
            } else {
              console.log(stderr);
              return errorStatus.noFileCommits;
            }
          })
          .catch((err) => {
            console.log("Tracked file has been removed!", err);
            return errorStatus.noFileCommits;
          });
      })
    ));

  //Module to get totally tracked git artifacts

  isGitLogAvailable &&
    (await execPromised(`git ls-files`, {
      cwd: repoPath,
      windowsHide: true,
    }).then((res) => {
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
    gitLatestCommit = errorStatus.noCommits;
    gitTrackedFiles = errorStatus.noTrackedFiles;
    gitFileBasedCommit = errorStatus.noFileCommits;
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
    gitAllBranchList,
  };

  console.log(gitRepoDetails);

  return gitRepoDetails;
};

module.exports.getGitStatus = getGitStatus;
