const { healthCheckHandler } = require("../API/healthcheckApi");

test("Test module for - healthCheckApi", async () => {
  let { osCheck, gitCheck, nodeCheck } = await healthCheckHandler();

  osCheck = JSON.parse(osCheck);
  nodeCheck = JSON.parse(nodeCheck);
  gitCheck = JSON.parse(gitCheck);

  expect(osCheck).toBeTruthy();
  expect(gitCheck).toBeTruthy();
  expect(nodeCheck).toBeTruthy();

  expect(osCheck.status.includes("PASSED")).toBeTruthy();
  expect(gitCheck.status.includes("PASSED")).toBeTruthy();
  expect(nodeCheck.status.includes("PASSED")).toBeTruthy();
});
