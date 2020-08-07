const fs = require("fs");
const path = require("path");
const { getEnvData } = require("../utils/getEnvData");

/**
 * @returns {Object} - details about all the repos stored in the data file
 */

async function fetchRepoHandler() {
  var repoDSContent = fs.readFileSync(getEnvData().DATABASE_FILE);

  repoDSContent = repoDSContent.toString();

  if (repoDSContent !== "") {
    let parsedData = JSON.parse(JSON.parse(JSON.stringify(repoDSContent)));

    let repoId = [];
    let repoName = [];
    let repoPath = [];

    parsedData.forEach((item) => {
      if (item) {
        repoId.push(item.id);
        repoName.push(item.repoName);
        repoPath.push(item.repoPath);
      }
    });

    const fetchRepo = {
      repoId,
      repoName,
      repoPath,
    };
    return fetchRepo;
  } else {
    return [];
  }
}

module.exports.fetchRepoHandler = fetchRepoHandler;
