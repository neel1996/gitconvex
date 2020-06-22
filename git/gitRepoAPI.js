const getGitStatus = require("./gitRepoStatus").getGitStatus;
const fetchRepoPath = require("../global/fetchGitRepoPath");

async function getGitRepoStatus(repoId) {
  const repoPath = fetchRepoPath.getRepoPath(repoId);

  const repoDetails = await getGitStatus(repoPath).then((result) => {
    if (result) {
      return result;
    }
  });

  return repoDetails;
}

module.exports.getGitRepoStatus = getGitRepoStatus;
