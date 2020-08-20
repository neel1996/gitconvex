const fs = require("fs");
const path = require("path");
const { getEnvData } = require("../utils/getEnvData");

/**
 * @returns {Object} - Path where the Data file is stored
 */

const fetchDatabaseFile = async () => {
  const dbPath = getEnvData().DATABASE_FILE || "NO_DATABASE_FILE";
  console.log("Database File", dbPath);
  return {
    settingsDatabasePath: dbPath.toString(),
  };
};

/**
 * @returns {Array} - The array containing the list of all repos stored in the data file
 */

const fetchRepoDetails = async () => {
  return await fs.promises
    .readFile(getEnvData().DATABASE_FILE)
    .then((res) => {
      const fileContent = res.toString();
      let parsedFileContent = [];
      if (fileContent.length > 0) {
        try {
          parsedFileContent = JSON.parse(fileContent);
        } catch (e) {
          console.log(e);
          parsedFileContent = [];
        }
      } else {
        fs.writeFileSync(getEnvData().DATABASE_FILE, "[]");

        return {
          settingsRepoDetails: [],
        };
      }
      return {
        settingsRepoDetails: parsedFileContent,
      };
    })
    .catch((err) => {
      if (err) {
        console.log(err);
        return {
          settingsRepoDetails: [],
        };
      }
    });
};

/**
 * @param  {String} newFileName - A new JSON file for storing repo information
 * @returns {String} - updates the data file in config file and returns the status as string
 */

const updateDbFile = async (newFileName) => {
  console.log("FILE NAME : ", newFileName);

  return await fs.promises
    .access(newFileName)
    .then(async (res) => {
      const envContent = fs
        .readFileSync(path.join(__dirname, "..", "env_config.json"))
        .toString();

      const parsedEnvContent = JSON.parse(envContent)[0];
      parsedEnvContent.databaseFile = newFileName.toString();

      console.log(JSON.stringify(parsedEnvContent));

      return await fs.promises
        .writeFile(
          path.join(__dirname, "..", "env_config.json"),
          JSON.stringify([parsedEnvContent])
        )
        .then(() => {
          return "DATAFILE_UPDATE_SUCCESS";
        })
        .catch((err) => {
          console.log(err);
          return "DATAFILE_UPDATE_FAILED";
        });
    })
    .catch((err) => {
      console.log(err);
      return "DATAFILE_UPDATE_FAILED";
    });
};

const getPortDetails = async () => {
  return { settingsPortDetails: Number(getEnvData().GITCONVEX_PORT) };
};

/**
 * @param  {number} newPort - new port number for running gitconvex
 * @returns {String} - status string for updating the port
 */

const updatePortDetails = async (newPort) => {
  if (!isNaN(newPort)) {
    const envContent = fs
      .readFileSync(path.join(__dirname, "..", "env_config.json"))
      .toString();

    const parsedEnvContent = JSON.parse(envContent)[0];
    parsedEnvContent.port = Number(newPort);

    console.log(JSON.stringify(parsedEnvContent));

    return await fs.promises
      .writeFile(
        path.join(__dirname, "..", "env_config.json"),
        JSON.stringify([parsedEnvContent])
      )
      .then(() => {
        return "PORT_UPDATE_SUCCESS";
      })
      .catch((err) => {
        console.log(err);
        return "PORT_UPDATE_FAILED";
      });
  } else {
    return "PORT_UPDATE_FAILED";
  }
};

module.exports.updateDbFile = updateDbFile;
module.exports.fetchDatabaseFile = fetchDatabaseFile;
module.exports.fetchRepoDetails = fetchRepoDetails;
module.exports.updatePortDetails = updatePortDetails;
module.exports.getPortDetails = getPortDetails;
