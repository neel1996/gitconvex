import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../../../util/env_config";

export default function BranchListComponent({ repoId, currentBranch }) {
  library.add(fas);

  const [branchList, setBranchList] = useState([]);
  const [listError, setListError] = useState(false);
  const [switchSuccess, setSwitchSuccess] = useState(false);
  const [switchError, setSwitchError] = useState(false);
  const [switchedBranch, setSwitchedBranch] = useState("");
  const [errorBranch, setErrorBranch] = useState("");
  const [deleteSuccess, setDeleteSuccess] = useState(false);
  const [deleteError, setDeleteError] = useState(false);

  function resetStates() {
    setListError(false);
    setSwitchSuccess(false);
    setSwitchError(false);
    setSwitchedBranch("");
    setErrorBranch("");
    setDeleteError(false);
    setDeleteSuccess(false);
  }

  useEffect(() => {
    const token = axios.CancelToken;
    const source = token.source();

    setBranchList([]);

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      cancelToken: source.token,
      data: {
        query: `
          query
          {
            gitRepoStatus(repoId:"${repoId}"){
                gitAllBranchList  
                gitCurrentBranch
            }
          }
        `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          let {
            gitAllBranchList,
            gitCurrentBranch,
          } = res.data.data.gitRepoStatus;

          gitAllBranchList = gitAllBranchList.map((branch) => {
            if (branch === gitCurrentBranch) {
              return "*" + branch;
            }
            return branch;
          });

          setBranchList([...gitAllBranchList]);
        } else {
          setListError(true);
        }
      })
      .catch((err) => {
        if (err) {
          console.log("API error occurred : " + err);
          setListError(true);
        }
      });

    return () => source.cancel;
  }, [repoId, switchedBranch]);

  function switchBranchHandler(branchName) {
    resetStates();

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            mutation{
              checkoutBranch(repoId: "${repoId}", branchName: "${branchName}")
            }
          `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          setSwitchSuccess(true);
          setSwitchedBranch(branchName);
        } else {
          setSwitchError(true);
        }
      })
      .catch((err) => {
        if (err) {
          setSwitchError(true);
          setErrorBranch(branchName);
        }
      });
  }

  function deleteBranchHandler(branchName, forceFlag) {
    resetStates();

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            mutation{
              deleteBranch(repoId: "${repoId}", branchName: "${branchName}", forceFlag: ${forceFlag}){
                status
              }
            }
          `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          if (res.data.data.deleteBranch.status === "BRANCH_DELETE_SUCCESS") {
            setDeleteSuccess(true);
            setSwitchedBranch(branchName);
          } else {
            setDeleteError(true);
            setErrorBranch(branchName);
          }
        }
      })
      .catch((err) => {
        if (err) {
          setDeleteError(true);
          setErrorBranch(branchName);
        }
      });
  }

  function errorComponent(errorString, branchError = false) {
    return (
      <div className="text-center p-4 rounded bg-red-300 text-xl font-sans mt-10">
        {errorString}
        {branchError ? (
          <span className="font-semibold border-b border-dashed mx-2">
            {errorBranch}
          </span>
        ) : null}
      </div>
    );
  }

  function successComponent(successString) {
    return (
      <div className="text-center p-4 rounded bg-green-300 text-xl font-sans mt-10">
        {successString}
        <span className="font-semibold border-b border-dashed mx-2">
          {switchedBranch}
        </span>
      </div>
    );
  }

  return (
    <div className="repo-backdrop--branchlist xl:w-3/4 lg:w-3/4 md:w-11/12 sm:w-11/12">
      <div className="branchlist--header">Available Branches</div>
      <div className="branchlist--warn">
        <div className="mx-2">
          <FontAwesomeIcon
            icon={["fas", "exclamation-circle"]}
          ></FontAwesomeIcon>
        </div>
        <div className="font-sans font-semibold">
          Note that this section also lets you delete the branches, so be
          cautious!
        </div>
      </div>
      <div className="branchlist--infotext">
        Click on a branch to checkout to that branch
      </div>
      <div className="branchlist--list-area" style={{ height: "400px" }}>
        {branchList.length === 0 ? (
          <div className="list-area--message">Collecting branch list...</div>
        ) : null}
        {!listError &&
          branchList &&
          branchList.map((branch) => {
            const branchPickerComponent = (icon, branchType, branchName) => {
              let activeSwitchStyle = "";
              let activeBranchFlag = false;
              if (branchName.includes("*")) {
                activeBranchFlag = true;
                branchName = branchName.replace("*", "");
              }

              if (activeBranchFlag) {
                activeSwitchStyle =
                  "border-dashed border-b-2 text-indigo-700 text-2xl";
              }
              return (
                <div
                  className="list-area--branches"
                  key={branchType + branchName}
                >
                  <div className="list-area--branches--icon">
                    <FontAwesomeIcon icon={["fas", icon]}></FontAwesomeIcon>
                  </div>
                  <div className="xl:block lg:block md:block sm:hidden list-area--branches--type">
                    {branchType}
                  </div>
                  <div
                    className={`list-area--branches--name ${activeSwitchStyle}`}
                    title={branchName}
                    onClick={() => {
                      if (!activeBranchFlag) {
                        switchBranchHandler(branchName);
                      }
                    }}
                  >
                    {branchName}
                  </div>
                  {!activeBranchFlag && branchType === "Local Branch" ? (
                    <div className="list-area--branches--active">
                      <div
                        className="list-area--branches--delete"
                        title="Will delete only if the branch is clean and safe"
                        onClick={() => {
                          if (!activeBranchFlag) {
                            deleteBranchHandler(branchName, false);
                          }
                        }}
                      >
                        <div className="list-area--branches--delete--btn">
                          <FontAwesomeIcon
                            icon={["fas", "trash-alt"]}
                          ></FontAwesomeIcon>
                        </div>
                        <div className="list-area--branches--delete--type">
                          Normal
                        </div>
                      </div>
                      <div
                        className="list-area--branches--delete"
                        title="Will delete the branch forcefully.Be careful!"
                        onClick={() => {
                          if (!activeBranchFlag) {
                            deleteBranchHandler(branchName, true);
                          }
                        }}
                      >
                        <div className="list-area--branches--delete--btn">
                          <FontAwesomeIcon
                            icon={["fas", "minus-square"]}
                          ></FontAwesomeIcon>
                        </div>
                        <div className="list-area--branches--delete--type">
                          Force
                        </div>
                      </div>
                    </div>
                  ) : (
                    <>
                      {activeBranchFlag ? (
                        <div className="list-area--branches--pill bg-blue-200 border-blue-800">
                          Active
                        </div>
                      ) : (
                        <div className="list-area--branches--pill bg-orange-200 border-orange-800">
                          Remote
                        </div>
                      )}
                    </>
                  )}
                </div>
              );
            };

            if (!branch.includes("remotes/")) {
              return branchPickerComponent(
                "code-branch",
                "Local Branch",
                branch
              );
            } else {
              const splitBranch = branch.split("/");
              if (splitBranch.length <= 2) {
                return null;
              }
              const remoteName = splitBranch[1];
              const remoteBranch = splitBranch
                .slice(2, splitBranch.length)
                .join("/");

              return branchPickerComponent("wifi", remoteName, remoteBranch);
            }
          })}
      </div>
      {listError
        ? errorComponent("Error occurred while listing branches!")
        : null}

      {switchError
        ? errorComponent("Error occurred while switching to", true)
        : null}

      {switchedBranch && switchSuccess
        ? successComponent("Active branch has been switched to")
        : null}

      {deleteSuccess
        ? successComponent("Selected branch has been removed")
        : null}

      {deleteError ? errorComponent("Branch deletion failed for", true) : null}
    </div>
  );
}
