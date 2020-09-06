const { exec } = require("child_process");
const path = require("path");
const util = require("util");
const execPromisified = util.promisify(exec);
const fs = require("fs");
const { getEnvData } = require("../utils/getEnvData");
const { gitCommitLogToDb } = require("../utils/sqliteDbAccess");

/**
 * @param  {String} repoName - Name of the repository
 * @param  {String} repoPath - Path were the repo is residing or where it should be cloned / initialized
 * @param  {boolean} initCheck - Switch to check if user has selected repo init option
 * @param  {boolean} cloneCheck - Switch to check if user has selected repo cloning option
 * @param  {String} cloneUrl - The If cloning switch is true, then this holds the URL of the remote repo
 * @returns {Object} - created a new entry in the data file and retusn the status
 */

async function addRepoHandler(
  repoName,
  repoPath,
  initCheck,
  cloneCheck,
  cloneUrl
) {
  const timeStamp = new Date().toUTCString();
  const id = new Date().getTime();

  function errorResponse() {
    return {
      message: "REPO_WRITE_FAILED",
    };
  }

  function successResponse() {
    gitCommitLogToDb();
    return {
      message: "REPO_DATA_UPDATED",
      repoId: id,
    };
  }

  if (cloneCheck) {
    try {
      if (cloneUrl.match(/[^a-zA-Z0-9-_.~@#$%:/]/gi)) {
        throw new Error("Invalid clone URL string!");
      }

      if (repoName.match(/[^a-zA-Z0-9-_.\s]/gi)) {
        throw new Error("Invalid repo name string!");
      }

      const cloneStatus = await execPromisified(
        `git clone "${cloneUrl}" "./${repoName}"`,
        {
          cwd: repoPath,
          windowsHide: true,
        }
      )
        .then(({ stdout, stderr }) => {
          console.log(stdout);
          console.log(stderr);
          if (stdout || stderr) {
            console.log(stdout);
            return true;
          } else {
            return false;
          }
        })
        .catch((err) => {
          console.log(err);
          return false;
        });

      console.log("CLONE STAT : ", cloneStatus);

      if (cloneStatus) {
        if (repoPath.includes("\\")) {
          repoPath = repoPath + "\\" + repoName;
        } else {
          repoPath = repoPath + "/" + repoName;
        }
      } else {
        return errorResponse();
      }
    } catch (err) {
      console.log(err);
      return errorResponse();
    }
  }

  const repoObject = {
    id,
    timeStamp,
    repoName,
    repoPath,
  };

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
