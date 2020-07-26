const { fetchRepoHandler } = require("../API/fetchRepoApi");

test("Test module for - fetchRepoApi", async () => {
  const repoDetails = await fetchRepoHandler();
  expect(repoDetails).toBeTruthy();

  expect(repoDetails.repoId.length).toBeGreaterThan(0);
  expect(repoDetails.repoName.length).toBeGreaterThan(0);
  expect(repoDetails.repoPath.length).toBeGreaterThan(0);
});
