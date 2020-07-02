const { exec } = require("child_process");
const util = require("util");
const execPromise = util.promisify(exec);
const os = require("os");

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
      commandString = `echo 'NO_COMMAND'`;
  }

  return await execPromise(commandString).then((res) => {
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
