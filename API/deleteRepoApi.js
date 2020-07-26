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

async function deleteRepoApi(repoId) {
  const dataStoreFile = getEnvData().DATABASE_FILE;

  return await fs.promises
    .readFile(dataStoreFile)
    .then(async (data) => {
      const fileContent = JSON.parse(data.toString());

      if (fileContent && fileContent.length > 0) {
        let updatedData = fileContent.filter(({ id }) => {
          if (id.toString() === repoId.toString()) {
            return false;
          }

          return true;
        });

        return await fs.promises
          .writeFile(dataStoreFile, JSON.stringify(updatedData))
          .then(() => {
            return {
              status: "DELETE_SUCCESS",
              repoId,
            };
          });
      } else {
        return {
          status: "DELETE_FAILED",
        };
      }
    })
    .catch((err) => {
      console.log(err);
      return {
        status: "DELETE_FAILED",
      };
    });
}

module.exports.deleteRepoApi = deleteRepoApi;
