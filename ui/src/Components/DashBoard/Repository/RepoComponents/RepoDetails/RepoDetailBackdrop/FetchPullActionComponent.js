import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useRef, useState } from "react";
import { v4 as uuid } from "uuid";
import { globalAPIEndpoint } from "../../../../../../util/env_config";
import InfiniteLoader from "../../../../../Animations/InfiniteLoader";

export default function FetchFromRemoteComponent(props) {
  library.add(fas);
  const { repoId, actionType } = props;

  const [remoteData, setRemoteData] = useState();
  const [isRemoteSet, setIsRemoteSet] = useState(false);
  const [isBranchSet, setIsBranchSet] = useState(false);
  const [result, setResult] = useState([]);
  const [loading, setLoading] = useState(false);

  const remoteRef = useRef();
  const branchRef = useRef();

  useEffect(() => {
    const cancelToken = axios.CancelToken;
    const source = cancelToken.source();

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      cancelToken: source.token,
      headers: {
        "Content-type": "application/json",
      },
      data: {
        query: `
                query
                {
                    gitRepoStatus(repoId:"${props.repoId}") {
                      gitRemoteData
                      gitCurrentBranch
                      gitRemoteHost
                      gitBranchList 
                    }
                }
              `,
      },
    })
      .then((res) => {
        const repoDetails = res.data.data.gitRepoStatus;
        setRemoteData(repoDetails);
      })
      .catch((err) => {
        setLoading(false);
      });

    return () => {
      return source.cancel();
    };
  }, [props]);

  function remoteHostGenerator() {
    if (remoteData) {
      const { gitRemoteData } = remoteData;
      if (gitRemoteData.includes("||")) {
        return gitRemoteData.split("||").map((item) => {
          return (
            <option value={item} key={item}>
              {item}
            </option>
          );
        });
      } else {
        return <option>{gitRemoteData}</option>;
      }
    }
  }

  function branchListGenerator() {
    if (remoteData) {
      const { gitBranchList } = remoteData;

      return gitBranchList.map((branch) => {
        if (branch !== "NO_BRANCH") {
          return (
            <option value={branch} key={branch}>
              {branch}
            </option>
          );
        }
        return null;
      });
    }
  }

  function actionHandler(remote = "", branch = "") {
    setLoading(true);

    const getAxiosRequestBody = (remote, branch) => {
      let gqlQuery = "";
      if (actionType === "fetch") {
        gqlQuery = `mutation {
          fetchFromRemote(repoId: "${repoId}", remoteUrl: "${remote}", remoteBranch: "${branch}"){
            status
            fetchedItems
          }
        }
      `;
      } else {
        gqlQuery = `mutation {
          pullFromRemote(repoId: "${repoId}", remoteUrl: "${remote}", remoteBranch: "${branch}"){
            status
            pulledItems
          }
        }
      `;
      }

      return gqlQuery;
    };

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: getAxiosRequestBody(remote, branch),
      },
    })
      .then((res) => {
        setLoading(false);
        if (res.data.data && !res.data.error) {
          let actionResponse = {};

          if (actionType === "fetch") {
            actionResponse = res.data.data.fetchFromRemote;
          } else {
            actionResponse = res.data.data.pullFromRemote;
          }

          if (actionResponse.status.match(/ABSENT/gi)) {
            setResult([
              <div className="text-xl text-center border-2 border-dashed border-gray-800 p-2 text-gray-700 font-semibold">
                No changes to {actionType === "fetch" ? "Fetch" : "Pull"} from
                remote
              </div>,
            ]);
          } else if (actionResponse.status.match(/ERROR/gi)) {
            setResult([
              <div className="text-xl p-2 text-pink-800 border border-pink-200 shadow rounded font-semibold">
                Error while {actionType === "fetch" ? "Fetching" : "Pulling"}{" "}
                from remote!
              </div>,
            ]);
          } else {
            let resArray = [];
            if (actionType === "fetch") {
              resArray = actionResponse.fetchedItems;
            } else {
              resArray = actionResponse.pulledItems;
            }
            setResult([
              <div className="text-xl text-center border-2 border-dashed border-green-600 p-2 text-green-700 bg-green-200 font-semibold rounded shadow">
                {resArray[0]}
              </div>,
            ]);
          }
        }
      })
      .catch((err) => {
        setLoading(false);
        console.error(err);
        setResult([
          <div className="text-xl p-2 text-pink-800 border border-pink-200 shadow rounded font-semibold">
            Error while {actionType === "fetch" ? "Fetching" : "Pulling"} from
            remote!
          </div>,
        ]);
      });
  }

  return (
    <>
      <div className="w-3/4 mx-auto my-auto shadow rounded-lg bg-white pt-4">
        {actionType === "fetch" ? (
          <div>
            <div className="text-center font-sans font-semibold mx-auto w-full p-3 text-2xl border-b-2 border-dashed text-gray-800">
              Fetch from Remote
            </div>
            <div
              className="flex justify-center items-center w-11/12 mx-auto my-4 text-center p-1 font-sans font-medium text-lg cursor-pointer text-indigo-400 hover:text-indigo-500 xl:w-3/5 lg:w-3/4 md:w-3/4 sm:w-11/12"
              onClick={() => {
                actionHandler();
              }}
            >
              <div className="text-2xl text-indigo-800 mx-4">
                <FontAwesomeIcon
                  icon={["fas", "exclamation-circle"]}
                ></FontAwesomeIcon>
              </div>
              <div>Click to Fetch from default remote and branch</div>
            </div>
          </div>
        ) : null}
        {actionType === "pull" ? (
          <div className="text-center font-sans font-semibold mx-auto w-full p-3 text-2xl border-b-2 border-dashed text-gray-800">
            Pull from Remote
          </div>
        ) : null}
        <div className="flex flex-wrap w-3/4 mx-auto my-4 justify-around items-center align-middle gap-4">
          <div className="w-full font-sans text-xl font-semibold text-gray-600">
            Available remotes
          </div>
          <div className="w-full mb-6">
            <select
              className="border p-3 text-lg rounded shadow font-sans outline-none"
              defaultValue="checked"
              disabled={remoteData ? false : true}
              onChange={() => {
                setIsRemoteSet(true);
              }}
              onClick={() => {
                setResult([]);
              }}
              ref={remoteRef}
            >
              <option disabled hidden value="checked">
                {remoteData
                  ? "Select the remote repo"
                  : "Loading available remotes..."}
              </option>
              {remoteData ? remoteHostGenerator() : null}
            </select>
          </div>
        </div>

        {isRemoteSet ? (
          <div className="flex flex-wrap w-3/4 mx-auto my-4 justify-around items-center align-middle gap-4">
            <div className="w-full font-sans text-xl font-semibold text-gray-600">
              Available Branches
            </div>
            <div className="w-full mb-6">
              <select
                className="border p-3 text-lg rounded shadow font-sans outline-none"
                defaultValue="checked"
                onChange={() => {
                  setIsBranchSet(true);
                }}
                onClick={() => {
                  setResult([]);
                }}
                ref={branchRef}
              >
                <option disabled hidden value="checked">
                  Select upstream branch
                </option>
                {remoteData ? branchListGenerator() : null}
              </select>
            </div>
          </div>
        ) : null}

        {isRemoteSet && isBranchSet && !loading ? (
          <div
            className="text-center font-semibold text-xl p-4 mx-auto cursor-pointer bg-indigo-400 text-white mt-10 rounded-b hover:bg-indigo-500"
            onClick={(event) => {
              const remoteHost = remoteRef.current.value;
              const branchName = branchRef.current.value;

              if (remoteHost && branchName) {
                actionHandler(remoteHost, branchName);
              } else {
                event.target.style.display = "none";
              }
            }}
          >
            {actionType === "pull" ? "PULL FROM REMOTE" : null}
            {actionType === "fetch" ? "FETCH FROM REMOTE" : null}
          </div>
        ) : null}
        <div>
          {loading ? (
            <>
              <div className="my-4 text-center border-2 border-dashed border-gray-500 font-light font-sans text-xl p-2 mx-auto text-gray-500 bg-gray-50 shadow rounded-md">
                {actionType === "pull" ? "Pulling changes" : "Fetching"} from
                remote...
              </div>
              <div className="flex mx-auto my-6 text-center justify-center">
                <InfiniteLoader loadAnimation={loading}></InfiniteLoader>
              </div>
            </>
          ) : null}
        </div>

        {true || (!loading && result && result.length > 0) ? (
          <>
            {result.map((result) => {
              return (
                <div
                  className="my-1 mx-2 text-center text-xl font-sans shadow bg-gray-300"
                  key={result + `-${uuid()}`}
                >
                  {result}
                </div>
              );
            })}
          </>
        ) : null}
      </div>
    </>
  );
}
