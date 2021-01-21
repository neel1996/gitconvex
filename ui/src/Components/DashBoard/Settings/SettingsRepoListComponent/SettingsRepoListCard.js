import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useRef, useState } from "react";
import { globalAPIEndpoint } from "../../../../util/env_config";
import SettingsRepoDeleteBackDrop from "./SettingsRepoDeleteBackDrop";

export default function SettingsRepoListCard(props) {
  const { repoId, repoName, repoPath, timeStamp } = props;
  const repoNameRef = useRef();

  const [editMode, setEditMode] = useState(false);
  const [renameRepo, setRenameRepo] = useState("");
  const [repoNameState, setRepoNameState] = useState(repoName);
  const [backdropToggle, setBackdropToggle] = useState(false);
  const [deleteRepoData, setDeleteRepoData] = useState({ ...props });

  function saveRenamedRepo(repoId, updatedRepoName) {
    if (
      !updatedRepoName.match(/([a-zA-Z-_.])*([0-9])*/gi) ||
      updatedRepoName.match(/( )+/gi)
    ) {
      setEditMode(false);
      return;
    }
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          mutation {
            updateRepoName(repoId:"${repoId}", repoName:"${updatedRepoName}")
          }
        `,
      },
    })
      .then((res) => {
        if (res.data && !res.data.error) {
          setEditMode(false);
          setRepoNameState(updatedRepoName);
        }
      })
      .catch((err) => {
        console.log(err);
        setRenameRepo("ERROR!");
        repoNameRef.current.classList.add("border-2");
        repoNameRef.current.classList.add("border-red-400");
      });
  }

  return (
    <>
      {backdropToggle ? (
        <SettingsRepoDeleteBackDrop
          setBackdropToggle={setBackdropToggle}
          deleteRepoData={deleteRepoData}
        ></SettingsRepoDeleteBackDrop>
      ) : null}
      <div className="flex items-center my-1 w-full rounded bg-white shadow p-3 font-sans text-gray-800">
        {repoId ? (
          <>
            <div className="w-1/4 px-2 border-r font-sans break-all">
              {repoId}
            </div>
            <div className="w-1/2 px-2 border-r font-bold font-sans break-all">
              {editMode ? (
                <input
                  type="text"
                  placeholder={repoNameState}
                  className="w-full p-2 rounded shadow border border-blue-200 font-sans font-medium outline-white"
                  onChange={(e) => setRenameRepo(e.target.value)}
                  value={renameRepo}
                  ref={repoNameRef}
                ></input>
              ) : (
                repoNameState
              )}
            </div>
            <div className="w-1/2 px-2 border-r font-sans break-all text-sm font-light text-blue-600">
              {repoPath}
            </div>
            <div className="w-1/2 px-2 border-r font-sans break-all text-sm font-light">
              {timeStamp}
            </div>
            <div className="w-1/2 px-2 border-r font-sans break-all flex flex-row justify-center">
              {editMode ? (
                <>
                  <div
                    className="bg-indigo-500 p-2 mx-2 my-2 rounded shadow text-center w-1/2 hover:bg-indigo-400 cursor-pointer"
                    onClick={() => saveRenamedRepo(repoId, renameRepo)}
                  >
                    <FontAwesomeIcon
                      color="white"
                      icon={["fas", "save"]}
                    ></FontAwesomeIcon>
                  </div>
                  <div
                    className="bg-gray-600 p-2 mx-2 my-2 rounded shadow text-center w-1/2 hover:bg-gray-400 cursor-pointer"
                    onClick={() => {
                      setEditMode(false);
                    }}
                  >
                    <FontAwesomeIcon
                      color="white"
                      icon={["fas", "times"]}
                    ></FontAwesomeIcon>
                  </div>
                </>
              ) : (
                <div
                  className="bg-gray-600 p-2 mx-2 my-2 rounded shadow text-center w-1/2 hover:bg-gray-400 cursor-pointer"
                  onClick={(event) => {
                    setEditMode(true);
                  }}
                >
                  <FontAwesomeIcon
                    color="white"
                    icon={["fas", "edit"]}
                  ></FontAwesomeIcon>
                </div>
              )}
              <div
                className="bg-red-600 p-2 mx-2 my-2 rounded shadow text-center w-1/2 hover:bg-red-400 cursor-pointer"
                onClick={(event) => {
                  setBackdropToggle(true);
                  setDeleteRepoData({ ...props });
                }}
              >
                <FontAwesomeIcon
                  color="white"
                  icon={["fas", "trash-alt"]}
                ></FontAwesomeIcon>
              </div>
            </div>
          </>
        ) : null}
      </div>
    </>
  );
}
