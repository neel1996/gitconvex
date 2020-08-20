const forever = require("forever-monitor");
const path = require("path");

function commitLogCrawler() {
  const foreverMonitor = new forever.Monitor(
    path.join(__dirname, ".", "sqliteDbAccess.js"),
    {
      spinSleepTime: 5000,
      silent: false,
      killTree: true,
    }
  );

  foreverMonitor.on("exit", () => {
    console.log("INFO: Commit log Database module has stopped!");
  });

  foreverMonitor.on("error", (err) => {
    console.log(err);
  });

  foreverMonitor.start();

  setInterval(() => {
    console.log("INFO: Restarting commit log DB crawler");
    foreverMonitor.restart();
  }, 60 * 1000);
}

module.exports.commitLogCrawler = commitLogCrawler;
