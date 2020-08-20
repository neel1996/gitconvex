const { exec } = require("child_process");
const util = require("util");
const execPromise = util.promisify(exec);
const os = require("os");

/**
 * @returns {Object} - Version of the required software and the platform on which the app is running
 */

async function healthCheckHandler() {
  let healthCheckResults = {
    osCheck: "",
    gitCheck: "",
    nodeCheck: "",
  };

  healthCheckResults.osCheck = await checkStatus("OS").then((res) =>
    JSON.stringify(res)
  );
  healthCheckResults.gitCheck = await checkStatus("GIT").then((res) =>
    JSON.stringify(res)
  );
  healthCheckResults.nodeCheck = await checkStatus("NODE").then((res) =>
    JSON.stringify(res)
  );

  return { ...healthCheckResults };
}

/**
 * @param  {String} param
 * @returns {Object} - Executes the shell command and returns the results 
 */

async function checkStatus(param) {
  var commandString = "";

  switch (param) {
    case "GIT":
      commandString = `git --version`;
      break;
    case "NODE":
      commandString = `node --version`;
      break;
    default:
      commandString = ` `;
  }

  return await execPromise(commandString, { windowsHide: true }).then((res) => {
    if (param === "OS") {
      return {
        code: "SUCCESS",
        status: `${param}_CHECK_PASSED`,
        message: os.platform().toString(),
      };
    } else {
      if (res.stderr) {
        return {
          code: "ERR",
          status: `${param}_CHECK_FAILURE`,
          message: res.stderr,
        };
      } else {
        return {
          code: "SUCCESS",
          status: `${param}_CHECK_PASSED`,
          message: res.stdout.trim(),
        };
      }
    }
  });
}

module.exports.healthCheckHandler = healthCheckHandler;
