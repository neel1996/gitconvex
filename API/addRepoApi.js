const express = require("express");
const app = express();
const fs = require("fs");

async function addRepoHandler(repoName, repoPath) {
  const timeStamp = new Date().toUTCString();
  const id = new Date().getTime();

  var repoObject = {
    id,
    timeStamp,
    repoName,
    repoPath,
  };

  const dataStoreFile = "./database/repo-datastore.json";

  let fileData = fs.readFileSync(dataStoreFile);
  const repoData = fileData.toString();

  if (repoData) {
    const existingData = JSON.parse(repoData);

    existingData.push(repoObject);

    fs.writeFileSync(dataStoreFile, JSON.stringify(existingData));

    return {
      addRepo: {
        message: "REPO_DATA_UPDATED",
      },
    };
  } else {
    return {
      addRepo: {
        message: "REPO_WRITE_FAILED",
      },
    };
  }
}

module.exports.addRepoHandler = addRepoHandler;
