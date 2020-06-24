const { exec } = require("child_process");
const dotenv = require("dotenv").config();
const fs = require("fs");

const util = require("util");
const execPromised = util.promisify(exec);

const fetchDatabaseFile = async () => {
  const dbPath = process.env.DATABASE_FILE || "NO_DATABASE_FILE";
  console.log(dbPath);
  return {
    settingsDatabasePath: dbPath,
  };
};

const fetchRepoDetails = async () => {
  return await fs.promises
    .readFile(process.env.DATABASE_FILE)
    .then((res) => {
      const fileData = JSON.parse(res.toString());
      return {
        settingsRepoDetails: fileData,
      };
    })
    .catch((err) => {
      if (err) {
        return {
          settingsRepoDetails: [],
        };
      }
    });
};

const deleteRepo = async (repoId) => {
  return await fs.promises.readFile(process.env.DATABASE_FILE).then((res) => {
    console.log(res);
  });
};

module.exports.fetchDatabaseFile = fetchDatabaseFile;
module.exports.fetchRepoDetails = fetchRepoDetails;
module.exports.deleteRepo = deleteRepo;
