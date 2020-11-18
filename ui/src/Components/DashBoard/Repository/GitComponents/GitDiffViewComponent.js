import { LangLine } from "@itassistors/langline";
import axios from "axios";
import * as Prism from "prismjs";
import React, { useContext, useEffect, useState } from "react";
import { v4 as uuidv4 } from "uuid";
import { GIT_TRACKED_FILES } from "../../../../actionStore";
import { ContextProvider } from "../../../../context";
import "../../../../prism.css";
import { globalAPIEndpoint } from "../../../../util/env_config";
import "../../../styles/GitDiffView.css";

export default function GitDiffViewComponent() {
  const { state, dispatch } = useContext(ContextProvider);
  const repoId = state.globalRepoId;

  const [changedFiles, setChangedFiles] = useState([]);
  const [diffStatState, setDiffStatState] = useState(
    "Click on a file item to see the difference"
  );
  const [fileLineDiffState, setFileLineDiffState] = useState([]);
  const [activeFileName, setActiveFileName] = useState("");
  const [isApiCalled, setIsApiCalled] = useState(false);
  const [warnStatus, setWarnStatus] = useState("");
  const [lang, setLang] = useState("");

  useEffect(() => {
    let repoId = state.globalRepoId;

    setActiveFileName("");
    setFileLineDiffState("Click on a file item to see the difference");
    setDiffStatState("Click on a file item to see the difference");
    setWarnStatus("");
    setChangedFiles([]);

    let apiEndPoint = globalAPIEndpoint;
    if (repoId) {
      axios({
        url: apiEndPoint,
        method: "POST",
        headers: {
          "Content-type": "application/json",
        },
        data: {
          query: `
            query {
              gitChanges(repoId: "${repoId}"){
                gitChangedFiles
              }
            }
          `,
        },
      })
        .then((res) => {
          if (res) {
            if (res.data.data && !res.data.error) {
              var apiData = res.data.data.gitChanges;
              let { gitChangedFiles } = apiData;

              gitChangedFiles = gitChangedFiles.filter((fileEntry) => {
                if (fileEntry.split(",")[0] === "D") {
                  return false;
                }
                return true;
              });

              if (gitChangedFiles.length >= 1) {
                setChangedFiles([...gitChangedFiles]);
              } else {
                setWarnStatus(
                  "No modified or new files found in the repo. All the files are either removed or not present in the repo!"
                );
              }

              setIsApiCalled(true);
              dispatch({ type: GIT_TRACKED_FILES, payload: gitChangedFiles });
            }
          }
        })
        .catch((err) => {
          return err;
        });
    }
  }, [state.globalRepoId, dispatch]);

  function getDiffFiles() {
    const directorySplit = (fileEntry) => {
      if (fileEntry.includes("/")) {
        let splitEntry = fileEntry.split("/");
        let dirSplit = splitEntry
          .map((entry, index) => {
            if (index === splitEntry.length - 1) {
              return null;
            } else {
              return entry;
            }
          })
          .join("/");
        let fileName = splitEntry[splitEntry.length - 1];

        return (
          <div className="git-diff--files--list" title={fileEntry}>
            <div className="bg-gray-100 p-1 rounded">
              <div className="git-diff--files--label">Directory:</div>
              <div className="git-diff--files--directory">{dirSplit}</div>
            </div>
            <div className="git-diff--files--filename_sm">{fileName}</div>
          </div>
        );
      } else {
        return (
          <span className="git-diff--files--filename_lg">{fileEntry}</span>
        );
      }
    };

    return (
      <>
        {changedFiles.length >= 1 &&
          changedFiles.map((entry) => {
            if (entry && entry.split(",")[0] === "M") {
              let fileEntry = entry.split(",")[1];
              const styleSelector = " bg-indigo-100 border-b border-indigo-400";
              return (
                <div
                  className={`p-2 text-sm break-words hover:bg-indigo-100 cursor-pointer ${
                    fileEntry === activeFileName ? styleSelector : ""
                  }`}
                  onClick={() => {
                    setActiveFileName(fileEntry);
                    fileDiffStatComponent(repoId, fileEntry);
                  }}
                  key={fileEntry}
                >
                  {directorySplit(fileEntry)}
                </div>
              );
            } else {
              return null;
            }
          })}
      </>
    );
  }

  function fileDiffStatComponent(repoId, fileName) {
    const apiEndPoint = globalAPIEndpoint;
    setWarnStatus("Loading file difference...");

    axios({
      url: apiEndPoint,
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      data: {
        query: `
          query{
              gitFileLineChanges(repoId: "${repoId}", fileName: "${fileName}"){
                diffStat
                fileDiff
              }
          }
        `,
      },
    })
      .then(async (res) => {
        if (res.data.data && !res.data.error) {
          const { diffStat, fileDiff } = res.data.data.gitFileLineChanges;

          if (diffStat === "NO_STAT" || fileDiff[0] === "NO_DIFF") {
            setWarnStatus(
              <>
                <div>
                  No difference could be found. Please check if the file is
                  present.
                </div>
                <div className="font-sans text-sm font-light text-yellow-700">
                  Note : Space based changes will not be displayed here even if
                  it is considered as a change!
                </div>
              </>
            );
          } else {
            setWarnStatus("");
            setDiffStatState(diffStat);
            setFileLineDiffState(fileDiff);

            let language;
            let langName = new LangLine().withFileName(fileName).prismIndicator;

            if (langName && typeof langName != undefined) {
              language = langName;
            } else {
              language = "markdown";
            }

            if (language) {
              await import("prismjs/components/prism-" + language + ".js")
                .then(() => {})
                .catch((err) => {
                  console.log(err);
                  setLang("markdown");
                });

              setLang(language);
            }
          }
        } else {
          setWarnStatus(
            "Error while fetching the file difference. Please try reloading the view!"
          );
        }
      })
      .catch((err) => {
        console.log(err);
        setWarnStatus(
          "Error while fetching the file difference. Please try reloading the view!"
        );
      });
  }

  function statFormat() {
    if (diffStatState) {
      let splitStat = diffStatState.split(",");

      return (
        <div className="text-xl text-center w-3/4 mx-auto block p-2 border border-gray-500 rounded-md shadow-md border-dotted">
          {splitStat.map((parts) => {
            if (parts.match(/insert/i)) {
              return (
                <span key={`${parts}-${new Date().getTime()}`}>
                  <span className="px-2">{parts.toString().split(" ")[0]}</span>
                  <span className="text-green-700 font-sans font-semibold">
                    insertions (+)
                  </span>
                </span>
              );
            } else {
              return (
                <span key={`${parts}-${new Date().getTime()}`}>
                  <span className="px-2">{parts.toString().split(" ")[0]}</span>
                  <span className="text-red-700 font-sans font-semibold">
                    deletions (-)
                  </span>
                </span>
              );
            }
          })}
        </div>
      );
    }
  }

  function fileLineDiffComponent() {
    let splitLines = [];
    let lineCounter = 0;
    if (fileLineDiffState && lang) {
      splitLines = fileLineDiffState.map((line) => {
        if (line.match(/\\ No newline at end of file/gi)) {
          return "";
        }
        if (line[0] && line[0] === "+") {
          return (
            <div
              className="git-diff--codeview--change bg-green-200"
              key={`${line}-${uuidv4()}`}
            >
              <div className="git-diff--codeview--linenumber border-green-500 text-green-500">
                {++lineCounter}
              </div>
              <pre className="w-5/6 mx-2">
                <code
                  dangerouslySetInnerHTML={{
                    __html: Prism.highlight(
                      line.replace("+", ""),
                      Prism.languages[lang],
                      lang
                    ),
                  }}
                ></code>
              </pre>
            </div>
          );
        } else if (line[0] && line[0] === "-") {
          return (
            <div
              className="git-diff--codeview--change bg-red-200"
              key={`${line}-${uuidv4()}`}
            >
              <div className="git-diff--codeview--linenumber border-red-500 text-red-500">
                -
              </div>
              <pre className="w-5/6 mx-2">
                <code
                  dangerouslySetInnerHTML={{
                    __html: Prism.highlight(
                      line.replace("-", ""),
                      Prism.languages[lang],
                      lang
                    ),
                  }}
                ></code>
              </pre>
            </div>
          );
        } else {
          if (line[0]) {
            return (
              <div
                className="git-diff--codeview--change bg-white-200 "
                key={`${line}-${uuidv4()}`}
              >
                <div className="git-diff--codeview--linenumber">
                  {++lineCounter}
                </div>
                <pre className="w-5/6">
                  <code
                    dangerouslySetInnerHTML={{
                      __html: Prism.highlight(
                        line,
                        Prism.languages[lang],
                        lang
                      ),
                    }}
                  ></code>
                </pre>
              </div>
            );
          } else {
            return "";
          }
        }
      });
    }

    return (
      <div className="break-all my-6 mx-auto">
        <code>{splitLines}</code>
      </div>
    );
  }

  return (
    <>
      {changedFiles && changedFiles.length > 0 ? (
        <>
          <div className="git-diff--wrapper">
            <div className="git-diff--files" style={{ height: "880px" }}>
              {getDiffFiles()}
            </div>

            {!activeFileName ? (
              <div className="git-diff--msg">
                Click on a file to see difference information
              </div>
            ) : (
              ""
            )}

            {warnStatus ? (
              <div className="git-diff--warn">{warnStatus}</div>
            ) : null}

            {diffStatState &&
            diffStatState !== "Click on a file item to see the difference" &&
            !warnStatus ? (
              <div className="git-diff--codeview">
                {diffStatState ? statFormat() : ""}
                {fileLineDiffState &&
                fileLineDiffState !==
                  "Click on a file item to see the difference" ? (
                  <div
                    className="git-diff--codeview--content"
                    style={{ height: "800px" }}
                  >
                    {fileLineDiffComponent()}
                  </div>
                ) : (
                  ""
                )}
              </div>
            ) : (
              ""
            )}
          </div>
        </>
      ) : (
        <>
          {isApiCalled ? (
            <div className="git-diff--error">No File changes in the repo</div>
          ) : (
            <div className="git-diff--loader">
              <span className="text-gray-400">
                Fetching changed files from the server...
              </span>
            </div>
          )}
        </>
      )}
    </>
  );
}
