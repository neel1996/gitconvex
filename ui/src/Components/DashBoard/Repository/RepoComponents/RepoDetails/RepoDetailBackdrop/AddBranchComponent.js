import React, { useState, useRef } from "react";
import axios from "axios";
import { globalAPIEndpoint } from "../../../../../../util/env_config";
import "../../../../../styles/RepositoryDetailsBackdrop.css";

export default function AddBranchComponent(props) {
  const { repoId } = props;
  const [branchName, setBranchName] = useState("");
  const [branchAddStatus, setBranchAddStatus] = useState("");
  const branchNameRef = useRef();

  function resetBranchNameText() {
    branchNameRef.current.value = "";
    setBranchName("");
  }

  function addBranchClickHandler() {
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            mutation GitConvexMutation{
              addBranch(repoId: "${repoId}", branchName: "${branchName}")
            }
          `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          const branchStatus = res.data.data.addBranch;
          setBranchAddStatus(branchStatus);
          resetBranchNameText();
        } else {
          setBranchAddStatus("BRANCH_ADD_FAILED");
          resetBranchNameText();
        }
      })
      .catch((err) => {
        setBranchAddStatus("BRANCH_ADD_FAILED");
        resetBranchNameText();
      });
  }

  return (
    <div className="repo-backdrop--addbranch">
      <div className="my-auto">
        <div className="mx-auto">
          <input
            type="text"
            ref={branchNameRef}
            id="branchName"
            placeholder="Branch Name"
            className="addbranch--input"
            onChange={(event) => {
              const branchNameVal = event.target.value;
              if (
                event.target.id === "branchName" &&
                branchNameVal.match(/[^a-zA-Z0-9_.:^\\/]/gi)
              ) {
                event.target.value = branchNameVal.replace(
                  /[^a-zA-Z0-9_.:^\\/]/gi,
                  "-"
                );
              }
              setBranchName(event.target.value);
            }}
            onClick={() => {
              setBranchAddStatus("");
            }}
          ></input>
        </div>
        <div
          className="addbranch--btn"
          onClick={(event) => {
            if (branchName) {
              addBranchClickHandler();
            } else {
              setBranchAddStatus("BRANCH_ADD_FAILED");
            }
          }}
        >
          Add Branch
        </div>
        {branchAddStatus === "BRANCH_CREATION_SUCCESS" ? (
          <div className="backdrop--alert bg-green-200 text-green border-green-500">
            New branch has been added to your repo successfully
          </div>
        ) : null}
        {branchAddStatus === "BRANCH_ADD_FAILED" ? (
          <div className="backdrop--alert bg-red-200 text-red-600 border-red-400">
            New branch addition failed!
          </div>
        ) : null}
      </div>
    </div>
  );
}
