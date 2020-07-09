const { exec } = require("child_process");
const path = require("path");
const util = require("util");
const execPromisified = util.promisify(exec);
const fs = require("fs");

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

async function addRepoHandler(repoName, repoPath, initCheck) {
  const timeStamp = new Date().toUTCString();
  const id = new Date().getTime();

  var repoObject = {
    id,
    timeStamp,
    repoName,
    repoPath,
  };

  function errorResponse() {
    return {
      addRepo: {
        message: "REPO_WRITE_FAILED",
      },
    };
  }

  function successResponse() {
    return {
      addRepo: {
        message: "REPO_DATA_UPDATED",
        repoId: id,
      },
    };
  }

  const dataStoreFile = getEnvData().DATABASE_FILE;

  let fileData = fs.readFileSync(dataStoreFile);
  const repoData = fileData.toString();

  return await fs.promises
    .access(repoPath)
    .then(async () => {
      if (initCheck) {
        await execPromisified(`git init`, { cwd: repoPath, windowsHide: true });
      }

      if (repoData) {
        const existingData = JSON.parse(repoData);

        existingData.push(repoObject);

        fs.writeFileSync(dataStoreFile, JSON.stringify(existingData));
        return successResponse();
      } else {
        fs.writeFileSync(dataStoreFile, JSON.stringify([repoObject]));
        return successResponse();
      }
    })
    .catch((err) => {
      console.log(err);
      return errorResponse();
    });
}

module.exports.addRepoHandler = addRepoHandler;
