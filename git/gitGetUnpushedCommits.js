const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);

const fetchRepopath = require("../global/fetchGitRepoPath");

const gitGetUnpushedCommits = async (repoId, remoteURL) => {
  return await execPromisified(
    `git log --branches --not --remotes --pretty=format:"%h||%an||%ad||%s"`,
    { cwd: fetchRepopath.getRepoPath(repoId) }
  ).then((res) => {
    const { stdout, stderr } = res;

    if (stderr) {
      console.log(stderr);
      return [];
    } else {
      if (stdout) {
        const splitCommits = stdout.trim().split("\n");
        console.log(splitCommits);

        return splitCommits;
      }
    }
  });
};

module.exports.gitGetUnpushedCommits = gitGetUnpushedCommits;
