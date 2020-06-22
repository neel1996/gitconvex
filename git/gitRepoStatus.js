const { exec } = require("child_process");

const fs = require("fs");

const util = require("util");
const { stderr } = require("process");
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
    .then((res) => {
      if (res) {
        isGitLogAvailable = true;
        return isGitLogAvailable;
      }
    })
    .catch((err) => {
      // console.log(err);
      isGitLogAvailable = false;
      return isGitLogAvailable;
    });

  // Module to get git remote repo URL
  await execPromised(
    `${currentDir} if [ ! -z "\`git remote\`" ]; then git remote | xargs -L 1 git remote get-url; else echo "NO_REMOTE"; fi`
  )
    .then((res) => {
      const { stdout, stderr } = res;
      if (stderr !== "") {
        console.log(stderr);
        gitRemoteData = "NO_REMOTE";
      } else {
        gitRemoteData = stdout.trim();
        console.log("REMOTE : " + gitRemoteData);
        if (gitRemoteData.split("\n").length > 0) {
          const splitRemote = gitRemoteData.split("\n");
          gitRemoteData = splitRemote.join("||");
          console.log(gitRemoteData);
        }
      }
    })
    .catch((err) => {
      console.log("Error GIT : " + err);
      gitRemoteData = "NO_REMOTE";
    });

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
          return;
        }
        return entry.trim();
      })
      .filter((entry) => (entry !== "" ? entry : null));

  gitBranchList = [gitCurrentBranch, ...gitBranchList];

  // Module to get total number of commits to current branch
  isGitLogAvailable &&
    (await execPromised(`${currentDir} git log --oneline | wc -l`)
      .then((res) => {
        const { stdout, stderr } = res;
        if (stderr) {
          console.log(stderr);
        }
        if (res && !res.stderr) {
          gitTotalCommits = res.stdout.trim();
        }
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

  isGitLogAvailable &&
    (await execPromised(
      `${currentDir} for i in \`git ls-tree --name-status HEAD\`; do if [ -f $i ] || [ -d $i ] ; then file $i; fi; done`
    ).then((res) => {
      const { stdout, stderr } = res;
      if (res && !stderr) {
        gitTrackedFiles = stdout.trim().split("\n");
      } else {
        console.log(stderr);
      }
    }));

  //Module to fetch commit for each file and folder

  var gitFileBasedCommit = [];

  isGitLogAvailable &&
    (await execPromised(
      `${currentDir} for i in \`git ls-tree --name-status HEAD\`; do git log -1 --oneline $i; done 2> /dev/null`
    ).then((res) => {
      const { stdout, stderr } = res;

      if (res && !stderr) {
        gitFileBasedCommit = stdout
          .split("\n")
          .filter((elm) => (elm ? elm : null));
      } else {
        console.log(stderr);
      }
    }));

  //Module to get totally tracked git artifacts

  isGitLogAvailable &&
    (await execPromised(`${currentDir} git ls-files | wc -l`).then((res) => {
      const { stdout, stderr } = res;
      if (res && !stderr) {
        gitTotalTrackedFiles = Number(stdout.trim());
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
