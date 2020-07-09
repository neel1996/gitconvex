const fs = require("fs");
const path = require("path");

function getEnvData() {
  try {
    const envFileData = fs.readFileSync(
      path.join(__dirname, "..", "env_config.json")
    );

    const envContent = envFileData.toString();
    let envData = JSON.parse(envContent)[0];

    return {
      DATABASE_FILE: envData.databaseFile,
      GITCONVEX_PORT: envData.port,
    };
  } catch (e) {
    return {
      DATABASE_FILE: "",
      GITCONVEX_PORT: "",
    };
  }
}

module.exports.getEnvData = getEnvData;
