#!/usr/bin/env node

const globalAPI = require("./global/globalAPIHandler");
const express = require("express");
const path = require("path");
const fs = require("fs");
const dotenv = require("dotenv").config();
const { updateDbFile } = require("./API/settingsApi");
const app = globalAPI;
const log = console.log;
var envConfigFilename = "env_config.json";
var envConfigFilePath = path.join(__dirname, envConfigFilename);

app.use(express.static(path.join(__dirname, "build")));

log("INFO: Checking for config file");

let configStatus = "";
try {
  configStatus = fs.accessSync(envConfigFilePath);
} catch (e) {
  log("ERROR: No config file found. Falling back to config creation module");
  const configData = [
    {
      databaseFile: path.join(__dirname, "database/repo-datastore.json"),
      port: 9001,
    },
  ];
  log(
    "INFO: Creating config file with default config -> " +
      JSON.stringify(configData)
  );
  fs.writeFileSync(envConfigFilePath, JSON.stringify(configData), {
    flag: "w",
  });
}

let {
  DATABASE_FILE,
  GITCONVEX_PORT,
} = require("./global/envConfigReader").getEnvData();
const port = GITCONVEX_PORT;

log("INFO: Config file is present");
log("INFO: Reading from config file " + envConfigFilePath);

app.get("/*", (req, res) => {
  res.sendFile(path.join(__dirname, "build", "index.html"));
});

globalAPI.listen(port || 9001, async (err) => {
  if (err) {
    log(err);
  }

  log("GitConvex API connected!");

  log("\n#Checking data file availability...");

  await fs.promises
    .access(DATABASE_FILE)
    .then(() => {
      log(
        `INFO: Data file ${DATABASE_FILE} is present and it will be used as the active data file!\n\n## You can change this under the settings menu
        `
      );
    })
    .catch(async (err) => {
      const dataFileCreator = async () => {
        return await fs.promises
          .writeFile(DATABASE_FILE, "[]")
          .then((res) => {
            log(
              "\nINFO: New data file created and it will be used as the active file\n\n## You can change this under the settings menu"
            );
          })
          .catch((err) => {
            log(
              "INFO: New data file creation failed!\nINFO: Falling back to directory creation module"
            );
          });
      };

      log(
        `INFO: Data file is missing\nCreating new file under ${DATABASE_FILE}`
      );

      await dataFileCreator();

      if (!fs.existsSync(DATABASE_FILE)) {
        log("INFO: Database directory is missing");
        await fs.promises
          .mkdir(path.join(__dirname, ".", "/database"))
          .then(async () => {
            log(
              "INFO: Created database directory\nINFO: Setting up new data file in database directory"
            );
            await dataFileCreator();
          })
          .catch((err) => {
            log("ERROR: database directory creation failed!");
          });
      }
    });

  log(
    `\n## Gitconvex is running on port ${port}
     
    Open http://localhost:${port}/ to access gitconvex
    `
  );
});
