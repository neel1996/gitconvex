const sqlite = require("sqlite3").verbose();
const path = require("path");
const db = new sqlite.Database(
  path.join(__dirname, "..", "/database/commitLogs.sqlite")
);
const { gitCommitLogSearchHandler } = require("../git/gitCommitLogSearchApi");

async function gitCommitLogDbSerchApi(repoId, searchCategory, searchKey) {
  let searchQuery = "";

  switch (searchCategory) {
    case "hash":
      searchQuery = `SELECT * FROM commitLog_${repoId} WHERE hash LIKE "%${searchKey}%" LIMIT 10`;
      break;
    case "message":
      searchQuery = `SELECT * FROM commitLog_${repoId} WHERE commit_message LIKE "%${searchKey}%" LIMIT 10`;
      break;
    case "user":
      searchQuery = `SELECT * FROM commitLog_${repoId} WHERE author LIKE "%${searchKey}%" LIMIT 10`;
      break;
    default:
      searchQuery = `SELECT * FROM commitLog_${repoId} WHERE commit_message LIKE "%${searchKey}%" LIMIT 10`;
      break;
  }

  return new Promise((resolve, reject) => {
    db.all(searchQuery, [], async (err, rows) => {
      if (err) {
        console.log(err);
        reject("Database fetch error");
      }
      if (rows) {
        const hashArray = rows.map((row) => {
          return row.hash;
        });
        const commits = await gitCommitLogSearchHandler(repoId, hashArray)
          .then(async (res) => {
            console.log(await res.commits);
            return Promise.all(res.commits).then((commit) => {
              return commit;
            });
          })
          .catch((err) => {
            console.log(err);
            return [];
          });
        resolve(commits);
      }
    });
  });
}

module.exports.gitCommitLogDbSerchApi = gitCommitLogDbSerchApi;
