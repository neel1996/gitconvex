const fs = require("fs");
const dotenv = require("dotenv").config();

const getRepoPath = (repoId) => {
  const dataEntry = fs
    .readFileSync(process.env.DATABASE_FILE || "./database/repo-datastore.json")
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
