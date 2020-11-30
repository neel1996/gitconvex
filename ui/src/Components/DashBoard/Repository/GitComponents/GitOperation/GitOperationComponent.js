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
      border: "border-blue-500",
      text: "text-blue-700",
      hoverBg: "bg-blue-500",
      key: "stage",
    },
    {
      label: "Commit Changes",
      border: "border-green-500",
      text: "text-green-700",
      hoverBg: "bg-green-500",
      key: "commit",
    },
    {
      label: "Push to remote",
      border: "border-pink-500",
      text: "text-pink-700",
      hoverBg: "bg-pink-500",
      key: "push",
    },
  ];

  const tableColumns = ["Changes", "File Status", "Action"];

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
          <div className="git-ops--file--status text-yellow-700 border-yellow-500 ">
            Modified
          </div>
        );
      } else if (status === "D") {
        return (
          <div className="git-ops--file--status text-red-700 bg-red-200 border-red-500 ">
            Removed
          </div>
        );
      } else {
        return (
          <div className="git-ops--file--status text-indigo-700 border-indigo-500 ">
            Untracked
          </div>
        );
      }
    };

    let actionButton = (stageItem) => {
      return (
        <div
          className="git-ops--stageitem--add"
          onClick={(event) => {
            stageGitComponent(stageItem, event);
            setUnStageFailed(false);
          }}
          key={`add-btn-${stageItem}`}
        >
          Add
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
        <div className="git-ops--stagearea">
          <div className="git-ops--stagearea--top">
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
            <div className="w-3/4">Staged File</div>
            <div className="w-1/2">Action</div>
          </div>
          <div className="git-ops--unstage--table" style={{ height: "450px" }}>
            {stageItems.map((item) => {
              if (item) {
                return (
                  <div className="git-ops--unstage--table--data" key={item}>
                    <div className="git-ops--unstage--table--item">{item}</div>
                    <div className="w-1/2 mx-auto">
                      <div
                        className="git-ops--unstage--remove--btn"
                        onClick={(event) => {
                          removeStagedItem(item, event);
                          setUnStageFailed(false);
                        }}
                        key={`remove-btn-${item}`}
                      >
                        Remove
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
      <div className="git-ops--actions">
        {actionButtons.map((item) => {
          const { label, border, text, hoverBg, key } = item;
          return (
            <div
              className={`git-ops--actions--btn ${border} ${text} hover:${hoverBg}`}
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
          <div className="git-ops--file-table--header">
            {tableColumns.map((column, index) => {
              return (
                <div
                  key={column}
                  className={`git-ops--file-table--cols ${
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
                                className={`break-all items-center align-middle my-auto ${
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
