const fs = require("fs");
const path = require("path");
const { getRepoPath } = require("../global/fetchGitRepoPath");
const { gitFileBasedCommit } = require("../git/gitFileBasedCommit");
const { isText } = require("istextorbinary");
const { LangLine } = require("@itassistors/langline");

async function codeFileViewApi(repoId, fileName) {
  const repoPath = await getRepoPath(repoId);
  const targetFile = path.join(repoPath, fileName);
  const langData = await new LangLine().withFile(targetFile);
  let fileContent = [];

  if (isText(targetFile)) {
    const commit = await gitFileBasedCommit(repoPath, targetFile);

    let fileData = await fs.promises
      .readFile(targetFile)
      .then((res) => {
        return res.toString();
      })
      .catch((err) => {
        console.log(err);
        return "";
      });

    if (langData && langData.name) {
      if (fileData) {
        fileContent = fileData.split("\n");
        return {
          codeFileDetails: {
            language: langData.name,
            fileData: fileContent,
            fileCommit: commit,
            prism: langData.prismIndicator
              ? langData.prismIndicator
              : "markdown",
          },
        };
      }
    } else {
      return {
        codeFileDetails: {
          language: "",
          fileData: fileContent,
          fileCommit: commit,
          prism: "markdown",
        },
      };
    }
  }
}

module.exports.codeFileViewApi = codeFileViewApi;
