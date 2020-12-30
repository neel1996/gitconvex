import { library } from "@fortawesome/fontawesome-svg-core";
import { fab } from "@fortawesome/free-brands-svg-icons";
import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../../util/env_config";
import "../../../../styles/GitOperations.css";
import CommitComponent from "./CommitComponent";
import PushComponent from "./PushComponent";
import StageComponent from "./StageComponent";

export default function GitOperationComponent(props) {
  library.add(fab);
  const { repoId } = props;
  const { stateChange } = props;

  const [gitTrackedFiles, setGitTrackedFiles] = useState([]);
  const [gitUntrackedFiles, setGitUntrackedFiles] = useState([]);

  const [action, setAction] = useState("");
  const [list, setList] = useState([]);
  const [viewReload, setViewReload] = useState(0);
  const [currentStageItem, setCurrentStageItem] = useState("");
  const [stageItems, setStagedItems] = useState([]);
  const [unStageFailed, setUnStageFailed] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    setIsLoading(true);
    setStagedItems([]);
    setCurrentStageItem("");

    const cancelToken = axios.CancelToken;
    const source = cancelToken.source();

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      cancelToken: source.token,
      data: {
        query: `
            query {
                gitChanges(repoId: "${props.repoId}"){
                  gitUntrackedFiles
                  gitChangedFiles
                  gitStagedFiles
                }
            }
        `,
      },
    })
      .then((res) => {
        if (res.data.data) {
          var apiData = res.data.data.gitChanges;

          setGitTrackedFiles([...apiData.gitChangedFiles]);
          setGitUntrackedFiles([...apiData.gitUntrackedFiles]);
          setStagedItems([...apiData.gitStagedFiles]);

          const apiTrackedFiles = [...apiData.gitChangedFiles];
          const apiUnTrackedFiles = [...apiData.gitUntrackedFiles];

          let componentList = [];

          apiTrackedFiles &&
            apiTrackedFiles.forEach((item) => {
              if (item.split(",").length > 0) {
                const trackedItem = item.split(",")[1];
                componentList.push(trackedItem);
              }
            });

          apiUnTrackedFiles &&
            apiUnTrackedFiles.forEach((item) => {
              if (item) {
                item = item.replace("NO_DIR", "");
                item.split(",")
                  ? componentList.push(item.split(",").join(""))
                  : componentList.push(item);
              }
            });

          setList([...componentList]);
          setIsLoading(false);
        }
      })
      .catch((err) => {
        setIsLoading(false);
      });

    return () => {
      source.cancel();
    };
  }, [props.repoId, viewReload, currentStageItem]);

  const actionButtons = [
    {
      label: "Stage all changes",
      border: "border-indigo-300",
      text: "text-indigo-700",
      bg: "bg-indigo-50",
      hoverBg: "hover:bg-indigo-100",
      key: "stage",
    },
    {
      label: "Commit Changes",
      border: "border-green-300",
      text: "text-green-700",
      bg: "bg-green-50",
      hoverBg: "hover:bg-green-100",
      key: "commit",
    },
    {
      label: "Push to remote",
      border: "border-pink-300",
      text: "text-pink-700",
      bg: "bg-pink-50",
      hoverBg: "hover:bg-pink-100",
      key: "push",
    },
  ];

  const tableColumns = ["CHANGES", "FILE STATUS", "ACTION"];

  function stageGitComponent(stageItem, event) {
    let localViewReload = viewReload + 1;

    event.target.innerHTML = "Staging...";
    event.target.style.backgroundColor = "gray";
    event.target.disabled = true;

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          mutation {
            stageItem(repoId: "${repoId}", item: "${stageItem}")
          }
        `,
      },
    })
      .then((res) => {
        setViewReload(localViewReload);

        if (res.data.data && !res.data.error) {
          if (res.data.data.stageItem === "ADD_ITEM_SUCCESS") {
            setCurrentStageItem(stageItem);
          }
        }
      })
      .catch((err) => {
        console.log(err);
        setViewReload(localViewReload);
      });
  }

  function getTableData() {
    let tableDataArray = [];

    let statusPill = (status) => {
      if (status === "M") {
        return (
          <div className="git-ops--file--status w-full p-3 text-yellow-600 border-2 border-yellow-200 font-sans font-semibold bg-yellow-50">
            MODIFIED
          </div>
        );
      } else if (status === "D") {
        return (
          <div className="git-ops--file--status w-full p-3 text-red-500 border-2 border-red-200 font-sans font-semibold bg-red-50">
            REMOVED
          </div>
        );
      } else {
        return (
          <div className="git-ops--file--status w-full p-3 text-indigo-600 border-2 border-indigo-200 font-sans font-semibold bg-indigo-50">
            UNTRACKED
          </div>
        );
      }
    };

    let actionButton = (stageItem) => {
      return (
        <div
          className="git-ops--stageitem--add p-3 bg-green-300 text-xl font-sans font-semibold text-center shadow-md hover:bg-green-400 hover:shadow-sm"
          onClick={(event) => {
            stageGitComponent(stageItem, event);
            setUnStageFailed(false);
          }}
          key={`add-btn-${stageItem}`}
        >
          ADD
        </div>
      );
    };

    gitTrackedFiles &&
      gitTrackedFiles.forEach((item) => {
        if (item.split(",").length > 0) {
          const trackedItem = item.split(",")[1];
          tableDataArray.push([
            trackedItem,
            statusPill(item.split(",")[0]),
            actionButton(trackedItem),
          ]);
        }
      });

    gitUntrackedFiles &&
      gitUntrackedFiles.forEach((item) => {
        if (item) {
          item = item.replace("NO_DIR", "");
          item.split(",")
            ? tableDataArray.push([
                item.split(",").join(""),
                statusPill("N"),
                actionButton(item.split(",").join("")),
              ])
            : tableDataArray.push([item, statusPill("N"), actionButton(item)]);
        }
      });
    return tableDataArray;
  }

  function getStagedFilesComponent() {
    function removeStagedItem(item, event) {
      let localViewReload = viewReload + 1;
      stateChange();

      event.target.innerHTML = "removing...";
      event.target.style.backgroundColor = "gray";
      event.target.disabled = true;

      axios({
        url: globalAPIEndpoint,
        method: "POST",
        data: {
          query: `
            mutation{
              removeStagedItem(repoId: "${repoId}", item: "${item}")
            }
          `,
        },
      })
        .then((res) => {
          setViewReload(localViewReload);
          if (res.data.data && !res.data.error) {
            if (res.data.data.removeStagedItem === "STAGE_REMOVE_SUCCESS") {
              let localStagedItems = stageItems;

              localStagedItems = localStagedItems.filter((filterItem) => {
                if (filterItem === item) {
                  return false;
                }
                return true;
              });

              setStagedItems([...localStagedItems]);
            } else {
              setUnStageFailed(true);
            }
          }
        })
        .catch((err) => {
          console.log(err);
          setViewReload(localViewReload);
          setUnStageFailed(true);
        });
    }

    function removeAllStagedItems(event) {
      let localViewReload = viewReload + 1;
      setStagedItems();
      stateChange();

      event.target.innerHTML = "Removing...";
      event.target.style.backgroundColor = "gray";
      event.target.disabled = true;

      axios({
        url: globalAPIEndpoint,
        method: "POST",
        data: {
          query: `
            mutation {
              removeAllStagedItem(repoId: "${repoId}")
            }
          `,
        },
      })
        .then((res) => {
          setViewReload(localViewReload + 1);
          if (res.data.data && !res.data.error) {
            if (
              res.data.data.removeAllStagedItem === "STAGE_ALL_REMOVE_SUCCESS"
            ) {
              setStagedItems([]);
            }
          }
        })
        .catch((err) => {
          console.log(err);
          setViewReload(localViewReload + 1);
        });
    }

    if (stageItems && stageItems.length > 0) {
      return (
        <div className="git-ops--stagearea border border-blue-100 shadow-md rounded-lg">
          <div className="git-ops--stagearea--top items-center">
            <div className="text-5xl font-sans text-gray-800 mx-4">
              Staged Files
            </div>
            <div
              className="git-ops--unstage--btn"
              onClick={(event) => {
                removeAllStagedItems(event);
                setUnStageFailed(false);
              }}
            >
              Remove All Items
            </div>
          </div>
          {unStageFailed ? (
            <div className="my-4 mx-auto text-center shadow-md rounded p-4 border border-red-200 text-red-400 font-sans font-semibold">
              Remove item failed. Note that deleted files cannot be removed as a
              single entity. Use
              <i className="border-b border-gray-600 border-dashed mx-2">
                Remove All Items
              </i>
              to perform a complete git reset
            </div>
          ) : null}
          <div className="git-ops--unstage--header">
            <div className="w-11/12 font-sans font-semibold text-xl text-center text-white border-r-2 border-indigo-300">
              STAGED FILE
            </div>
            <div className="w-1/4 font-sans font-semibold text-xl text-center text-white border-r-2 border-indigo-300">
              ACTION
            </div>
          </div>
          <div className="git-ops--unstage--table" style={{ height: "450px" }}>
            {stageItems.map((item) => {
              if (item) {
                return (
                  <div className="git-ops--unstage--table--data" key={item}>
                    <div className="w-11/12 block px-1 border-r-2 border-gray-300">
                      <div className="git-changed-item overflow-x-auto w-full font-sans font-light text-base text-gray-700">
                        {item}
                      </div>
                    </div>
                    <div className="w-1/4 mx-auto">
                      <div
                        className="git-ops--unstage--remove--btn p-3 bg-red-400 text-xl font-sans font-semibold text-center shadow-md hover:bg-red-500 hover:shadow-sm"
                        onClick={(event) => {
                          removeStagedItem(item, event);
                          setUnStageFailed(false);
                        }}
                        key={`remove-btn-${item}`}
                      >
                        REMOVE
                      </div>
                    </div>
                  </div>
                );
              }
              return null;
            })}
          </div>
        </div>
      );
    }
  }

  function actionComponent(action) {
    switch (action) {
      case "stage":
        if (list && list.length > 0) {
          return (
            <StageComponent
              repoId={repoId}
              stageComponents={list}
            ></StageComponent>
          );
        } else {
          return (
            <div className="w-1/2 mx-auto my-auto bg-gray-200 p-6 rounded-md">
              <div className="bg-white p-6 font-sans text-3xl font-light text-gray-500 border-b-4 border-dashed rounded-lg shadow-lg border-gray-500 text-center">
                No Changes for staging...
              </div>
            </div>
          );
        }
      case "commit":
        return <CommitComponent repoId={repoId}></CommitComponent>;
      case "push":
        return <PushComponent repoId={repoId}></PushComponent>;
      default:
        return null;
    }
  }

  function noChangesComponent() {
    return (
      <div className="git-ops--nochange">
        {!isLoading ? (
          <span>No files changes found in the repo.</span>
        ) : (
          <span className="text-gray-600">Fetching results...</span>
        )}
      </div>
    );
  }

  return (
    <>
      {action ? (
        <div
          className="git-ops--backdrop"
          id="operation-backdrop"
          style={{ background: "rgba(0,0,0,0.6)", zIndex: "99" }}
          onClick={(event) => {
            if (event.target.id === "operation-backdrop") {
              setAction("");
              let closeViewCount = viewReload + 1;
              setViewReload(closeViewCount);
              setStagedItems([]);
              setGitTrackedFiles([]);
              setGitUntrackedFiles([]);
              setList([]);
            }
          }}
        >
          {actionComponent(action)}

          <div
            className="git-ops--backdrop--close"
            onClick={() => {
              setAction("");
              const localReload = viewReload + 1;
              setViewReload(localReload);
            }}
          >
            X
          </div>
        </div>
      ) : null}
      <div className="flex justify-between w-full mx-auto items-center my-10">
        {actionButtons.map((item) => {
          const { label, border, text, hoverBg, key, bg } = item;
          return (
            <div
              className={`mx-4 shadow-lg w-1/3 text-center font-sans text-xl font-semibold rounded-lg p-3 border-b-4 rounded-b-xl transition hover:shadow-sm cursor-pointer ${bg} ${border} ${text} ${hoverBg}`}
              key={key}
              onClick={() => setAction(key)}
            >
              {label}
            </div>
          );
        })}
      </div>
      {getTableData() && getTableData().length > 0 ? (
        <div className="git-ops--file-table">
          <div className="flex items-center p-4 bg-gray-100 rounded">
            {tableColumns.map((column, index) => {
              return (
                <div
                  key={column}
                  className={`font-sans font-semibold text-xl text-center text-gray-600 border-r-2 border-white  ${
                    index === 0 ? "w-3/4" : "w-1/4"
                  }`}
                >
                  {column}
                </div>
              );
            })}
          </div>

          <div
            className="git-ops--file-table--data"
            style={{ height: "400px" }}
          >
            <div>
              {isLoading ? (
                <div className="text-center font-sans font-light p-4 text-2xl text-gray-600">
                  Loading modified file items...
                </div>
              ) : (
                <>
                  {getTableData() &&
                    getTableData().map((tableData, index) => {
                      return (
                        <div
                          className="git-ops--file-table--items"
                          key={`tableItem-${index}`}
                        >
                          {tableData.map((data, index) => {
                            return (
                              <div
                                key={`${data}-${index}`}
                                title={data}
                                className={`git-changed-item mx-1 overflow-x-auto font-sans text-gray-700 items-center align-middle my-auto ${
                                  index === 0
                                    ? "w-3/4 text-left"
                                    : "w-1/4 text-center"
                                }`}
                              >
                                {data}
                              </div>
                            );
                          })}
                        </div>
                      );
                    })}
                </>
              )}
            </div>
          </div>
        </div>
      ) : (
        <>{noChangesComponent()}</>
      )}
      <>{stageItems && !isLoading ? getStagedFilesComponent() : null}</>
    </>
  );
}
