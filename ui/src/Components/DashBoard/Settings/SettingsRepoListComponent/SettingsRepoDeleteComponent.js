import axios from "axios";
import React, { useContext, useState } from "react";
import { globalAPIEndpoint } from "../../../../util/env_config";
import { SettingsContext } from "../Settings";

export default function SettingsRepoDeleteComponent(props) {
  const { deleteRepoData } = props;
  const { state, dispatch } = useContext(SettingsContext);

  const repoColumn = ["Repo ID", "Repo Name", "Repo Path", "Timestamp"];
  let repoArray = [];

  const [deleteRepoStatus, setDeleteRepoStatus] = useState("");

  Object.keys(deleteRepoData).forEach((key, index) => {
    repoArray.push({ label: repoColumn[index], value: deleteRepoData[key] });
  });

  function deleteRepoApiHandler(repoId) {
    setDeleteRepoStatus("loading");
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
        if (res.data.data && !res.data.error) {
          const { status } = res.data.data.deleteRepo;
          if (status === "DELETE_SUCCESS") {
            setDeleteRepoStatus("success");
            dispatch(state.viewReloadCount + 1);
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
            deleteRepoApiHandler(deleteRepoData.repoId);
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
