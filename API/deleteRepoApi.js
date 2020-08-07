const fs = require("fs");
const path = require("path");
const { getEnvData } = require("../utils/getEnvData");

/**
 * @param  {String} repoId - ID of the repo stored in the data file
 */


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
