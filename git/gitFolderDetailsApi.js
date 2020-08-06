const { exec } = require("child_process");
const fs = require("fs");
const util = require("util");
const execPromised = util.promisify(exec);
const fetchRepopath = require("../global/fetchGitRepoPath");
const path = require("path");

const gitFetchFolderContentApi = async (repoId, directoryName) => {
  const repoPath = fetchRepopath.getRepoPath(repoId);
  const targetPath = path.join(repoPath, directoryName);

  return await fs.promises.readdir(targetPath).then(async (folderContent) => {
    /* @gitcommit section */
    const gitCommits = folderContent.map(async (item) => {
      let commitCommand = "";
      if (directoryName) {
        commitCommand = `git log -1 --oneline "${directoryName + "/" + item}"`;
      } else {
        commitCommand = `git log -1 --oneline "${item}"`;
      }
      return await execPromised(commitCommand, {
        cwd: repoPath,
        windowsHide: true,
      })
        .then(({ stdout, stderr }) => {
          if (stdout) {
            return stdout.trim();
          } else {
            console.log(stderr);
          }
        })
        .catch((err) => {
          console.log(err);
        });
    });

    /* @folderContentTagging module */
    const folderObjects = folderContent.map(async (item, index) => {
      return await fs.promises
        .stat(path.join(targetPath, item))
        .then(async (content) => {
          if (await gitCommits[index]) {
            if (content.isFile()) {
              return `${item}: File`;
            } else if (content.isDirectory()) {
              return `${path.join(directoryName, item)}: directory`;
            } else {
              return `${item}: File`;
            }
          }
        })
        .catch((err) => {
          console.log(err);
          return [];
        });
    });

    return {
      gitFolderContent: {
        gitTrackedFiles: folderObjects,
        gitFileBasedCommit: gitCommits,
      },
    };
  });
};

module.exports.gitFetchFolderContentApi = gitFetchFolderContentApi;
