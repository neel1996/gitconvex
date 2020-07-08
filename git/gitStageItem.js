const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitStageItem = async (repoId, item) => {
  return await execPromisified(`git add ${item}`, {
    cwd: fetchRepopath.getRepoPath(repoId),
    windowsHide: true,
  }).then(({ stdout, stderr }) => {
    if (stderr) {
      console.log(stderr);
      return "ADD_ITEM_FAILED";
    } else {
      console.log(stdout);
      return "ADD_ITEM_SUCCES";
    }
  });
};

module.exports.gitStageItem = gitStageItem;
