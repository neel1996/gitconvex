const chokidar = require("chokidar");
const path = require("path");
const { fetchRepoHandler } = require("../API/fetchRepoApi");
const { gitCommitLogToDb } = require("./sqliteDbAccess");

async function gitRepoListener() {
  console.log(
    "INFO: Repo path listener initiated. All changes in the configured repos will be tracked till gitconvex is stopped"
  );
  const { repoPath } = await fetchRepoHandler();

  if (repoPath) {
    repoPath.forEach((repo) => {
      chokidar
        .watch(path.join(repo, ".", ".git"))
        .on("change", (change, stats) => {
          console.log("INFO: change noticed in ", change);
          gitCommitLogToDb();
        });
    });
  }
}

module.exports.gitRepoListener = gitRepoListener;
