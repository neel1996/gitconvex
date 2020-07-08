const fs = require("fs");
const dotenv = require("dotenv").config();
const { DATABASE_FILE } = require("./envConfigReader").getEnvData();

const getRepoPath = (repoId) => {
  const dataEntry = fs
    .readFileSync(DATABASE_FILE || "./database/repo-datastore.json")
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
