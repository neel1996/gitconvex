import axios from "axios";
import React, { useEffect, useRef, useState } from "react";
import { globalAPIEndpoint } from "../../../../../util/env_config";
import InfiniteLoader from "../../../../Animations/InfiniteLoader";
import "../../../../styles/GitOperations.css";

export default function PushComponent(props) {
  const { repoId } = props;

  const [remoteData, setRemoteData] = useState();
  const [currentBranch, setCurrentBranch] = useState("");
  const [isRemoteSet, setIsRemoteSet] = useState(false);
  const [unpushedCommits, setUnpushedCommits] = useState([]);
  const [isCommitEmpty, setIsCommitEmpty] = useState(false);

  const [pushDone, setPushDone] = useState(false);
  const [pushFailed, setPushFailed] = useState(false);
  const [loading, setLoading] = useState(false);

  const remoteRef = useRef();
  const branchRef = useRef();

  useEffect(() => {
    axios({
      url: globalAPIEndpoint,
      method: "POST",
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
                }
            }
          `,
      },
    })
      .then((res) => {
        const repoDetails = res.data.data.gitRepoStatus;
        setCurrentBranch(repoDetails.gitCurrentBranch);
        setRemoteData(repoDetails);
      })
      .catch((err) => {
        console.log(err);
      });
  }, [props]);

  function getUnpushedCommits() {
    const remoteHost = remoteRef.current.value.trim();
    const branchName = currentBranch;

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          query 
          {
            gitUnPushedCommits(repoId: "${props.repoId}", remoteURL: "${remoteHost}", remoteBranch: "${branchName}")
          }
        `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          const commits = res.data.data.gitUnPushedCommits;
          if (commits.length === 0) {
            setIsCommitEmpty(true);
          }
          setUnpushedCommits([...commits]);
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }

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
      const { gitCurrentBranch } = remoteData;

      if (gitCurrentBranch) {
        return (
          <option
            disabled
            hidden
            value={gitCurrentBranch}
            key={gitCurrentBranch}
          >
            {gitCurrentBranch}
          </option>
        );
      }
      return null;
    }
  }

  function commitModel(commit) {
    const modelLabel = [
      "Commit Hash",
      "Commit Author",
      "Commit Timestamp",
      "Commit Message",
    ];
    const splitCommit = commit.split("||");

    const localModelFormat = (left, right) => {
      return (
        <div className="flex justify-evenly" key={left}>
          <div className="font-sans text-gray-900 font-bold mx-2 w-1/4 break-words">
            {left}
          </div>
          <div className="font-sans text-gray-800 mx-2 w-1/2 break-words">
            {right}
          </div>
        </div>
      );
    };

    return (
      <div className="block justify-evenly border-b-2 p-2">
        {modelLabel.map((label, index) => {
          return localModelFormat(label, splitCommit[index]);
        })}
      </div>
    );
  }

  function pushHandler(remote) {
    setLoading(true);
    setPushFailed(false);
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          mutation {
            pushToRemote(repoId: "${repoId}", remoteHost: "${remote}", branch: "${currentBranch}")
          }
        `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          const pushStatus = res.data.data.pushToRemote;
          if (pushStatus !== "PUSH_FAILED") {
            setPushDone(true);
            setLoading(false);
          } else {
            setPushFailed(true);
            setLoading(false);
          }
        } else {
          setPushFailed(true);
          setLoading(false);
        }
      })
      .catch((err) => {
        setPushFailed(true);
        setLoading(false);
      });
  }

  return (
    <>
      {!pushDone ? (
        <>
          <div className="w-11/12 xl:w-3/4 lg:w-3/4 mx-auto my-auto bg-gray-50 px-10 py-6 rounded-lg">
            <div className="text-center mb-4 font-sans font-semibold mx-auto w-full p-3 text-2xl border-b-2 border-dashed text-gray-800">
              Push To Remote
            </div>
            <div className="flex w-full mx-auto justify-around items-center align-middle gap-4 my-6">
              <div className="w-2/3 font-sans text-xl font-semibold text-gray-600">
                Available remotes
              </div>
              <div className="w-3/4">
                <select
                  className="outline-none p-3 font-sans text-xl text-gray-600 rounded shadow-md border border-blue-100"
                  defaultValue="checked"
                  onChange={() => {
                    setIsRemoteSet(true);
                    getUnpushedCommits();
                  }}
                  ref={remoteRef}
                  disabled={remoteData ? false : true}
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
              <div className="my-4 flex w-full mx-auto justify-around items-center align-middle gap-4">
                <div className="w-2/3 font-sans text-xl font-semibold text-gray-600">
                  Commits will be pushed to 
                </div>
                <div className="w-3/4">
                  <select
                    disabled
                    className="outline-none p-3 font-sans text-xl text-gray-600 rounded shadow-md border border-blue-100"
                    defaultValue={remoteData.gitCurrentBranch}
                    onChange={() => {
                      getUnpushedCommits();
                    }}
                    ref={branchRef}
                  >
                    {remoteData ? branchListGenerator() : null}
                  </select>
                </div>
              </div>
            ) : null}

            {unpushedCommits && unpushedCommits.length > 0 ? (
              <div className="my-2 mx-auto w-full bg-purple-50 p-4 rounded-lg shadow border">
                <div className="font-sans font-semibold text-2xl text-gray-600 my-1">
                  {unpushedCommits.length !== 0 ? (
                    <span className="mx-1 border-b border-dashed border-gray-600">
                      {unpushedCommits.length}
                    </span>
                  ) : null}
                  {unpushedCommits.length === 1 ? "Commit " : "Commits "}
                  to be pushed
                </div>
                <div className="overflow-auto" style={{ height: "200px" }}>
                  {unpushedCommits.map((commits, index) => {
                    return (
                      <div key={`unpushed-commit-${index}`} className="p-2">
                        {commitModel(commits)}
                      </div>
                    );
                  })}
                </div>
              </div>
            ) : (
              <div></div>
            )}

            {pushFailed ? (
              <>
                <div className="my-4 bg-red-100 text-red-400 text-center text-2xl font-sans font-semibold border-b-4 border-dashed border-red-400 p-4 rounded-lg">
                  Failed to push changes!
                </div>
              </>
            ) : null}

            {isRemoteSet && unpushedCommits.length > 0 && !loading ? (
              <div
                className="w-full text-center font-sans font-semibold p-3 my-4 bg-indigo-300 text-xl text-white rounded shadow-md cursor-pointer hover:bg-indigo-400 hover:shadow-sm"
                onClick={() => {
                  const remoteHost = remoteRef.current.value.trim();

                  if (remoteHost) {
                    pushHandler(remoteHost);
                  }
                }}
              >
                PUSH CHANGES
              </div>
            ) : (
              <>
                {isCommitEmpty ? (
                  <div className="mt-4 rounded shadow-md p-3 bg-gray-400 text-white font-sans font-semibold text-xl text-center">
                    No Commits to Push
                  </div>
                ) : null}
              </>
            )}
            <>
              {loading ? (
                <div className="my-2 p-3 bg-green-50 rounded-lg shadow text-green-500 font-sans font-light text-center text-2xl border-2 border-dashed border-green-500">
                  <div className="text-green-500 text-2xl">
                    Pushing to remote...
                  </div>
                  <div className="flex mx-auto my-6 text-center justify-center">
                    <InfiniteLoader loadAnimation={loading}></InfiniteLoader>
                  </div>
                </div>
              ) : null}
            </>
          </div>
        </>
      ) : (
        <div className="w-11/12 xl:w-3/4 lg:w-3/4 mx-auto my-auto p-6 bg-white rounded">
          <div className="p-6 border-b-4 border-dashed bg-green-100 border-green-500 text-center rounded-lg shadow font-sans text-green-500 text-2xl font-semibold">
            Changes have been pushed to remote
          </div>
        </div>
      )}
    </>
  );
}
