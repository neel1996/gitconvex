import { library } from "@fortawesome/fontawesome-svg-core";
import { fab } from "@fortawesome/free-brands-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useMemo, useState } from "react";
import { v4 as uuid } from "uuid";
import { getIconForFile } from "vscode-icons-js";
import { globalAPIEndpoint } from "../../../../../util/env_config";
import InfiniteLoader from "../../../../Animations/InfiniteLoader";
import CodeFileViewComponent from "./RepoDetailBackdrop/CodeFileViewComponent";

export default function FileExplorerComponent(props) {
  library.add(fab, fas);

  const [gitRepoFiles, setGitRepoFiles] = useState([]);
  const [codeViewToggle, setCodeViewToggle] = useState(false);
  const [gitFileBasedCommits, setGitFileBasedCommits] = useState([]);
  const [directoryNavigator, setDirectoryNavigator] = useState([]);
  const [codeViewItem, setCodeViewItem] = useState("");
  const [selectionIndex, setSelectionIndex] = useState(0);
  const [isEmpty, setIsEmpty] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [cwd, setCwd] = useState("");

  const { repoIdState } = props;

  const memoizedCodeFileViewComponent = useMemo(() => {
    return (
      <CodeFileViewComponent
        repoId={repoIdState}
        fileItem={codeViewItem}
        commitMessage={gitFileBasedCommits[selectionIndex]}
      ></CodeFileViewComponent>
    );
  }, [repoIdState, codeViewItem, gitFileBasedCommits, selectionIndex]);

  function filterNullCommitEntries(gitTrackedFiles, gitFileBasedCommit) {
    let localGitCommits = gitFileBasedCommit;
    let localTrackedFiles = gitTrackedFiles.filter((item, index) => {
      if (item) {
        return true;
      } else {
        localGitCommits[index] = "";
        return false;
      }
    });

    localGitCommits = localGitCommits.filter((commit) => commit);

    setGitRepoFiles([...localTrackedFiles]);
    setGitFileBasedCommits([...localGitCommits]);
  }

  useEffect(() => {
    const repoId = props.repoIdState;
    setIsEmpty(false);
    setIsLoading(true);
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      data: {
        query: `
          query
          {
            gitFolderContent(repoId:"${repoId}", directoryName: ""){
              trackedFiles
              fileBasedCommits   
            }
          }
        `,
      },
    })
      .then((res) => {
        setIsLoading(false);
        const {
          trackedFiles,
          fileBasedCommits,
        } = res.data.data.gitFolderContent;

        if (trackedFiles.length === 0 || fileBasedCommits.length === 0) {
          setIsEmpty(true);
          return;
        }

        if (trackedFiles && fileBasedCommits) {
          filterNullCommitEntries(trackedFiles, fileBasedCommits);
        }
      })
      .catch((err) => {
        console.log(err);
        setIsLoading(false);
      });
  }, [props]);

  function directorySeparatorRemover(directoryPath) {
    if (directoryPath.match(/.\/./gi)) {
      directoryPath = directoryPath.split("/")[
        directoryPath.split("/").length - 1
      ];
    } else if (directoryPath.match(/[^\\]\\[^\\]/gi)) {
      directoryPath = directoryPath.split("\\")[
        directoryPath.split("\\").length - 1
      ];
    } else if (directoryPath.match(/.\\\\./gi)) {
      directoryPath = directoryPath.split("\\\\")[
        directoryPath.split("\\\\").length - 1
      ];
    }

    return directoryPath;
  }

  const fetchFolderContent = (
    directoryName,
    slicePosition,
    sliceIndicator,
    homeIndicator
  ) => {
    if (repoIdState) {
      setGitRepoFiles([]);
      setGitFileBasedCommits([]);
      let localDirNavigator = directoryNavigator;

      if (sliceIndicator) {
        let slicedDirectory = localDirNavigator.slice(0, slicePosition);
        if (slicedDirectory.length > 0) {
          directoryName = slicedDirectory.join("/") + "/" + directoryName;
        }
      }

      setCwd(directoryName);
      setIsLoading(true);

      axios({
        url: globalAPIEndpoint,
        method: "POST",
        headers: {
          "Content-type": "application/json",
        },
        data: {
          query: `
            query
            {
              gitFolderContent(repoId:"${repoIdState}", directoryName: "${directoryName}"){
                trackedFiles
                fileBasedCommits   
              }
            }
          `,
        },
      })
        .then((res) => {
          setIsLoading(false);
          if (res.data.data && !res.data.error) {
            const localFolderContent = res.data.data.gitFolderContent;

            filterNullCommitEntries(
              localFolderContent.trackedFiles,
              localFolderContent.fileBasedCommits
            );

            directoryName = directorySeparatorRemover(directoryName);

            if (homeIndicator) {
              setDirectoryNavigator([]);
              return;
            }

            if (directoryNavigator.length === 0) {
              setDirectoryNavigator([directoryName]);
            } else {
              if (
                sliceIndicator &&
                slicePosition < directoryNavigator.length - 1
              ) {
                const iterator =
                  directoryNavigator.length - (slicePosition + 1);

                for (let i = 0; i < iterator; i++) {
                  localDirNavigator.pop();
                }
                setDirectoryNavigator([...localDirNavigator]);
              } else {
                setDirectoryNavigator([...directoryNavigator, directoryName]);
              }
            }
          } else {
            setIsLoading(false);
            console.log(
              "ERROR: Error occurred while fetching the folder content!"
            );
          }
        })
        .catch((err) => {
          setIsLoading(false);
          if (err) {
            console.log(
              "ERROR: Error occurred while fetching the folder content!",
              err
            );
          }
        });
    }
  };

  const gitTrackedFileComponent = () => {
    var fileIcon;

    if (gitRepoFiles && gitRepoFiles.length > 0) {
      var formattedFiles = [];
      var directoryEntry = [];
      var fileEntry = [];

      gitRepoFiles.forEach(async (entry, index) => {
        const splitEntry = entry.split(":");

        if (splitEntry[1] && splitEntry[1].includes("directory")) {
          let directoryPath = directorySeparatorRemover(splitEntry[0]);

          directoryEntry.push(
            <div
              className="block w-full p-2 border-b border-gray-200"
              key={`directory-key-${uuid()}`}
            >
              <div className="flex cursor-pointer items-center">
                <div className="w-1/6">
                  <FontAwesomeIcon
                    icon={["fas", "folder"]}
                    className="font-sans text-xl"
                  ></FontAwesomeIcon>
                </div>
                <div
                  className="font-sans w-1/2 text-gray-600 text-xl mx-3 hover:text-indigo-400"
                  onClick={(event) => {
                    fetchFolderContent(splitEntry[0], 0, false);
                  }}
                >
                  {directoryPath}
                </div>

                <div className="w-3/5 p-2 overflow-hidden truncate rounded-lg mx-auto text-left bg-green-200 text-green-900">
                  {gitFileBasedCommits[index]}
                </div>
              </div>
            </div>
          );
        } else if (splitEntry[1] && splitEntry[1].includes("File")) {
          if (splitEntry[0] === "LICENSE") {
            fileIcon = require("../../../../../assets/icons/file_type_license.svg");
          } else {
            fileIcon = require("../../../../../assets/icons/" +
              getIconForFile(splitEntry[0]));
          }
          fileEntry.push(
            <div
              className="block w-full p-2 border-b border-gray-200"
              key={`file-key-${uuid()}`}
            >
              <div className="flex items-center align-middle cursor-pointer">
                <div className="w-1/6">
                  <img
                    src={fileIcon.default}
                    style={{
                      width: "26px",
                      filter: "grayscale(30%)",
                    }}
                    alt={fileIcon.default}
                  ></img>
                </div>
                <div
                  className="font-sans w-1/2 text-gray-700 text-xl mx-3 hover:text-indigo-400"
                  onClick={() => {
                    setSelectionIndex(index);
                    if (cwd === "" || cwd === "/") {
                      setCodeViewItem(splitEntry[0]);
                    } else {
                      setCodeViewItem(cwd + "/" + splitEntry[0]);
                    }
                    setCodeViewToggle(true);
                  }}
                >
                  {splitEntry[0]}
                </div>
                <div className="w-3/5 p-2 overflow-hidden truncate rounded-lg mx-auto text-left bg-indigo-200 text-indigo-900">
                  {gitFileBasedCommits[index]}
                </div>
              </div>
            </div>
          );
        }
      });

      formattedFiles.push(directoryEntry);
      formattedFiles.push(fileEntry);

      return (
        <div
          className="block mx-auto justify-center p-2"
          style={{
            color: "#79b8ff",
          }}
          key="repo-key"
        >
          <div className="flex w-full justify-around p-2 mx-auto pb-4 border-b border-blue-400">
            <div className="w-1/6"></div>
            <div className="w-2/4 font-sans font-semibold">
              File / Directory
            </div>
            <div className="w-2/4 font-sans font-semibold">Latest commit</div>
          </div>
          {formattedFiles}
        </div>
      );
    }
  };

  return (
    <>
      {isLoading ? (
        <>
          <div className="flex justify-center mx-auto my-2 w-3/4">
            <div className="w-full mx-auto text-2xl text-center font-sans font-semibold text-gray-700 border-b-2 border-dashed border-gray-400">
              Loading tracked files...
            </div>
          </div>
          <div className="flex mx-auto my-6 text-center justify-center">
            <InfiniteLoader
              loadAnimation={!gitRepoFiles.length}
            ></InfiniteLoader>
          </div>
        </>
      ) : (
        <>
          {codeViewToggle ? (
            <div
              className="fixed flex w-full h-full top-0 left-0 right-0 overflow-auto"
              id="code-view__backdrop"
              style={{ background: "rgba(0,0,0,0.5)", zIndex: 99 }}
              onClick={(event) => {
                if (event.target.id === "code-view__backdrop") {
                  setCodeViewToggle(false);
                }
              }}
            >
              <div
                className="w-14 h-14 mr-5 mt-6 rounded-full bg-red-500 text-white flex justify-center items-center shadow cursor-pointer fixed right-0 top-0"
                onClick={() => {
                  setCodeViewToggle(false);
                }}
              >
                <FontAwesomeIcon
                  className="flex text-center text-3xl my-auto"
                  icon={["fas", "times"]}
                ></FontAwesomeIcon>
              </div>
              <div className="w-full h-full mx-auto my-auto block mt-4 mb-10">
                {memoizedCodeFileViewComponent}
              </div>
            </div>
          ) : null}

          <div>
            <div
              className="w-1/6 text-gray-600 cursor-pointer border-b-2 border-dashed text-center flex justify-center gap-4 p-3 rounded mx-6 myy-auto items-center align-middle text-xl"
              onClick={() => {
                fetchFolderContent("", 0, false, true);
              }}
            >
              <div>
                <FontAwesomeIcon icon={["fas", "home"]}></FontAwesomeIcon>
              </div>
              <div>Home</div>
              <div className="text-2xl font-sans text-blue-400">./</div>
            </div>
            {directoryNavigator && gitRepoFiles && gitRepoFiles.length > 0 ? (
              <div className="mx-6 p-3 flex font-sans justify-start items-center gap-4">
                <div
                  className="flex p-4 w-3/4 gap-4 items-center break-words overflow-x-auto"
                  id="repoFolderNavigator"
                >
                  {directoryNavigator.map((item, index) => {
                    return (
                      <div
                        className="flex items-center justify-start gap-2"
                        key={item + "-" + index}
                      >
                        <div
                          className={`${
                            index !== directoryNavigator.length - 1
                              ? "text-blue-400 cursor-pointer font-semibold hover:text-blue-500 underline"
                              : ""
                          } text-xl`}
                          onClick={() => {
                            if (index !== directoryNavigator.length - 1) {
                              fetchFolderContent(item, index, true);
                            }
                          }}
                        >
                          {item}
                        </div>
                        <div>/</div>
                      </div>
                    );
                  })}
                </div>
              </div>
            ) : null}

            <div className="block w-11/12 my-6 mx-auto justify-center p-6 rounded-lg bg-white shadow overflow-auto border">
              {!isEmpty ? (
                gitTrackedFileComponent()
              ) : (
                <div className="flex w-3/4 gap-4 mx-auto items-center align-middle justify-center rounded-lg text-gray-500 text-center border-b-4 text-2xl border-dashed p-1 border-gray-300">
                  <div>
                    <FontAwesomeIcon icon={["fas", "unlink"]}></FontAwesomeIcon>
                  </div>
                  <div>No Tracked Files in the directory!</div>
                </div>
              )}
            </div>
          </div>
        </>
      )}
    </>
  );
}
