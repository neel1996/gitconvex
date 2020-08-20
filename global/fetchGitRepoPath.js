const fs = require("fs");
const path = require("path");
const { getEnvData } = require("../utils/getEnvData");

const getRepoPath = (repoId) => {
  const { DATABASE_FILE } = getEnvData();

  const dataEntry = fs.readFileSync(DATABASE_FILE).toString();

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
