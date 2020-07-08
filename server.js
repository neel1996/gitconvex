#!/usr/bin/env node

const globalAPI = require("./global/globalAPIHandler");
const express = require("express");
const path = require("path");
const fs = require("fs");
const dotenv = require("dotenv").config();
const {
  DATABASE_FILE,
  GITCONVEX_PORT,
} = require("./global/envConfigReader").getEnvData();

const app = globalAPI;
const port = GITCONVEX_PORT;

app.use(express.static(path.join(__dirname, "build")));

app.get("/*", (req, res) => {
  res.sendFile(path.join(__dirname, "build", "index.html"));
});

globalAPI.listen(port || 9001, async (err) => {
  if (err) {
    console.log(err);
  }

  console.log("GitConvex API connected!");

  console.log("\n#Checking data file availability...");

  await fs.promises
    .access(DATABASE_FILE)
    .then(() => {
      console.log(
        `INFO: Data file ${DATABASE_FILE} is present and it will be used as the active data file!\n\n## You can change this under the settings menu
        `
      );
    })
    .catch(async (err) => {
      const dataFileCreator = async () => {
        return await fs.promises
          .writeFile(DATABASE_FILE, "[]")
          .then((res) => {
            console.log(
              "\nINFO: New data file created and it will be used as the active file\n\n## You can change this under the settings menu"
            );
          })
          .catch((err) => {
            console.log(
              "INFO: New data file creation failed!\nINFO: Falling back to directory creation module"
            );
          });
      };

      console.log(
        `INFO: Data file is missing\nCreating new file under ${DATABASE_FILE}`
      );

      await dataFileCreator();

      if (fs.existsSync()) {
      } else {
        console.log("INFO: Database directory is missing");
        await fs.promises
          .mkdir("./database")
          .then(async () => {
            console.log(
              "INFO: Created database directory\nINFO: Setting up new data file in database directory"
            );
            await dataFileCreator();
          })
          .catch((err) => {
            console.log("ERROR: database directory creation failed!");
          });
      }
    });

  console.log(
    `\n## Gitconvex is running on port ${port}
     
    Open http://localhost:${port}/ to access gitconvex
    `
  );
});
