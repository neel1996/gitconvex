const dotenv = require("dotenv").config();
const fs = require("fs");
const { parse, stringify } = require("envfile");

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

const updateDbFile = async (newFileName) => {
  console.log("FILE NAME : ", newFileName);
  // DATABASE_FILE = "database/repo-datastore.json"
  // GITCONVEX_PORT = 9001

  return await fs.promises
    .access(newFileName)
    .then(async (res) => {
      const envContent = fs.readFileSync(".env").toString();

      const parsedEnvContent = parse(envContent);
      parsedEnvContent.DATABASE_FILE = `"${newFileName.toString()}"`;

      console.log(stringify(parsedEnvContent));

      return await fs.promises
        .writeFile(".env", stringify(parsedEnvContent))
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
  const envContent = fs.readFileSync(".env").toString();

  const parsedEnvContent = parse(envContent);

  console.log(stringify(parsedEnvContent));

  return { settingsPortDetails: Number(parsedEnvContent.GITCONVEX_PORT) };
};

const updatePortDetails = async (newPort) => {
  if (!isNaN(newPort)) {
    const envContent = fs.readFileSync(".env").toString();

    const parsedEnvContent = parse(envContent);
    parsedEnvContent.GITCONVEX_PORT = Number(newPort);

    console.log(stringify(parsedEnvContent));

    return await fs.promises
      .writeFile(".env", stringify(parsedEnvContent))
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
module.exports.deleteRepo = deleteRepo;
module.exports.updatePortDetails = updatePortDetails;
module.exports.getPortDetails = getPortDetails;
