import { library } from "@fortawesome/fontawesome-svg-core";
import { fab } from "@fortawesome/free-brands-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useContext, useEffect, useRef, useState } from "react";
import { DELETE_PRESENT_REPO } from "../../../actionStore";
import { ContextProvider } from "../../../context";
import { globalAPIEndpoint } from "../../../util/env_config";

export default function Settings(props) {
  library.add(fab, fas);

  const dbPathTextRef = useRef();

  const { state, dispatch } = useContext(ContextProvider);
  const { presentRepo } = state;

  const [dbPath, setDbPath] = useState("");
  const [port, setPort] = useState(0);
  const [repoDetails, setRepoDetails] = useState([]);
  const [backdropToggle, setBackdropToggle] = useState(false);
  const [deleteRepo, setDeleteRepo] = useState({});
  const [deleteRepoStatus, setDeleteRepoStatus] = useState("");
  const [viewReload, setViewReload] = useState(0);
  const [newDbPath, setNewDbPath] = useState("");
  const [dbUpdateFailed, setDbUpdateFailed] = useState(false);
  const [portUpdateFailed, setPortUpdateFailed] = useState(false);
  const [editMode, setEditMode] = useState({status: false, idx: null});
  const [renameRepo, setRenameRepo] = useState('');

  useEffect(() => {
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          query {
            settingsData{
              settingsDatabasePath
              settingsPortDetails
            }
          }
        `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          const {
            settingsDatabasePath,
            settingsPortDetails,
          } = res.data.data.settingsData;

          setDbPath(settingsDatabasePath);
          setNewDbPath(settingsDatabasePath);
          setPort(settingsPortDetails);
          dbPathTextRef.current.value = settingsDatabasePath;
        }
      })
      .catch((err) => {
        console.log(err);
      });

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            query {
              fetchRepo{
                repoId
                repoName
                repoPath
                timeStamp
              }
            }
          `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          const repoDetails = res.data.data.fetchRepo;
          setRepoDetails(repoDetails);
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }, [props, viewReload]);

  const databasePathSettings = () => {
    const updateDbFileHandler = () => {
      if (newDbPath) {
        const localViewReload = viewReload + 1;

        axios({
          url: globalAPIEndpoint,
          method: "POST",
          data: {
            query: `
              mutation {
                updateRepoDataFile(newDbFile: "${newDbPath.toString()}")
              }
            `,
          },
        })
          .then((res) => {
            if (res.data.data && !res.data.error) {
              const updateStatus = res.data.data.updateRepoDataFile;
              if (updateStatus === "DATAFILE_UPDATE_SUCCESS") {
                setDbUpdateFailed(false);
                setViewReload(localViewReload);
              } else {
                setDbUpdateFailed(true);
                setViewReload(localViewReload);
              }
            } else {
              setDbUpdateFailed(true);
              setViewReload(localViewReload);
            }
          })
          .catch((err) => {
            console.log("Datafile update error", err);
            setDbUpdateFailed(true);
            setViewReload(localViewReload);
          });
      }
    };
    return (
      <div className="settings-data">
        <div className="text-xl text-gray-700 font-sans font-semibold">
          Server data file (file which stores repo details)
        </div>
        <div className="my-4">
          <input
            type="text"
            className="p-2 rounded border border-gray-500 bg-gray-200 text-gray-800 w-2/3"
            ref={dbPathTextRef}
            onChange={(event) => {
              setNewDbPath(event.target.value);
              setDbUpdateFailed(false);
            }}
            onClick={() => {
              setDbUpdateFailed(false);
            }}
          ></input>
          <div className="text-justify font-sand font-light text-sm my-4 text-gray-500 italic w-2/3">
            The data file can be updated. The data file must be an accessible
            JSON file with read / write permissions set to it. Also make sure
            you enter the full path for the file
            <pre className="my-2">E.g: /opt/my_data/data-file.json</pre>
          </div>
          {dbPath !== newDbPath ? (
            <div
              className="my-4 text-center p-2 font-sans text-white border-green-400 border-2 bg-green-500 rounded-md cursor-pointer shadow w-1/4 hover:bg-green-600"
              onClick={() => {
                updateDbFileHandler();
                setDbUpdateFailed(false);
              }}
            >
              Update Data file
            </div>
          ) : null}
          {dbUpdateFailed ? (
            <div className="my-2 p-2 rounded border border-red-300 text-red-700 font-sans font-semibold w-2/3 text-center">
              Data file update failed
            </div>
          ) : null}
        </div>
      </div>
    );
  };

  function deleteRepoHandler() {
    const repoColumn = ["Repo ID", "Repo Name", "Repo Path", "Timestamp"];
    let repoArray = [];

    Object.keys(deleteRepo).forEach((key, index) => {
      repoArray.push({ label: repoColumn[index], value: deleteRepo[key] });
    });

    return (
      <div className="w-3/4 p-6 mx-auto my-auto rounded shadow bg-white">
        <div className="mx-4 my-2 text-3xl font-sans text-gray-900">
          The repo below will be removed from Gitconvex records.
        </div>
        <div className="mx-4 my-1 text-md font-light w-5/6 font-sans italic text-gray-800">
          This will not delete the actual git folder. Just the record from the
          gitconvex server will be removed
        </div>
        <div className="my-2 mx-auto block justify-center w-3/4 p-2">
          {repoArray.map((item) => {
            return (
              <div className="mx-auto flex p-2 font-sans" key={item.label}>
                <div className="w-2/4 font-semibold">{item.label}</div>
                <div className="w-2/4">{item.value}</div>
              </div>
            );
          })}
        </div>

        {deleteRepoStatus !== "lodaing" && deleteRepoStatus !== "success" ? (
          <div
            className="cursor-pointer mx-auto my-4 text-center p-3 rounded shadow bg-red-400 hover:bg-red-500 text-white text-xl"
            onClick={() => {
              deleteRepoApiHandler();
              setDeleteRepoStatus("");
            }}
          >
            Confirm Delete
          </div>
        ) : null}

        {deleteRepoStatus === "loading" ? (
          <div className="cursor-pointer mx-auto my-4 text-center p-3 rounded shadow bg-gray-400 hover:bg-gray-500 text-white text-xl">
            Deletion in progress
          </div>
        ) : null}
        {deleteRepoStatus === "success" ? (
          <div className="p-4 mx-auto text-center font-sans font-semibold bg-green-300 text-green-600 my-4">
            Repo has been deleted!
          </div>
        ) : null}
        {deleteRepoStatus === "failed" ? (
          <div className="p-4 mx-auto text-center font-sans font-semibold bg-red-300 my-4">
            Repo deletion failed!
          </div>
        ) : null}
      </div>
    );
  }

  function deleteRepoApiHandler() {
    setDeleteRepoStatus("loading");
    const { repoId } = deleteRepo;
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          mutation {
            deleteRepo(repoId: "${repoId}"){
              status
              repoId
            }
          }
        `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.eror) {
          const { status, repoId } = res.data.data.deleteRepo;
          if (status === "DELETE_SUCCESS") {
            if (presentRepo && presentRepo.length > 0) {
              let localState = presentRepo[0];

              localState = localState.map((item) => {
                if (item.id.toString() === repoId.toString()) {
                  return null;
                } else {
                  return item;
                }
              });

              dispatch({
                action: DELETE_PRESENT_REPO,
                payload: [...localState],
              });
            }

            setDeleteRepoStatus("success");
          } else {
            setDeleteRepoStatus("failed");
          }
        }
      })
      .catch((err) => {
        console.log(err);
        setDeleteRepoStatus("failed");
      });
  }

  const saveRenamedRepo = (renamedRepoObject) => {
    console.log(renamedRepoObject)
    // Perfom rename repo operation api call
  }

  const cancelRenameRepo = (e) => {
    setEditMode({status: false, idx: null})
  }

  const repoDetailsSettings = () => {
    return (
      <div className="repo-data my-10">
        <div className="text-xl text-gray-700 font-sans font-semibold">
          Saved Repos
        </div>
        <>
          {repoDetails && repoDetails.repoId && repoDetails.repoId.length ? (
            <>
              <div className="flex my-4 bg-indigo-500 w-full rounded text-white shadow p-3 font-sand text-xl font-semibold">
                <div className="w-1/2 border-r text-center">Repo ID</div>
                <div className="w-1/2 border-r text-center">Repo Name</div>
                <div className="w-1/2 border-r text-center">Repo Path</div>
                <div className="w-1/2 border-r text-center">Timestamp</div>
                <div className="w-1/2 border-r text-center">Action</div>
              </div>
              {repoDetails.repoId.map((repoId, idx) => {
                return (
                  <div
                    className="flex my-1 w-full rounded bg-white shadow p-3 font-sans text-gray-800"
                    key={repoId}
                  >
                    <div className="w-1/2 px-2 border-r font-sans break-all">
                      {repoId}
                    </div>
                    <div className="w-1/2 px-2 border-r font-bold font-sans break-all">
                      {editMode.status && editMode.idx === repoId ? ( <input
                                    id={repoId.concat('repoeditinput')}
                                    type="text"
                                    placeholder={repoDetails.repoName[idx]}
                                    className="repo-form--input mt-0 mb-0"
                                    onChange={(e) => setRenameRepo(e.target.value)}
                                    value={renameRepo}></input>) 
                      :  repoDetails.repoName[idx]}
                    </div>
                    <div className="w-1/2 px-2 border-r font-sans break-all text-sm font-light text-blue-600">
                      {repoDetails.repoPath[idx]}
                    </div>
                    <div className="w-1/2 px-2 border-r font-sans break-all text-sm font-light">
                      {repoDetails.timeStamp[idx]}
                    </div>
                    <div className="w-1/2 px-2 border-r font-sans break-all flex flex-row justify-center">
                      {editMode.idx === repoId && editMode.status ? (<><div
                        className="bg-indigo-500 p-2 mx-2 my-2 rounded shadow text-center w-1/2 hover:bg-indigo-400 cursor-pointer"
                        onClick={() => saveRenamedRepo({
                            repoId: repoId,
                            repoName: renameRepo,
                            repoPath: repoDetails.repoPath[idx],
                            timeStamp: repoDetails.timeStamp[idx]
                        })}>
                        <FontAwesomeIcon
                          color="white"
                          icon={["fas", "save"]}
                        ></FontAwesomeIcon>
                      </div> <div
                        className="bg-gray-600 p-2 mx-2 my-2 rounded shadow text-center w-1/2 hover:bg-gray-400 cursor-pointer"
                        onClick={cancelRenameRepo}>
                        <FontAwesomeIcon
                          color="white"
                          icon={["fas", "times"]}
                        ></FontAwesomeIcon>
                      </div></>) : (<div
                        className="bg-gray-600 p-2 mx-2 my-2 rounded shadow text-center w-1/2 hover:bg-gray-400 cursor-pointer"
                        onClick={(event) => {
                          setEditMode({
                            status: true,
                            idx: repoId
                          })
                        }}>
                        <FontAwesomeIcon
                          color="white"
                          icon={["fas", "edit"]}
                        ></FontAwesomeIcon>
                      </div>)}
                      <div
                        className="bg-red-600 p-2 mx-2 my-2 rounded shadow text-center w-1/2 hover:bg-red-400 cursor-pointer"
                        onClick={(event) => {
                          setBackdropToggle(true);
                          setDeleteRepo({
                            repoId: repoId,
                            repoName: repoDetails.repoName[idx],
                            repoPath: repoDetails.repoPath[idx],
                            timeStamp: repoDetails.timeStamp[idx],
                          });
                        }}
                      >
                        <FontAwesomeIcon
                          color="white"
                          icon={["fas", "trash-alt"]}
                        ></FontAwesomeIcon>
                      </div>
                    </div>
                  </div>
                );
              })}
            </>
          ) : (
            <div className="my-4 mx-auto bg-gray-200 text-center p-3 rounded shadow w-3/4">
              No repos are being managed by Gitconvex. You can add one from the
              dashboard
            </div>
          )}
        </>
      </div>
    );
  };

  function portDetailsSettings() {
    function portUpdateHandler() {
      if (port) {
        axios({
          url: globalAPIEndpoint,
          method: "POST",
          data: {
            query: `
              mutation {
                settingsEditPort(newPort: "${port}")
              }
            `,
          },
        })
          .then((res) => {
            const { settingsEditPort } = res.data.data;
            if (settingsEditPort === "PORT_UPDATED") {
              window.location.reload();
            } else {
              portUpdateFailed(true);
            }
          })
          .catch((err) => {
            console.log(err);
            setPortUpdateFailed(true);
          });
      }
    }

    return (
      <div className="my-2 mx-auto">
        <div className="text-xl font-sans text-gray-800 my-2">
          Active Gitconvex port
        </div>
        <div className="flex my-4">
          <input
            type="text"
            className="p-2 rounded border border-gray-500 bg-gray-200 text-gray-800 xl:w-1/2 lg:w-1/3 md:w-1/2 sm:w-1/2 w-1/2"
            value={port}
            onChange={(event) => {
              setPort(event.target.value);
            }}
          ></input>
          <div
            className="p-2 text-center mx-4 rounded border text-white bg-indigo-500 xl:w-1/6 lg:w-1/6 md:w-1/5 sm:w-1/4 w-1/4 hover:bg-indigo-600 cursor-pointer"
            onClick={() => {
              portUpdateHandler();
            }}
          >
            Update Port
          </div>
        </div>
        <div className="text-justify font-sand font-light text-sm my-4 text-gray-500 italic w-2/3">
          Make sure to restart the app and to change the port in the URL after
          updating it
        </div>
        {portUpdateFailed ? (
          <div className="my-2 p-2 rounded border border-red-300 text-red-700 font-sans font-semibold w-1/2 text-center">
            Port update failed
          </div>
        ) : null}
      </div>
    );
  }

  return (
    <>
      {backdropToggle ? (
        <div
          className="fixed w-full h-full top-0 left-0 right-0 flex xl:overflow-auto lg:overflow-auto md:overflow-none sm:overflow-none"
          id="settings-backdrop"
          style={{ background: "rgba(0,0,0,0.7)" }}
          onClick={(event) => {
            if (event.target.id === "settings-backdrop") {
              setDeleteRepoStatus("");
              setBackdropToggle(false);
              let localViewReload = viewReload + 1;
              setViewReload(localViewReload);
            }
          }}
        >
          {deleteRepo ? deleteRepoHandler() : null}
          <div
            className="top-0 right-0 fixed float-right font-semibold my-2 bg-red-500 text-3xl cursor-pointer text-center text-white align-middle rounded-full w-12 h-12 items-center shadow-md mr-5"
            onClick={() => {
              setDeleteRepoStatus("");
              setBackdropToggle(false);
              let localViewReload = viewReload + 1;
              setViewReload(localViewReload);
            }}
          >
            X
          </div>
        </div>
      ) : null}
      <div className="block w-full h-full pt-5 pb-10 overflow-auto">
        <div className="font-sans text-6xl my-4 mx-10 text-gray-700 block items-center align-middle">
          <FontAwesomeIcon
            className="text-5xl"
            icon={["fas", "cogs"]}
          ></FontAwesomeIcon>
          <span className="mx-10">Settings</span>
        </div>
        <div className="block my-10 justify-center mx-auto w-11/12">
          {dbPath ? databasePathSettings() : null}
          {repoDetails ? repoDetailsSettings() : null}
          {portDetailsSettings()}
        </div>
      </div>
    </>
  );
}
