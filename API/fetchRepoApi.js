const fs = require("fs");
const dotnev = require("dotenv").config();
const { DATABASE_FILE } = require("../global/envConfigReader").getEnvData();

async function fetchRepoHandler() {
  var repoDSContent = fs.readFileSync(
    DATABASE_FILE || "./database/repo-datastore.json"
  );

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
