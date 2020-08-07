const path = require("path");
const fs = require("fs");
/**
 * @returns {Object} - reads the env_config json file ans returns the results
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
