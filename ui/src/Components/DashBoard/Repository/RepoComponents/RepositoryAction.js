import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useContext, useEffect, useMemo, useState } from "react";
import { GIT_GLOBAL_REPOID, PRESENT_REPO } from "../../../../actionStore";
import { ContextProvider } from "../../../../context";
import { globalAPIEndpoint } from "../../../../util/env_config";
import "../../../styles/RepositoryAction.css";
import GitTrackedComponent from "../GitComponents/GitTrackedComponent";

export default function RepositoryAction() {
  library.add(fas);

  const { state, dispatch } = useContext(ContextProvider);
  const { presentRepo } = state;
  const [loading, setLoading] = useState(false);
  const [selectedFlag, setSelectedFlag] = useState(false);
  const [defaultRepo, setDefaultRepo] = useState({});
  const [availableRepos, setAvailableRepos] = useState([]);
  const [activeBranch, setActiveBranch] = useState("");
  const [selectedRepoDetails, setSelectedRepoDetails] = useState({
    gitBranchList: "",
    gitCurrentBranch: "",
    gitTotalCommits: 0,
    gitTotalTrackedFiles: 0,
  });
  const [branchError, setBranchError] = useState(false);

  const memoizedGitTracker = useMemo(() => {
    if (defaultRepo && defaultRepo.id) {
      return (
        <GitTrackedComponent
          repoId={defaultRepo.id}
          resetBranchError={() => {
            setBranchError(false);
          }}
        ></GitTrackedComponent>
      );
    }
  }, [defaultRepo]);

  useEffect(() => {
    const token = axios.CancelToken;
    const source = token.source();

    function fetchSelectedRepoStatus() {
      const repoId = defaultRepo && defaultRepo.id;

      if (repoId) {
        setLoading(true);
        axios({
          url: globalAPIEndpoint,
          method: "POST",
          headers: {
            "Content-type": "application/json",
          },
          cancelToken: source.token,
          data: {
            query: `
              query 
              {
                  gitRepoStatus(repoId: "${repoId}") {
                    gitBranchList
                    gitCurrentBranch
                    gitTotalCommits
                    gitTotalTrackedFiles 
                  }
              }
            `,
          },
        })
          .then((res) => {
            setLoading(false);

            setSelectedRepoDetails(res.data.data.gitRepoStatus);
            setActiveBranch(res.data.data.gitRepoStatus.gitCurrentBranch);
          })
          .catch((err) => {
            setLoading(false);

            if (err) {
              console.log("API GitStatus error occurred : " + err);
            }
          });
      }
    }

    //Effect dep function
    async function invokeRepoFetchAPI() {
      setLoading(true);
      return await axios({
        url: globalAPIEndpoint,
        method: "POST",
        cancelToken: source.token,
        data: {
          query: `
              query {
                  fetchRepo{
                    repoId
                    repoName
                    repoPath
                  }
              }
          `,
        },
      }).then((res) => {
        setLoading(false);

        const apiResponse = res.data.data.fetchRepo;

        if (apiResponse) {
          const repoContent = apiResponse.repoId.map((entry, index) => {
            return {
              id: apiResponse.repoId[index],
              repoName: apiResponse.repoName[index],
              repoPath: apiResponse.repoPath[index],
            };
          });

          dispatch({
            type: PRESENT_REPO,
            payload: [...repoContent],
          });

          setDefaultRepo(repoContent[0]);
          setAvailableRepos(repoContent);
          return repoContent;
        }
      });
    }

    if (presentRepo && presentRepo[0]) {
      setAvailableRepos(presentRepo[0]);
      fetchSelectedRepoStatus();
    } else {
      invokeRepoFetchAPI();
      fetchSelectedRepoStatus();
    }

    return () => {
      source.cancel();
    };
  }, [defaultRepo, activeBranch, presentRepo, dispatch, branchError]);

  function setTrackingBranch(branchName, event) {
    setLoading(true);
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          mutation{
            checkoutBranch(repoId: "${defaultRepo.id}", branchName: "${branchName}")
          }
        `,
      },
    })
      .then((res) => {
        setLoading(false);

        if (res.data.data && !res.data.error) {
          setActiveBranch(branchName);
        }
      })
      .catch((err) => {
        setLoading(false);

        if (err) {
          setBranchError(true);
          event.target.selectedIndex = 0;
        }
      });
  }

  function activeRepoPane() {
    return (
      <div className="top-pane">
        <div className="flex items-center">
          <div className="select--label">Choose saved repository</div>
          <select
            className="top-pane--select bg-green-50 text-green-700 border-green-400"
            defaultValue={"checked"}
            onClick={() => {
              setBranchError(false);
            }}
            onChange={(event) => {
              setActiveBranch("...");
              if (event.currentTarget.value !== defaultRepo.repoName) {
                setSelectedRepoDetails({
                  ...selectedRepoDetails,
                  gitCurrentBranch: "",
                  gitBranchList: ["..."],
                });
              }
              setSelectedFlag(true);
              availableRepos.length &&
                availableRepos.forEach((elm) => {
                  if (event.target.value === elm.repoName) {
                    setDefaultRepo(elm);
                    dispatch({ type: GIT_GLOBAL_REPOID, payload: elm.id });
                  }
                });
            }}
          >
            <option defaultChecked value="checked" hidden disabled>
              Select a repo
            </option>
            {availableRepos.length &&
              availableRepos.map((entry) => {
                return (
                  <option value={entry.repoName} key={entry.repoName}>
                    {entry.repoName}
                  </option>
                );
              })}
          </select>
        </div>
        {selectedFlag ? (
          <div className="flex items-center">
            <div className="select--label">Branch</div>
            <select
              value={activeBranch}
              defaultChecked={activeBranch}
              className="top-pane--select bg-indigo-50 border-indigo-300 text-indigo-700"
              disabled={activeBranch ? false : true}
              onChange={(event) => {
                event.persist();
                setActiveBranch("...");
                setTrackingBranch(event.target.value, event);
              }}
              onClick={() => {
                setBranchError(false);
              }}
            >
              <option key={activeBranch} value={activeBranch}>
                {activeBranch}
              </option>
              {availableBranch()}
            </select>
          </div>
        ) : null}
      </div>
    );
  }

  function getTopPaneComponent(icon, value) {
    return (
      <>
        <div className="top-pane--component">
          <div className="mx-2">
            <FontAwesomeIcon icon={["fas", icon]}></FontAwesomeIcon>
          </div>
          <div className="mx-2">{value}</div>
        </div>
      </>
    );
  }

  function availableBranch() {
    if (selectedRepoDetails && selectedRepoDetails.gitBranchList) {
      const { gitBranchList } = selectedRepoDetails;

      return gitBranchList.map((branch, index) => {
        if (branch !== "NO_BRANCH" && branch !== activeBranch) {
          return (
            <option key={branch} value={branch}>
              {branch}
            </option>
          );
        }

        return null;
      });
    }
  }

  return (
    <div className="repository-action">
      {availableRepos ? (
        <div>
          <div className="active-repo">
            {activeRepoPane()}
            {selectedRepoDetails && selectedFlag ? (
              <div className="my-auto flex justify-around p-3 mx-auto">
                {loading ? (
                  <div className="text-center font-sans font-semibold text-gray-600 text-xl">
                    Loading repo details...
                  </div>
                ) : (
                  <>
                    {getTopPaneComponent(
                      "code-branch",
                      selectedRepoDetails.gitBranchList &&
                        selectedRepoDetails.gitBranchList.length > 0 &&
                        !selectedRepoDetails.gitBranchList[0].match(
                          /NO_BRANCH/gi
                        ) ? (
                        <>
                          {selectedRepoDetails.gitBranchList.length === 1
                            ? 1 + " branch"
                            : selectedRepoDetails.gitBranchList.length +
                              " branches"}
                        </>
                      ) : (
                        "No Branches"
                      )
                    )}
                    {getTopPaneComponent(
                      "sort-amount-up",
                      selectedRepoDetails.gitTotalCommits + " Commits"
                    )}
                    {getTopPaneComponent(
                      "archive",
                      selectedRepoDetails.gitTotalTrackedFiles +
                        " Tracked Files"
                    )}
                  </>
                )}
              </div>
            ) : null}
          </div>
          {!selectedFlag ? (
            <>
              <div className="alert--jumbotron">
                Select a configured repo from the dropdown to perform git
                related operations
              </div>
              <div className="alert--message">
                <div>
                  <FontAwesomeIcon
                    icon={["fas", "mouse-pointer"]}
                    className="alert--message--icon"
                  ></FontAwesomeIcon>
                </div>
                <div className="alert--message--label xl:text-6xl lg:text-3xl md:text-2xl">
                  No repositories selected
                </div>
              </div>
            </>
          ) : null}
          <div>
            {branchError ? (
              <div className="alert--failure">
                Branch switching failed.Commit your changes and try again
              </div>
            ) : null}
            {selectedRepoDetails && selectedFlag && defaultRepo.id
              ? memoizedGitTracker
              : null}
          </div>
        </div>
      ) : null}
    </div>
  );
}
