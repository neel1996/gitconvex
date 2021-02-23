import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useContext, useEffect, useMemo, useState } from "react";
import { GIT_GLOBAL_REPOID, PRESENT_REPO } from "../../../../actionStore";
import { ContextProvider } from "../../../../context";
import { globalAPIEndpoint } from "../../../../util/env_config";
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
  const [toggleSearchSelect, setToggleSearchSelect] = useState(false);
  const [searchBranchValue, setSearchBranchValue] = useState("");
  const [filteredBranchList, setFilteredBranchList] = useState([]);

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
          setSearchBranchValue("");
          setFilteredBranchList([]);
          setToggleSearchSelect(!toggleSearchSelect);
          handleScreenEvents();
        }
      })
      .catch((err) => {
        setLoading(false);
        if (err) {
          setBranchError(true);
          event.target.innerText = activeBranch;
        }
      });
  }

  const handleScreenEvents = () => {
    if (!toggleSearchSelect) {
      document
        .getElementById("repository-action")
        .addEventListener("scroll", () => {
          setToggleSearchSelect(false);
        });
    } else {
      document
        .getElementById("repository-action")
        .removeEventListener("scroll", () => {});
    }
  };

  const searchBranchHandler = (e) => {
    const searchBranch = e.target.value;
    setSearchBranchValue(searchBranch);
    if (
      searchBranch !== "" &&
      selectedRepoDetails &&
      selectedRepoDetails.gitBranchList
    ) {
      const { gitBranchList } = selectedRepoDetails;
      const filteredBranches = gitBranchList.filter((branchName) =>
        branchName.toLowerCase().includes(searchBranch)
      );
      setFilteredBranchList(filteredBranches);
    } else {
      setFilteredBranchList([]);
    }
  };

  const cancelSearchBranch = () => {
    setSearchBranchValue("");
    setFilteredBranchList([]);
  };

  function activeRepoPane() {
    return (
      <div className="flex items-center justify-around my-4 mx-auto align-middle">
        <div className="flex items-center">
          <div className="font-sans font-semibold my-1 text-gray-900">
            Choose saved repository
          </div>
          <select
            className="cursor-pointer rounded-lg border-dashed border-b-2 font-sans font-light text-xl mx-4 outline-none p-2 shadow bg-green-50 text-green-700 border-green-400"
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
          <div className="flex items-center gap-4">
            <div className="font-sans font-semibold my-1 text-gray-900">
              Branch
            </div>
            <div className="flex-1 flex flex-col justify-center">
              <div
                className="flex-auto cursor-pointer inline-flex items-center justify-center px-4 py-2 shadow-md bg-indigo-50 border-indigo-400 text-indigo-700 border-dashed border-b-2 truncate"
                onClick={(e) => {
                  let target = e.currentTarget;
                  if (!toggleSearchSelect) {
                    target.style.width = "17.5rem";
                  } else {
                    target.style.width = "auto";
                  }
                  setToggleSearchSelect(!toggleSearchSelect);
                  handleScreenEvents();
                }}
              >
                <span className="mr-2">{activeBranch}</span>
                <FontAwesomeIcon
                  className="text-sm m-1"
                  icon={["fas", "chevron-down"]}
                ></FontAwesomeIcon>
              </div>
              {toggleSearchSelect ? (
                <div className="flex-auto flex flex-row justify-center">
                  <div className="bg-white border-indigo-300 text-indigo-700 px-4 py-4 shadow-md rounded-md z-20 absolute">
                    <div className="flex flex-row mt-1 mb-3">
                      <div className="b-1 text-center px-2 py-1 text-white bg-blue-400 rounded-l-md">
                        <FontAwesomeIcon
                          icon={["fas", "search"]}
                        ></FontAwesomeIcon>
                      </div>
                      <input
                        id="branchSearchInput"
                        type="text"
                        placeholder="Search..."
                        className="px-2 py-1 bg-indigo-100 text-indigo-700 shadow-sm rounded-sm focus:outline-none outline-none"
                        onChange={searchBranchHandler}
                        value={searchBranchValue}
                      ></input>
                      <div
                        className="b-1 text-center px-2 py-1 text-white cursor-pointer bg-red-400 rounded-r-md"
                        onClick={cancelSearchBranch}
                      >
                        <FontAwesomeIcon
                          icon={["fas", "times"]}
                        ></FontAwesomeIcon>
                      </div>
                    </div>
                    {availableBranch()}
                  </div>
                </div>
              ) : null}
            </div>
          </div>
        ) : null}
      </div>
    );
  }

  function getTopPaneComponent(icon, value) {
    return (
      <>
        <div className="border-indigo-400 border-dashed border-b-2 flex justify-between font-sans text-lg mx-2 p-2 text-gray-600">
          <div className="mx-2">
            <FontAwesomeIcon icon={["fas", icon]}></FontAwesomeIcon>
          </div>
          <div className="mx-2">{value}</div>
        </div>
      </>
    );
  }

  const branchCardComponent = (branch) => {
    return (
      <div
        key={branch}
        value={branch}
        className="cursor-pointer text-sm border-b border-dotted p-2 mt-1 mb-1"
        onClick={(event) => {
          event.persist();
          setActiveBranch("...");
          setTrackingBranch(event.target.innerText, event);
        }}
      >
        {branch}
      </div>
    );
  };

  function availableBranch() {
    if (selectedRepoDetails && selectedRepoDetails.gitBranchList) {
      const { gitBranchList } = selectedRepoDetails;
      if (searchBranchValue !== "") {
        if (filteredBranchList.length > 0) {
          return filteredBranchList.map((branch, index) => {
            if (branch !== "NO_BRANCH") {
              return branchCardComponent(branch);
            }
            return null;
          });
        } else {
          return (
            <div className="text-center font-sans font-light text-base my-2 text-indigo-800 border-b border-dotted">
              <span className="mx-1 font-semibold border-b border-dashed">
                {searchBranchValue}
              </span>
              branch is not available!
            </div>
          );
        }
      } else {
        return gitBranchList.map((branch, index) => {
          if (branch !== "NO_BRANCH") {
            return branchCardComponent(branch);
          }
          return null;
        });
      }
    }
  }

  return (
    <div
      className="block justify-center mx-auto overflow-x-hidden w-full"
      id="repository-action"
    >
      {availableRepos ? (
        <div>
          <div className="w-11/12 border-gray-200 rounded border my-6 mx-auto shadow">
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
              <div className="text-center mx-auto font-sans p-10 text-2xl bg-yellow-100 w-11/12 font-light rounded-lg shadow border-dashed border-2 border-yellow-200">
                Select a configured repo from the dropdown to perform git
                related operations
              </div>
              <div className="w-3/4 border-gray-100 rounded-lg border-2 block my-20 mx-auto p-6">
                <div>
                  <FontAwesomeIcon
                    icon={["fas", "mouse-pointer"]}
                    className="flex font-bold h-full text-6xl m-auto text-center text-gray-300 w-full"
                  ></FontAwesomeIcon>
                </div>
                <div className="block my-4 mx-auto text-center text-gray-200 xl:text-6xl lg:text-3xl md:text-2xl">
                  No repositories selected
                </div>
              </div>
            </>
          ) : null}
          <div>
            {branchError ? (
              <div className="bg-red-100 rounded font-sans my-2 mx-auto p-2 text-center text-red-700">
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
