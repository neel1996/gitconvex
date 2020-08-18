const { getEnvData } = require("./getEnvData");
const { fetchRepoHandler } = require("../API/fetchRepoApi");
const { gitCommitLogHandler } = require("../git/gitCommitLogsAPI");
const sqlite = require("sqlite3").verbose();
const path = require("path");

async function gitCommitLogToDb() {
  const db = new sqlite.Database(
    getEnvData().COMMITLOG_DB ||
      path.join(__dirname, "..", "/database/commitLogs.sqlite"),
    (err) => {
      if (err) {
        console.log(err);
      }
    }
  );

  console.log("INFO: Initiaitng SQLite DB module for commit logs");

  db.serialize(async () => {
    const repoList = await fetchRepoHandler();

    repoList &&
      repoList.repoId.forEach(async (repoId) => {
        db.run(
          `CREATE TABLE IF NOT EXISTS commitLog_${repoId} (hash TEXT NOT NULL PRIMARY KEY, author TEXT, commit_date TEXT, commit_message TEXT, commit_relative_time TEXT)`
        );

        await gitCommitLogHandler(repoId, 0)
          .then(async (res) => {
            inserToDbHandler(Promise.all(res.commits), db, repoId);
            let totalCommits = res.totalCommits / 10 + 1;
            let skipLimit = 10;

            if (totalCommits > 500) {
              console.log(
                "WARN: The commit log volume is huge, so the log crawling process will take a while to complete!"
              );
            }

            db.get(
              `SELECT count(*) as rowCount from commitLog_${repoId}`,
              [],
              async (err, row) => {
                if (err) {
                  console.log("ERROR: ", err);
                } else {
                  console.log(
                    "INFO: Total commits in repo : ",
                    repoId,
                    res.totalCommits
                  );
                  console.log(
                    "INFO: Total commits in Database commitLog_",
                    repoId,
                    row.rowCount
                  );

                  if (res.totalCommits === row.rowCount) {
                    console.log(
                      `INFO: commitLog_${repoId} Database is up to date`
                    );
                  } else {
                    console.log(
                      `INFO: Inserting new logs to commitLog_${repoId} Database`
                    );
                    for (let i = 0; i <= totalCommits; i++) {
                      skipLimit = i * 10;

                      await gitCommitLogHandler(repoId, skipLimit)
                        .then((res) => {
                          inserToDbHandler(
                            Promise.all(res.commits),
                            db,
                            repoId
                          );
                        })
                        .catch((err) => {
                          console.log(err);
                        });
                    }
                  }
                }
              }
            );
          })
          .catch((err) => {
            console.log(err);
          });
      });
  });
}

/**
 * @param  {Array} commitArray - Array containing the commit logs
 * @param  {Sqlite.Database} db - Sqlite DB instance
 * @param  {String} repoId - Repo ID stored in the data file
 */

async function inserToDbHandler(commitArray, db, repoId) {
  // let commitArray = commitLogs && Promise.all(commitLogs);

  let commitList = [];

  await commitArray
    .then(async (res) => {
      res.forEach((item) => {
        commitList.push(item);
      });
    })
    .catch((err) => {
      console.log(err);
    });

  commitList = commitList.filter((item) => (item ? true : false));

  commitList &&
    commitList.forEach(async (commitData) => {
      let {
        hash,
        author,
        commitTime,
        commitMessage,
        commitRelativeTime,
      } = commitData;
      let insertTracker = 0;
      let errorTracker = [];

      commitMessage = commitMessage.split('"').join('""');

      if (hash) {
        db.get(
          `SELECT hash from commitLog_${repoId} WHERE hash="${hash}"`,
          [],
          (err, rows) => {
            if (err) {
              console.log("ERROR: ", err);
            }
            if (rows === undefined) {
              insertTracker++;

              !rows &&
                db.run(
                  `INSERT INTO commitLog_${repoId}(hash,author,commit_date,commit_message,commit_relative_time) VALUES("${hash}", "${author}", "${commitTime}", "${commitMessage}", "${commitRelativeTime}")`,
                  (err) => {
                    if (err) {
                      errorTracker.push(err);
                    }
                  }
                );
            }
          }
        );
      } else {
        console.log("ERROR: Commit hash not received for validation");
      }
    });
}

gitCommitLogToDb();

module.exports.gitCommitLogToDb = gitCommitLogToDb;
