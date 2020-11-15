import React, { useState, useRef } from "react";
import axios from "axios";
import { globalAPIEndpoint } from "../../../../../../util/env_config";

export default function AddRemoteRepoComponent({ repoId }) {
  const remoteNameRef = useRef();
  const remoteUrlRef = useRef();

  const [paramMissing, setParamMissing] = useState(false);
  const [addRemoteStatus, setAddRemoteStatus] = useState("");

  const remoteFormTextComponent = (formId, label, placeholder) => {
    return (
      <div className="addremote--form">
        <label htmlFor={formId} className="addremote--form--label">
          {label}
        </label>
        <div className="w-5/6">
          <input
            id={formId}
            onClick={() => {
              setParamMissing(false);
              setAddRemoteStatus("");
            }}
            className="backdrop--input"
            placeholder={placeholder}
            ref={formId === "remoteName" ? remoteNameRef : remoteUrlRef}
            onChange={(event) => {
              const remoteNameVal = event.target.value;
              if (
                event.target.id === "remoteName" &&
                remoteNameVal.match(/[^a-zA-Z0-9_]/gi)
              ) {
                event.target.value = remoteNameVal.replace(
                  /[^a-zA-Z0-9_]/gi,
                  "-"
                );
              }
            }}
          ></input>
        </div>
      </div>
    );
  };

  function addRemoteClickHandler() {
    let repoName = remoteNameRef.current.value;
    let repoUrl = remoteUrlRef.current.value;

    if (repoId && repoName && repoUrl) {
      axios({
        url: globalAPIEndpoint,
        method: "POST",
        data: {
          query: `
              mutation {
                addRemote(repoId: "${repoId}", remoteName: "${repoName}", remoteUrl: "${repoUrl}")
              }
            `,
        },
      })
        .then((res) => {
          if (res.data.data && !res.data.error) {
            const remoteAddStatus = res.data.data.addRemote;

            if (remoteAddStatus === "REMOTE_ADD_SUCCESS") {
              setAddRemoteStatus("success");
              remoteNameRef.current.value = "";
              remoteUrlRef.current.value = "";
            } else {
              setAddRemoteStatus("failed");
            }
          } else {
            setAddRemoteStatus("failed");
          }
        })
        .catch((err) => {
          console.log(err);
          setAddRemoteStatus("failed");
        });
    } else {
      setParamMissing(true);
    }
  }

  const statusPillComponent = (color, message) => {
    return (
      <div className={`addremote--alert border-${color}-900 bg-${color}-200`}>
        {message}
      </div>
    );
  };

  return (
    <div className="xl:w-1/2 lg:w-3/4 md:w-11/12 sm:w-11/12 repo-backdrop--addremote">
      <div className="addremote--header">Enter new remote details</div>
      <div className="my-4 mx-6">
        {remoteFormTextComponent(
          "remoteName",
          "Enter Remote Name",
          "Give a name for your new remote"
        )}
        {remoteFormTextComponent(
          "remoteUrl",
          "Enter Remote URL",
          "Provide the URL for your remote repo"
        )}
      </div>
      {paramMissing
        ? statusPillComponent(
            "orange",
            "One or more required parameters are empty!"
          )
        : null}
      {addRemoteStatus === "success"
        ? statusPillComponent(
            "green",
            "Remote repo has been added successfully!"
          )
        : null}
      {addRemoteStatus === "failed"
        ? statusPillComponent("red", "Failed to add new repo!")
        : null}
      <div
        className="addremote--btn"
        onClick={() => {
          addRemoteClickHandler();
        }}
      >
        Add New Remote
      </div>
    </div>
  );
}
