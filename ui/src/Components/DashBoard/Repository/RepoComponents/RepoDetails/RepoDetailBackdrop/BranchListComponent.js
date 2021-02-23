import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../../../util/env_config";

export default function BranchListComponent({ repoId }) {
  library.add(fas);

  const [branchList, setBranchList] = useState([]);
  const [listError, setListError] = useState(false);
  const [switchSuccess, setSwitchSuccess] = useState(false);
  const [switchError, setSwitchError] = useState(false);
  const [switchedBranch, setSwitchedBranch] = useState("");
  const [errorBranch, setErrorBranch] = useState("");
  const [deleteSuccess, setDeleteSuccess] = useState(false);
  const [deleteError, setDeleteError] = useState(false);
  const [loading, setLoading] = useState(false);
  const [branchSearchTerm, setBranchSearchTerm] = useState("");
  const [filteredBranchList, setFilteredBranchList] = useState([]);

  function resetStates() {
    setListError(false);
    setSwitchSuccess(false);
    setSwitchError(false);
    setSwitchedBranch("");
    setErrorBranch("");
    setDeleteError(false);
    setDeleteSuccess(false);
    setBranchSearchTerm("");
    setFilteredBranchList([]);
  }

  useEffect(() => {
    const token = axios.CancelToken;
    const source = token.source();

    setLoading(true);
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
        setLoading(false);

        if (res.data.data && !res.data.error) {
          let {
            gitAllBranchList,
            gitCurrentBranch,
          } = res.data.data.gitRepoStatus;

          if (gitCurrentBranch === "Repo HEAD is nil") {
            setBranchList([]);
            setListError(true);
            return;
          }
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
        setLoading(false);

        if (err) {
          console.log("API error occurred : " + err);
          setListError(true);
        }
      });

    return () => source.cancel;
  }, [repoId, switchedBranch]);

  function switchBranchHandler(branchName) {
    resetStates();
    setLoading(true);
    setBranchList([]);
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
        setLoading(false);

        if (res.data.data && !res.data.error) {
          const checkoutStatus = res.data.data.checkoutBranch;
          if (checkoutStatus === "CHECKOUT_FAILED") {
            setSwitchSuccess(false);
            setErrorBranch(branchName);
            setSwitchError(true);
            return;
          } else {
            setSwitchSuccess(true);
            setSwitchedBranch(branchName);
          }
        } else {
          setSwitchError(true);
          setErrorBranch(branchName);
        }
      })
      .catch((err) => {
        setLoading(false);

        if (err) {
          setSwitchError(true);
          setErrorBranch(branchName);
        }
      });
  }

  function deleteBranchHandler(branchName, forceFlag) {
    resetStates();
    setLoading(true);

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
        setLoading(false);

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
        setLoading(false);

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

  const searchBranchFromList = (event) => {
    const searchBranch = event.target.value;
    setBranchSearchTerm(searchBranch);
    if (searchBranch !== "") {
      const filteredBranches = branchList.filter((branchName) =>
        branchName.toLowerCase().includes(searchBranch)
      );
      setFilteredBranchList(filteredBranches);
    } else {
      setFilteredBranchList([]);
    }
  };

  const cancelSearchBranchFromList = () => {
    setBranchSearchTerm("");
    setFilteredBranchList([]);
  };

  const renderBranchListComponent = (branch) => {
    const branchPickerComponent = (icon, branchType, branchName) => {
      let activeSwitchStyle = "";
      let activeBranchFlag = false;
      if (branchName.includes("*")) {
        activeBranchFlag = true;
        branchName = branchName.replace("*", "");
      }

      if (activeBranchFlag) {
        activeSwitchStyle = "border-dashed border-b-2 text-indigo-700 text-2xl";
      }
      return (
        <div className="flex items-center justify-center px-14 py-4 mx-auto border-dashed border-b" key={branchType + branchName}>
          <div
            className={
              icon === "wifi"
                ? "mx-2 text-2xl text-blue-500 ml-0"
                : "mx-2 text-2xl text-blue-500"
            }
          >
            <FontAwesomeIcon icon={["fas", icon]}></FontAwesomeIcon>
          </div>
          <div className="xl:block lg:block md:block sm:hidden w-1/3 font-sans text-lg font-semibold text-center text-indigo-500">
            {branchType}
          </div>
          <div
            className={`w-1/2 font-sans font-semibold text-lg text-left cursor-pointer overflow-hidden text-gray-500 hover:text-blue-800 ${activeSwitchStyle}`}
            title={branchName}
            onClick={() => {
              if (!activeBranchFlag) {
                if (branchType !== "Local Branch") {
                  switchBranchHandler(branch);
                } else {
                  switchBranchHandler(branchName);
                }
              }
            }}
          >
            {branchName}
          </div>
          {!activeBranchFlag && branchType === "Local Branch" ? (
            <div className="flex mx-4 justify-between my-auto text-center items-center w-1/4 px-2">
              <div
                className="w-1/2 block mx-auto my-auto text-center px-2 justify-center bg-red-500 p-2 rounded-lg shadow-md cursor-pointer text-white font-sans font-semibold hover:bg-red-600"
                title="Will delete the branch forcefully.Be careful!"
                onClick={() => {
                  if (!activeBranchFlag) {
                    deleteBranchHandler(branchName, true);
                  }
                }}
              >
                <div>
                  <FontAwesomeIcon
                    icon={["fas", "trash-alt"]}
                  ></FontAwesomeIcon>
                </div>
                <div>DELETE</div>
              </div>
            </div>
          ) : (
            <>
              {activeBranchFlag ? (
                <div className="w-1/4 font-sans mx-4 text-sm px-2 font-light border border-dashed p-1 rounded-full text-center bg-blue-100 border-blue-800">
                  Active
                </div>
              ) : (
                <div className="w-1/4 font-sans mx-4 text-sm px-2 font-light border border-dashed p-1 rounded-full text-center bg-yellow-100 border-yellow-700">
                  Remote
                </div>
              )}
            </>
          )}
        </div>
      );
    };

    if (!branch.includes("remotes/")) {
      return branchPickerComponent("code-branch", "Local Branch", branch);
    } else {
      const splitBranch = branch.split("/");
      if (splitBranch.length <= 2) {
        return null;
      }
      const remoteName = splitBranch[1];
      const remoteBranch = splitBranch.slice(2, splitBranch.length).join("/");

      return branchPickerComponent("wifi", remoteName, remoteBranch);
    }
  };

  return (
    <div className="bg-gray-50 p-6 mx-auto my-auto items-center rounded-lg w-11/12 xl:w-3/4 lg:w-3/4 md:w-11/12 sm:w-11/12">
      <div className="text-4xl my-4 font-sans font-semibold text-gray-600">Available Branches</div>
      <div className="flex justify-start items-center font-sans text-sm font-semibold text-red-500">
        <div>
          <FontAwesomeIcon
            icon={["fas", "exclamation-circle"]}
          ></FontAwesomeIcon>
        </div>
        <div className="mx-2 font-sans font-semibold">
          Note that this section also lets you delete the branches, so be
          cautious!
        </div>
      </div>
      <div className="italic font-sans font-semibold text-lg my-2 border-b-2 border-dashed border-gray-300 text-gray-300">
        Click on a branch to checkout to that branch
      </div>
      <div className="w-full mx-auto my-6 overflow-y-auto overflow-x-hidden" style={{ height: "400px" }}>
        {loading ? (
          <div className="text-center font-sans font-light text-xl my-2 text-gray-600 border-b border-dotted">
            Collecting branch list...
          </div>
        ) : null}
        {!loading ? (
          <div className="flex flex-row mx-8 shadow-md rounded-md my-4">
            <div className="b-1 text-center p-4 text-white bg-blue-500 rounded-l-md">
              <FontAwesomeIcon icon={["fas", "search"]}></FontAwesomeIcon>
            </div>
            <div className="w-full">
              <input
                id="branchListSearchInput"
                type="text"
                placeholder="Search For Branch Name"
                className="border-0 outline-none w-full p-4 focus:outline-none"
                onChange={searchBranchFromList}
                value={branchSearchTerm}
              ></input>
            </div>
            <div
              className="b-1 text-center p-4 text-gray-500 cursor-pointer bg-white rounded-r-md"
              onClick={cancelSearchBranchFromList}
            >
              <FontAwesomeIcon icon={["fas", "times"]}></FontAwesomeIcon>
            </div>
          </div>
        ) : null}
        {!listError && branchList && branchSearchTerm !== "" ? (
          filteredBranchList.length > 0 ? (
            filteredBranchList.map((branch) => {
              return renderBranchListComponent(branch);
            })
          ) : (
            <div className="text-center font-sans font-light text-xl my-2 text-gray-600 border-b border-dotted">
              <span className="mx-2 font-semibold border-b border-dashed">
                {branchSearchTerm}
              </span>
              branch is not available!
            </div>
          )
        ) : (
          branchList.map((branch) => {
            return renderBranchListComponent(branch);
          })
        )}
      </div>
      {listError
        ? errorComponent("Error occurred while listing branches!")
        : null}

      {switchError
        ? errorComponent("Error occurred while switching to", true)
        : null}

      {switchedBranch.length > 0 && switchSuccess
        ? successComponent("Active branch has been switched to")
        : null}

      {deleteSuccess
        ? successComponent("Selected branch has been removed")
        : null}

      {deleteError ? errorComponent("Branch deletion failed for", true) : null}
    </div>
  );
}
