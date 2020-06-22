const fs = require("fs");

const getRepoPath = (repoId) => {
  const dataEntry = fs
    .readFileSync("./database/repo-datastore.json")
    .toString();

  const repoObject = JSON.parse(dataEntry);
  var repoPath = "";

  repoObject.forEach((entry) => {
    if (entry.id == repoId) {
      repoPath = entry.repoPath;
    }
  });

  return repoPath;
};

module.exports.getRepoPath = getRepoPath;
