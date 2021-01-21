import { library } from "@fortawesome/fontawesome-svg-core";
import { fab } from "@fortawesome/free-brands-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useMemo, useState } from "react";
import { globalAPIEndpoint } from "../../../../../util/env_config";
import LoadingHOC from "../../../../LoadingHOC";
import FileExplorerComponent from "./FileExplorerComponent";
import AddBranchComponent from "./RepoDetailBackdrop/AddBranchComponent";
import AddRemoteRepoComponent from "./RepoDetailBackdrop/AddRemoteRepoComponent";
import BranchListComponent from "./RepoDetailBackdrop/BranchListComponent";
import CommitLogComponent from "./RepoDetailBackdrop/CommitLogComponent";
import FetchPullActionComponent from "./RepoDetailBackdrop/FetchPullActionComponent";
import SwitchBranchComponent from "./RepoDetailBackdrop/SwitchBranchComponent";
import RepoInfoComponent from "./RepoInfoComponent";
import RepoLeftPaneComponent from "./RepoLeftPaneComponent";
import RepoRightPaneComponent from "./RepoRightPaneComponent";

export default function RepositoryDetails(props) {
  library.add(fab, fas);
  const [gitRepoStatus, setGitRepoStatus] = useState({});
  const [repoFetchFailed, setRepoFetchFailed] = useState(false);
  const [repoIdState, setRepoIdState] = useState("");
  const [showCommitLogs, setShowCommitLogs] = useState(false);
  const [isMultiRemote, setIsMultiRemote] = useState(false);
  const [multiRemoteCount, setMultiRemoteCount] = useState(0);
  const [backdropToggle, setBackdropToggle] = useState(false);
  const [reloadView, setReloadView] = useState(false);
  const [codeViewToggle, setCodeViewToggle] = useState(false);
  const [selectedBranch, setSelectedBranch] = useState("");
  const [currentBranch, setCurrentBranch] = useState("");
  const [action, setAction] = useState("");
  const [loading, setLoading] = useState(false);

  const closeBackdrop = (toggle) => {
    setBackdropToggle(!toggle);
  };

  const showCommitLogsView = () => {
    setShowCommitLogs(true);
  };

  const actionTrigger = (actionType) => {
    setAction(actionType);
    setBackdropToggle(true);
  };

  const memoizedFolderExplorer = useMemo(() => {
    return (
      <FileExplorerComponent repoIdState={repoIdState}></FileExplorerComponent>
    );
  }, [repoIdState]);

  const memoizedCommitLogComponent = useMemo(() => {
    return (
      <>
        <CommitLogComponent repoId={repoIdState}></CommitLogComponent>
      </>
    );
  }, [repoIdState]);

  const memoizedFetchRemoteComponent = useMemo(() => {
    return (
      <FetchPullActionComponent
        repoId={repoIdState}
        actionType="fetch"
      ></FetchPullActionComponent>
    );
  }, [repoIdState]);

  const memoizedPullRemoteComponent = useMemo(() => {
    return (
      <FetchPullActionComponent
        repoId={repoIdState}
        actionType="pull"
      ></FetchPullActionComponent>
    );
  }, [repoIdState]);

  const memoizedSwitchBranchComponent = useMemo(() => {
    return (
      <SwitchBranchComponent
        repoId={repoIdState}
        branchName={selectedBranch}
        closeBackdrop={closeBackdrop}
        switchReloadView={() => {
          setReloadView(true);
        }}
      ></SwitchBranchComponent>
    );
  }, [repoIdState, selectedBranch]);

  const memoizedBranchListComponent = useMemo(() => {
    return (
      <BranchListComponent
        repoId={repoIdState}
        currentBranch={currentBranch}
      ></BranchListComponent>
    );
  }, [repoIdState, currentBranch]);

  const memoizedAddRemoteRepoComponent = useMemo(() => {
    return (
      <AddRemoteRepoComponent repoId={repoIdState}></AddRemoteRepoComponent>
    );
  }, [repoIdState]);

  useEffect(() => {
    setReloadView(false);
    setCodeViewToggle(false);
    setLoading(true);
    const endpointURL = globalAPIEndpoint;

    if (props.parentProps.location) {
      const repoId = props.parentProps.location.pathname.split(
        "/repository/"
      )[1];

      setRepoIdState(repoId);

      axios({
        url: endpointURL,
        method: "POST",
        headers: {
          "Content-type": "application/json",
        },
        data: {
          query: `

            query
            {
                gitRepoStatus(repoId:"${repoId}"){
                  gitRemoteData
                  gitRepoName
                  gitBranchList
                  gitCurrentBranch
                  gitRemoteHost
                  gitTotalCommits
                  gitLatestCommit
                  gitTotalTrackedFiles    
                }
            }
          `,
        },
      })
        .then((res) => {
          setLoading(false);

          if (res.data && res.data.data && !res.data.error) {
            const localRepoStatus = res.data.data.gitRepoStatus;
            let gitRemoteLocal = localRepoStatus.gitRemoteData;
            setCurrentBranch(localRepoStatus.gitCurrentBranch);
            if (gitRemoteLocal.includes("||")) {
              setIsMultiRemote(true);
              localRepoStatus.gitRemoteData = gitRemoteLocal.split("||")[0];
              setIsMultiRemote(true);
              setMultiRemoteCount(gitRemoteLocal.split("||").length);
            }
            setGitRepoStatus(localRepoStatus);
          } else {
            setRepoFetchFailed(true);
          }
        })
        .catch((err) => {
          setLoading(false);

          if (err) {
            console.log("API GitStatus error occurred : " + err);
            setRepoFetchFailed(true);
          }
        });
    }
  }, [props.parentProps, reloadView]);

  let {
    gitRemoteData,
    gitRepoName,
    gitBranchList,
    gitCurrentBranch,
    gitRemoteHost,
    gitTotalCommits,
    gitLatestCommit,
  } = gitRepoStatus;

  const switchBranchHandler = (branchName) => {
    setBackdropToggle(true);
    setAction("switchbranch");
    setSelectedBranch(branchName);
  };

  const actionComponentPicker = () => {
    switch (action) {
      case "fetch":
        return memoizedFetchRemoteComponent;
      case "pull":
        return memoizedPullRemoteComponent;
      case "addRemoteRepo":
        return memoizedAddRemoteRepoComponent;
      case "addBranch":
        return <AddBranchComponent repoId={repoIdState}></AddBranchComponent>;
      case "switchbranch":
        return memoizedSwitchBranchComponent;
      case "listBranch":
        return memoizedBranchListComponent;
      default:
        return null;
    }
  };

  return (
    <>
      {loading ? (
        <LoadingHOC
          message="Fetching repo details..."
          loading={loading}
        ></LoadingHOC>
      ) : null}
      {showCommitLogs ? (
        <>
          <div
            className="fixed w-full h-full top-0 left-0 right-0 flex overflow-auto"
            id="commit-log__backdrop"
            style={{ background: "rgba(0,0,0,0.5)", zIndex: 99 }}
            onClick={(event) => {
              if (event.target.id === "commit-log__backdrop") {
                setShowCommitLogs(false);
              }
            }}
          >
            <div
              id="commit-log__cards"
              className="w-full xl:w-3/4 lg:w-5/6 md:w-11/12 sm:w-11/12 h-full block mx-auto my-auto mt-10 mb-10"
            >
              {memoizedCommitLogComponent}
            </div>
            <div
              className="w-14 h-14 mr-5 mt-6 rounded-full bg-red-500 text-white flex justify-center items-center shadow cursor-pointer fixed right-0 top-0"
              onClick={() => {
                setShowCommitLogs(false);
              }}
            >
              <FontAwesomeIcon
                className="flex text-center text-3xl my-auto"
                icon={["fas", "times"]}
              ></FontAwesomeIcon>
            </div>
          </div>
        </>
      ) : null}
      {backdropToggle || codeViewToggle ? (
        <div
          className="flex h-full overflow-auto fixed inset-x-0 top-0 w-full z-40"
          id="repo-backdrop"
          style={{ background: "rgba(0,0,0,0.7)", zIndex: "99" }}
          onClick={(event) => {
            if (event.target.id === "repo-backdrop") {
              setBackdropToggle(false);
              setAction("");
            }
          }}
        >
          <>{action ? actionComponentPicker() : null}</>
          <div
            className="w-14 h-14 mr-5 mt-6 rounded-full bg-red-500 text-white flex justify-center items-center shadow cursor-pointer fixed top-0 right-0"
            onClick={() => {
              setBackdropToggle(false);
              setCodeViewToggle(false);
              setReloadView(true);
              setAction("");
            }}
          >
            <FontAwesomeIcon
              className="flex text-center text-3xl my-auto"
              icon={["fas", "times"]}
            ></FontAwesomeIcon>
          </div>
        </div>
      ) : null}
      <>
        {!loading && gitRepoStatus && !repoFetchFailed ? (
          <div className="overflow-auto rounded-lg justify-evenly h-full mx-auto p-6 w-full">
            <div className="flex px-3 py-2">
              {gitRepoStatus ? (
                <RepoInfoComponent
                  gitRepoName={gitRepoName}
                  gitCurrentBranch={gitCurrentBranch}
                ></RepoInfoComponent>
              ) : null}
            </div>
            <div className="w-full">
              <div className="xl:w-11/12 lg:w-full w-full xl:flex lg:block md:block sm:block my-4 mx-auto justify-around">
                {gitRepoStatus ? (
                  <>
                    <RepoLeftPaneComponent
                      received={true}
                      actionTrigger={actionTrigger}
                      showCommitLogsView={showCommitLogsView}
                      gitRemoteHost={gitRemoteHost}
                      gitRemoteData={gitRemoteData}
                      isMultiRemote={isMultiRemote}
                      multiRemoteCount={multiRemoteCount}
                    ></RepoLeftPaneComponent>
                    <RepoRightPaneComponent
                      received={true}
                      switchBranchHandler={switchBranchHandler}
                      actionTrigger={actionTrigger}
                      gitBranchList={gitBranchList}
                      gitCurrentBranch={gitCurrentBranch}
                      gitLatestCommit={gitLatestCommit}
                      gitTotalCommits={gitTotalCommits}
                    ></RepoRightPaneComponent>
                  </>
                ) : null}
              </div>
            </div>

            {!loading && gitRepoStatus && repoIdState
              ? memoizedFolderExplorer
              : null}
          </div>
        ) : !loading ? (
          <div className="w-full h-full mx-auto text-center flex justify-center items-center">
            <div className="block mx-auto w-11/12">
              <div className="rounded-lg shadow border-2 border-dashed border-pink-400 text-red-300 text-3xl font-sans font-semibold p-4">
                Unable to fetch repo details
              </div>
              <div className="font-sans font-light text-xl my-4 text-gray-600">
                Please check if the repo is a valid git repo. If it is not a git
                repo, delete the entry from "Settings" menu and add the repo
                again by checking "Initialize a new repo" option
              </div>
              <div className="my-10 text-gray-200">
                <FontAwesomeIcon
                  icon={["fas", "unlink"]}
                  size="10x"
                ></FontAwesomeIcon>
              </div>
            </div>
          </div>
        ) : null}
      </>
    </>
  );
}
