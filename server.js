#!/usr/bin/env node

const globalAPI = require("./global/globalAPIHandler");
const express = require("express");
const path = require("path");
const fs = require("fs");

const { gitCommitLogToDb } = require("./utils/sqliteDbAccess");
const { gitRepoListener } = require("./utils/repoChangeListener");
const app = globalAPI;
const log = console.log;
var envConfigFilename = "env_config.json";
var envConfigFilePath = path.join(__dirname, envConfigFilename);
var isReactBundlePresent = false;

try {
  fs.accessSync(path.join(__dirname, ".", "build"));
  isReactBundlePresent = true;
} catch (err) {
  log("ERROR: ", err);
  isReactBundlePresent = false;
}

if (isReactBundlePresent) {
  app.use(express.static(path.join(__dirname, "build")));
}

function getEnvData() {
  try {
    const envFileData = fs.readFileSync(
      path.join(__dirname, "env_config.json")
    );
    const envContent = envFileData.toString();
    let envData = JSON.parse(envContent)[0];

    let insertFlag = false;

    if (!envData.databaseFile) {
      envData["databaseFile"] = path.join(
        __dirname,
        "database/repo-datastore.json"
      );
      insertFlag = true;
      log("INFO: Inserting new date --> DATABASE_FILE");
    }
    if (!envData.port) {
      envData["port"] = 9001;
      insertFlag = true;
      log("INFO: Inserting new date --> PORT");
    }
    if (!envData.commitLogDatabase) {
      envData["commitLogDatabase"] = path.join(
        __dirname,
        "database/commitLogs.sqlite"
      );
      insertFlag = true;
      log("INFO: Inserting new date --> COMMITLOG_DATABASE");
    }

    if (insertFlag) {
      writeConfigFile(insertFlag, envData);
    }

    return {
      DATABASE_FILE: envData.databaseFile,
      GITCONVEX_PORT: envData.port,
    };
  } catch (err) {
    console.log("ERROR: Error occurred while reading env_config file", err);
    writeConfigFile();
    return {
      DATABASE_FILE: path.join(__dirname, "database/repo-datastore.json"),
      GITCONVEX_PORT: 9001,
    };
  }
}

function writeConfigFile(insertFlag = false, envData = {}) {
  let configData = [
    {
      databaseFile: path.join(__dirname, "database/repo-datastore.json"),
      commitLogDatabase: path.join(__dirname, "database/commitLogs.sqlite"),
      port: 9001,
    },
  ];

  if (insertFlag) {
    log("INFO: Inserting new data to config file");
    configData = [
      {
        ...envData,
      },
    ];
  } else {
    log(
      "INFO: Creating config file with default config -> " +
        JSON.stringify(configData)
    );
  }
  fs.writeFileSync(envConfigFilePath, JSON.stringify(configData), {
    flag: "w",
  });
}

log("INFO: Checking for config file");

let configStatus = "";
try {
  configStatus = fs.accessSync(envConfigFilePath);
} catch (e) {
  log("ERROR: No config file found. Falling back to config creation module");
  writeConfigFile();
}

log("INFO: Config file is present");
log("INFO: Reading from config file " + envConfigFilePath);

app.get("/*", (req, res) => {
  if (isReactBundlePresent) {
    res.sendFile(path.join(__dirname, "build", "index.html"));
  }
});

globalAPI.listen(getEnvData().GITCONVEX_PORT || 9001, async (err) => {
  if (err) {
    log(err);
  }

  log("GitConvex API connected!");

  log("\n#Checking data file availability...");

  var DATABASE_FILE = getEnvData().DATABASE_FILE;

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
    });

  try {
    fs.accessSync(path.join(__dirname, ".", "database"));
    gitCommitLogToDb();
    gitRepoListener();
  } catch (err) {
    console.log(err);
    fs.mkdirSync(path.join(__dirname, ".", "database"));
    gitCommitLogToDb();
    gitRepoListener();
  }

  log(
    `\n## Gitconvex is running on port ${getEnvData().GITCONVEX_PORT}
     
    Open http://localhost:${getEnvData().GITCONVEX_PORT}/ to access gitconvex
    `
  );
});
