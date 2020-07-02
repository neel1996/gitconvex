const { exec } = require("child_process");
const util = require("util");
const execPromisified = util.promisify(exec);
const fs = require("fs");

const dotenv = require("dotenv").config();

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

  const dataStoreFile =
    process.env.DATABASE_FILE || "./database/repo-datastore.json";

  let fileData = fs.readFileSync(dataStoreFile);
  const repoData = fileData.toString();

  return await fs.promises
    .access(repoPath)
    .then(async () => {
      if (initCheck) {
        await execPromisified(`git init`, { cwd: repoPath });
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
