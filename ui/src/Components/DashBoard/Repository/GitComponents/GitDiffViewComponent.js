import { LangLine } from "@itassistors/langline";
import axios from "axios";
import * as Prism from "prismjs";
import React, { useContext, useEffect, useState } from "react";
import { v4 as uuidv4 } from "uuid";
import { GIT_TRACKED_FILES } from "../../../../actionStore";
import { ContextProvider } from "../../../../context";
import "../../../../prism.css";
import { globalAPIEndpoint } from "../../../../util/env_config";
import "prismjs/components/prism-markdown";

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
  const [diffWidth, setDiffWidth] = useState(0);

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
          <div
            className="my-2 block cursor-pointer border border-dashed border-blue-800"
            title={fileEntry}
          >
            <div className="bg-gray-50 p-1 rounded">
              <div className="font-sans font-semibold text-blue-400">
                Directory:
              </div>
              <div className="text-gray-500 truncate font-light">
                {dirSplit}
              </div>
            </div>
            <div className="mx-1 my-2 text-sm font-sans text-blue-800">
              {fileName}
            </div>
          </div>
        );
      } else {
        return (
          <span className="text-lg p-2 border-b-2 border-gray-200 w-full">
            {fileEntry}
          </span>
        );
      }
    };

    return (
      <>
        {changedFiles.length >= 1 &&
          changedFiles.map((entry) => {
            if (entry && entry.split(",")[0] === "M") {
              let fileEntry = entry.split(",")[1];
              const styleSelector = " bg-indigo-50 border-b border-indigo-300";
              return (
                <div
                  className={`p-2 text-sm break-words hover:bg-indigo-50 cursor-pointer ${
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
                <div className="font-sans text-lg text-center font-light text-yellow-700">
                  Note : Blank space changes will not be displayed here!
                </div>
              </>
            );
          } else {
            setWarnStatus("");
            setDiffStatState(diffStat);
            setFileLineDiffState(fileDiff);

            // Logic to set width for code lines based on max scroll width
            const target = document.getElementById("codeDiffWrapper");
            const targetWidth = target.scrollWidth;
            setDiffWidth(targetWidth);

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
              className="flex items-center gap-1 bg-green-100 w-screen"
              key={`${line}-${uuidv4()}`}
              style={diffWidth ? { width: `${diffWidth}px` } : null}
            >
              <div className="font-sans text-center w-16 mx-2 border-r border-green-400 text-green-500">
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
              className="flex items-center gap-1 bg-red-100 w-screen"
              key={`${line}-${uuidv4()}`}
              style={diffWidth ? { width: `${diffWidth}px` } : null}
            >
              <div className="font-sans text-center w-16 mx-2 border-r border-red-400 text-red-400">
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
                className="flex items-center gap-1 bg-white-200 w-screen"
                key={`${line}-${uuidv4()}`}
                style={diffWidth ? { width: `${diffWidth}px` } : null}
              >
                <div className="font-sans text-gray-300 text-center w-16 mx-2 border-r border-gray-200">
                  {++lineCounter}
                </div>
                <pre className="w-5/6 mx-2">
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
          <div className="w-full flex justify-center mx-auto">
            <div
              className="overflow-auto break-words p-3 bg-white border-2 border-dashed border-gray-300 w-1/4 rounded-lg shadow-md"
              style={{ height: "880px" }}
            >
              {getDiffFiles()}
            </div>

            {!activeFileName ? (
              <div className="mt-4 p-4 mb-auto bg-gray-200 rounded shadow w-1/2 font-sans font-semibold mx-auto flex justify-center text-center">
                Click on a file to see difference
              </div>
            ) : (
              ""
            )}

            {warnStatus ? (
              <div className="mx-auto my-auto block bg-yellow-200 text-yellow-700 p-5 text-xl font-semibold rounded-lg shadow border border-b-2 border-dashed border-yellow-300">
                {warnStatus}
              </div>
            ) : null}

            {diffStatState &&
            diffStatState !== "Click on a file item to see the difference" &&
            !warnStatus ? (
              <div className="w-3/4 mx-auto my-auto break-all p-3">
                {diffStatState ? statFormat() : ""}
                {fileLineDiffState &&
                fileLineDiffState !==
                  "Click on a file item to see the difference" ? (
                  <div
                    className="p-3 overflow-auto break-words w-full"
                    id="codeDiffWrapper"
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
            <div className="mx-auto bg-yellow-100 w-full my-2 p-3 border-dashed rounded-lg shadow border-b-4 border-yellow-400 text-xl font-sans font-semibold text-center text-yellow-600">
              No File changes in the repo
            </div>
          ) : (
            <div className="w-full p-4 mx-auto text-center font-sans font-semibold text-xl rounded-lg shadow my-2 bg-pink-50 border-b-4 border-dashed border-pink-400">
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
