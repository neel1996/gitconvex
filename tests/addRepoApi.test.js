const { addRepoHandler } = require("../API/addRepoApi");

describe("Test module for - addRepoApi", () => {
  it("Tests Add repo without init and clone", async () => {
    const addRepoResult = await addRepoHandler("JEST_REPO", ".", false, false, "");

    expect(addRepoResult).toBeTruthy();
    expect(addRepoResult.repoId).toBeTruthy();
    expect(addRepoResult.message).toBe("REPO_DATA_UPDATED");
  });

  it("Tests Add repo feature for negative scenario", async () => {
    const addRepoResult = await addRepoHandler("JEST_REPO", "/test", false, false, "");

    expect(addRepoResult).toBeTruthy();
    expect(addRepoResult.message).toBe("REPO_WRITE_FAILED");
  })
});
