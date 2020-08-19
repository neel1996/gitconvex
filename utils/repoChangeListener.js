const chokidar = require("chokidar");
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
        .watch(repo, { interval: 1500, usePolling: true })
        .on("change", (path, stats) => {
          console.log("INFO: change noticed in ", path);
          gitCommitLogToDb();
        });
    });
  }
}

module.exports.gitRepoListener = gitRepoListener;
