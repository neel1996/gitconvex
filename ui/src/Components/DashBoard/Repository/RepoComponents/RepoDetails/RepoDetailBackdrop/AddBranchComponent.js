import React, { useState, useRef } from "react";
import axios from "axios";
import { globalAPIEndpoint } from "../../../../../../util/env_config";

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
    <div className="w-11/12 xl:w-1/2 lg:w-8/12 md:w-3/4 mx-auto my-auto bg-gray-50 rounded-lg">
      <div className="my-auto">
        <div
          className={`w-full mb-20 text-center p-4 mx-auto font-sans font-semibold text-2xl border-b-2 border-dashed ${
            branchAddStatus === "BRANCH_CREATION_SUCCESS"
              ? "text-green-400 border-green-400"
              : "text-gray-600 border-gray-400"
          }
        ${
          branchAddStatus === "BRANCH_ADD_FAILED"
            ? "text-red-400 border-red-400"
            : "text-gray-600 border-gray-400"
        }
        `}
        >
          {branchAddStatus === "" ? "ADD A NEW BRANCH TO THE REPO" : null}
          {branchAddStatus === "BRANCH_CREATION_SUCCESS"
            ? "NEW BRANCH ADDED SUCCESSFULLY"
            : null}
          {branchAddStatus === "BRANCH_ADD_FAILED"
            ? "NEW BRANCH ADDITION FAILED"
            : null}
        </div>
        <div className="mx-auto">
          <input
            type="text"
            ref={branchNameRef}
            id="branchName"
            placeholder="Enter branch name"
            className="w-11/12 flex justify-center mx-auto p-4 text-center rounded-lg shadow-md text-xl text-gray-800 border-2 border-blue-100 mb-20 outline-none"
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
          className="w-full mt-4 rounded-b-lg bg-indigo-500 p-4 text-center text-xl font-sans font-semibold text-white cursor-pointer hover:bg-indigo-400"
          onClick={(event) => {
            if (branchName) {
              addBranchClickHandler();
            } else {
              setBranchAddStatus("BRANCH_ADD_FAILED");
            }
          }}
        >
          ADD BRANCH
        </div>
      </div>
    </div>
  );
}
