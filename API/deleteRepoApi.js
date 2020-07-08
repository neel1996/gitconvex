const fs = require("fs");
const dotenv = require("dotenv").config();
const { DATABASE_FILE } = require("../global/envConfigReader").getEnvData();

async function deleteRepoApi(repoId, name, pathName, time) {
  const dataStoreFile = DATABASE_FILE || "./database/repo-datastore.json";

  return await fs.promises
    .readFile(dataStoreFile)
    .then(async (data) => {
      const fileContent = JSON.parse(data.toString());

      if (fileContent && fileContent.length > 0) {
        let updatedData = fileContent.filter(({ id, repoName, repoPath }) => {
          if (id.toString() === repoId.toString()) {
            console.log("REPO DELETED");
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
