import { library } from "@fortawesome/fontawesome-svg-core";
import { fab } from "@fortawesome/free-brands-svg-icons";
import { far } from "@fortawesome/free-regular-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import { format } from "date-fns";
import debounce from "lodash.debounce";
import React, { useEffect, useRef, useState } from "react";
import ReactDOM from "react-dom";
import { globalAPIEndpoint } from "../../../../../../util/env_config";
import { relativeCommitTimeCalculator } from "../../../../../../util/relativeCommitTimeCalculator";
import InfiniteLoader from "../../../../../Animations/InfiniteLoader";
import CommitLogFileCard from "./CommitLogFileCard";

export default function RepositoryCommitLogComponent(props) {
  library.add(fab, fas, far);

  const [commitLogs, setCommitLogs] = useState([]);
  const [isCommitEmpty, setIsCommitEmpty] = useState(false);
  const [skipLimit, setSkipLimit] = useState(0);
  const [totalCommitCount, setTotalCommitCount] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const [excessCommit, setExcessCommit] = useState(false);
  const [searchKey, setSearchKey] = useState("");
  const [viewReload, setViewReload] = useState(0);
  const [searchWarning, setSearchWarning] = useState(false);
  const [referenceCommitHash, setReferenceCommitHash] = useState("");

  const searchRef = useRef();
  const searchOptionRef = useRef();

  const debouncedSearch = useRef(
    debounce(commitSearchHandler, 1500, { maxWait: 2000 })
  ).current;

  const searchOptions = ["Commit Hash", "Commit Message", "User"];

  useEffect(() => {
    setIsLoading(true);
    setSearchWarning(false);

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            query {
              gitCommitLogs(repoId: "${props.repoId}", referenceCommit: "") {
                  totalCommits
                  commits{
                      commitTime
                      hash
                      author
                      commitMessage
                      commitFilesCount
                  }  
              }
          }
          `,
      },
    })
      .then((res) => {
        setIsLoading(false);

        if (res.data.data) {
          const { commits, totalCommits } = res.data.data.gitCommitLogs;

          if (totalCommits <= 10) {
            setExcessCommit(false);
          } else {
            setExcessCommit(true);
          }

          setTotalCommitCount(totalCommits);
          if (commits && commits.length > 0) {
            setCommitLogs([...commits]);
            const len = commits.length;
            setReferenceCommitHash(commits[len - 1].hash);
          } else {
            setIsCommitEmpty(true);
          }
        }
      })
      .catch((err) => {
        setIsLoading(false);

        if (err) {
          setIsCommitEmpty(true);
          console.log(err);
        }
      });
  }, [props, viewReload]);

  function fetchCommitLogs() {
    setIsLoading(true);
    setSearchWarning(false);

    let localLimit = 0;
    localLimit = skipLimit + 10;

    setSkipLimit(localLimit);

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          query{
            gitCommitLogs(repoId:"${props.repoId}", referenceCommit: "${referenceCommitHash}"){
                totalCommits
                commits{
                    commitTime
                    hash
                    author
                    commitMessage
                    commitFilesCount
                }  
            }
        }
          `,
      },
    })
      .then((res) => {
        setIsLoading(false);

        if (totalCommitCount - localLimit <= 10) {
          setExcessCommit(false);
        }

        if (res.data.data) {
          const { commits, totalCommits } = res.data.data.gitCommitLogs;
          setTotalCommitCount(totalCommits);
          if (commits && commits.length > 0) {
            setCommitLogs([...commitLogs, ...commits]);

            const len = commits.length;
            setReferenceCommitHash(commits[len - 1].hash);
          } else {
            setIsCommitEmpty(true);
          }
        }
      })
      .catch((err) => {
        setIsLoading(false);

        if (err) {
          setIsCommitEmpty(true);
          console.log(err);
        }
      });
  }

  function fetchCommitFiles(commitHash, arrowTarget) {
    const parentDivId = `commitLogCard-${commitHash}`;
    const targetDivId = `commitFile-${commitHash}`;

    const targetDiv = document.createElement("div");
    targetDiv.id = targetDivId;

    const parentDiv = document.getElementById(parentDivId);
    parentDiv.append(targetDiv);

    const unmountHandler = () => {
      ReactDOM.unmountComponentAtNode(
        document.getElementById("closeBtn-" + commitHash)
      );
      ReactDOM.unmountComponentAtNode(document.getElementById(targetDivId));
      arrowTarget.classList.remove("hidden");
    };

    ReactDOM.render(
      <CommitLogFileCard
        repoId={props.repoId}
        commitHash={commitHash}
        unmountHandler={unmountHandler}
      ></CommitLogFileCard>,
      document.getElementById(targetDivId)
    );

    const closeArrow = (
      <div
        className="text-center mx-auto text-3xl font-sans font-light text-gray-600 items-center align-middle cursor-pointer"
        onClick={(event) => {
          unmountHandler();
        }}
      >
        <FontAwesomeIcon icon={["fas", "angle-up"]}></FontAwesomeIcon>
      </div>
    );

    const closeBtn = document.createElement("div");
    const closeBtnId = "closeBtn-" + commitHash;
    closeBtn.id = closeBtnId;
    parentDiv.append(closeBtn);

    ReactDOM.render(closeArrow, document.getElementById(closeBtnId));
  }

  function commitSearchHandler() {
    setIsLoading(true);
    setTotalCommitCount(0);
    setCommitLogs([]);
    const searchQuery = searchRef.current.value;
    let searchOption = "";

    switch (searchOptionRef.current.value) {
      case "Commit Hash":
        searchOption = "hash";
        break;
      case "Commit Message":
        searchOption = "message";
        break;
      case "User":
        searchOption = "user";
        break;
      default:
        searchOption = "message";
        break;
    }

    if (searchQuery) {
      axios({
        url: globalAPIEndpoint,
        method: "POST",
        data: {
          query: `
            query{
              searchCommitLogs(repoId:"${props.repoId}",searchType:"${searchOption}",searchKey:"${searchQuery}"){
                hash
                author
                commitTime
                commitMessage
                commitFilesCount
              }
            }
          `,
        },
      })
        .then((res) => {
          if (res.data.data) {
            const { searchCommitLogs } = res.data.data;
            if (searchCommitLogs && searchCommitLogs.length > 0) {
              setIsCommitEmpty(false);
              setExcessCommit(false);
              setCommitLogs([...searchCommitLogs]);
              setTotalCommitCount(searchCommitLogs.length);
              setIsLoading(false);
            } else {
              setIsCommitEmpty(true);
              setCommitLogs([]);
              setTotalCommitCount(0);
              setIsLoading(false);
              setSearchWarning(true);
            }
          }
        })
        .catch((err) => {
          console.log(err);
          setIsLoading(false);
          setCommitLogs([]);
        });
    } else {
      setViewReload(viewReload + 1);
      setIsLoading(false);
    }
  }

  function fallBackComponent(message) {
    return (
      <div className="p-6 rounded-md shadow-sm block justify-center mx-auto my-auto w-3/4 h-full text-center text-2xl text-indigo-500">
        <div className="flex w-full h-full mx-auto my-auto">
          <div className="block my-auto mx-auto bg-white w-full p-6 rounded-lg shadow">
            <div className="text-2xl text-center font-sans font-semibold text-indigo-800 border-b-2 border-dashed border-indigo-500 p-1">
              {message}
            </div>
            {searchWarning ? (
              <div className="my-4 mx-auto rounded shadow p-4 text-center font-sans text-yellow-800 font-light bg-yellow-50 border-b-4  border-dashed border-yellow-200 text-md">
                Make sure if you are searching with the right category and the
                right search query
              </div>
            ) : null}
            {isLoading ? (
              <div className="flex mx-auto my-6 text-center justify-center">
                <InfiniteLoader loadAnimation={isLoading}></InfiniteLoader>
              </div>
            ) : null}
          </div>
        </div>
      </div>
    );
  }

  function searchbarComponent() {
    return (
      <div className="my-4 w-full rounded-lg bg-white shadow-inner flex gap-4 justify-between items-center">
        <select
          defaultValue="default-search"
          id="searchOption"
          ref={searchOptionRef}
          className="w-1/4 flex p-4 items-center bg-indigo-400 text-white cursor-pointer rounded-l-md text-lg font-sans font-semibold outline-none"
        >
          <option value="default-search" hidden disabled>
            Search for...
          </option>
          {searchOptions.map((item) => {
            return (
              <option key={item} value={item}>
                {item}
              </option>
            );
          })}
        </select>

        <div className="w-3/4 rounded-r-md">
          <input
            ref={searchRef}
            type="text"
            className="w-5/6 outline-none text-lg font-light font-sans"
            placeholder="What are you looking for?"
            value={searchKey}
            onChange={(event) => {
              setSearchKey(event.target.value);
              debouncedSearch();
            }}
          />
        </div>
        <div
          className="w-20 bg-gray-200 p-3 mx-auto my-auto text-center rounded-r-lg hover:bg-gray-400 cursor-pointer"
          onClick={() => {
            commitSearchHandler();
          }}
        >
          <FontAwesomeIcon
            icon={["fas", "search"]}
            className="text-3xl text-gray-600"
          ></FontAwesomeIcon>
        </div>
      </div>
    );
  }

  return (
    <>
      {searchbarComponent()}
      {(isCommitEmpty || !commitLogs || !totalCommitCount) && !isLoading
        ? fallBackComponent("No Commit Logs found")
        : null}
      {commitLogs &&
        commitLogs.map((commit) => {
          const {
            hash,
            author,
            commitTime,
            commitMessage,
            commitFilesCount,
          } = commit;

          let commitRelativeTime = relativeCommitTimeCalculator(commitTime);
          const formattedCommitTime = format(
            new Date(commitTime),
            "MMMM dd, yyyy"
          );

          return (
            <div
              id={`commitLogCard-${hash}`}
              className="p-6 rounded-lg block shadow-md justify-center mx-auto my-4 bg-white w-full border-b-8 border-indigo-400"
              key={hash}
            >
              <div className="flex justify-between text-indigo-400">
                <div className="text-2xl font-sans mx-auto">
                  <FontAwesomeIcon
                    icon={["fas", "calendar-alt"]}
                  ></FontAwesomeIcon>
                  <span className="border-b-2 border-dashed mx-2">
                    {formattedCommitTime}
                  </span>
                </div>
                <div className="h-auto p-1 border-r"></div>
                <div className="text-2xl font-sans mx-auto">
                  <FontAwesomeIcon
                    icon={["fab", "slack-hash"]}
                  ></FontAwesomeIcon>
                  <span className="border-b-2 border-dashed mx-2">
                    {hash.substring(0, 7)}
                  </span>
                </div>
                <div className="h-auto p-1 border-r"></div>
                <div className="text-2xl font-sans mx-auto">
                  <FontAwesomeIcon
                    icon={["fas", "user-ninja"]}
                  ></FontAwesomeIcon>
                  <span className="border-b-2 border-dashed mx-2 truncate">
                    {author}
                  </span>
                </div>
              </div>

              <div className="font-sans font-semibold text-2xl my-4 text-gray-500 p-3 flex justify-evenly items-center">
                <div className="w-1/8">
                  <FontAwesomeIcon
                    icon={["fas", "code"]}
                    className="text-3xl"
                  ></FontAwesomeIcon>
                </div>
                <div className="w-5/6 mx-3">{commitMessage}</div>
              </div>

              <div className="w-11/12 flex justify-between mx-auto mt-4 font-sans text-xl text-gray-500">
                <div className="w-1/3 flex justify-center my-auto items-center align-middle">
                  <div>
                    <FontAwesomeIcon icon={["far", "clock"]}></FontAwesomeIcon>
                  </div>
                  <div className="mx-2 border-dashed border-b-4">
                    {commitRelativeTime}
                  </div>
                </div>
                <div
                  className="w-1/3 flex justify-around my-auto font-sans text-3xl font-light pt-10 cursor-pointer text-gray-500"
                  onClick={(event) => {
                    if (commitFilesCount) {
                      event.currentTarget.classList.add("hidden");
                      fetchCommitFiles(hash, event.currentTarget);
                    }
                  }}
                >
                  {commitFilesCount ? (
                    <FontAwesomeIcon
                      icon={["fas", "angle-down"]}
                    ></FontAwesomeIcon>
                  ) : (
                    <FontAwesomeIcon
                      icon={["fas", "dot-circle"]}
                      className="text-xl text-gray-200"
                    ></FontAwesomeIcon>
                  )}
                </div>
                <div className="w-1/3 flex justify-center my-auto items-center align-middle">
                  <div>
                    <FontAwesomeIcon
                      icon={["far", "plus-square"]}
                    ></FontAwesomeIcon>
                  </div>
                  <div className="mx-2 border-dashed border-b-4">
                    {commitFilesCount ? (
                      `${commitFilesCount} Files`
                    ) : (
                      <span className="text-gray-500">No Changed Files</span>
                    )}
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      {excessCommit ? (
        <div
          className="fixed flex bottom-0 right-0 w-16 h-16 mx-auto p-6 rounded-full shadow-md text-center bg-indigo-500 text-white text-2xl mb-6 mr-6 cursor-pointer"
          title="Click to load commits"
          onClick={() => {
            if (commitLogs.length > skipLimit) {
              fetchCommitLogs();
            }
          }}
        >
          <FontAwesomeIcon
            icon={["fas", "angle-double-down"]}
          ></FontAwesomeIcon>
        </div>
      ) : null}
      {isLoading && totalCommitCount ? (
        <div className="my-4 rounded-lg p-3 bg-gray-100 text-lg font-semibold font-sans text-gray-700 text-center mx-auto">
          Loading {totalCommitCount - skipLimit} more commits...
          <div className="flex mx-auto my-6 text-center justify-center">
            <InfiniteLoader loadAnimation={isLoading}></InfiniteLoader>
          </div>
        </div>
      ) : null}
      {!isCommitEmpty && commitLogs.length === 0 && isLoading
        ? fallBackComponent("Loading commits...")
        : null}
    </>
  );
}
