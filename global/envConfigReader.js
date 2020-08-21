const fs = require("fs");
const path = require("path");

/**
 * @returns {Object} - env config file content as JSON
 */

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

module.exports.getEnvData = getEnvData;
