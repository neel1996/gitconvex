const fs = require("fs");
const path = require("path");

function getEnvData() {
  const envFileData = fs.readFileSync(
    path.join(__dirname, "..", "env_config.json")
  );

  const envContent = envFileData.toString();
  let envData = JSON.parse(envContent)[0];

  return {
    DATABASE_FILE: envData.databaseFile,
    GITCONVEX_PORT: envData.port,
  };
}

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
