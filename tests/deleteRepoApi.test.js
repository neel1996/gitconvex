const { deleteRepoApi } = require("../API/deleteRepoApi");
const { fetchRepoHandler } = require("../API/fetchRepoApi");

let repoIdList = [];

describe("Test module for - deleteRepoApi", () => {
  it("Tests repo deletion scenario", async () => {
    const fetchResults = await fetchRepoHandler().then((res) => res);
    repoIdList = await fetchResults.repoName.map((name, index) => {
      if (name === "JEST_REPO") {
        return fetchResults.repoId[index];
      } else {
        return null;
      }
    });

    let selectedId = "";

    repoIdList.forEach((item) => {
      if (item) {
        selectedId = item;
      }
    });
    const { status } = await deleteRepoApi(selectedId);

    expect(status).toBe("DELETE_SUCCESS");
  });
});
