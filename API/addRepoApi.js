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
    return {
      message: "REPO_DATA_UPDATED",
      repoId: id,
    };
  }

  if (cloneCheck) {
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
