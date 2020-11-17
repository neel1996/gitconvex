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
      <div className="block justify-evenly border shadow rounded p-2">
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
          <div className="git-ops--push">
            <div className="flex mx-auto justify-around items-center align-middle gap-4">
              <div className="w-2/3 font-sans text-xl font-semibold text-gray-600">
                Available remotes
              </div>
              <div className="w-3/4">
                <select
                  className="git-ops--push--select"
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
              <div className="flex mx-auto justify-around items-center align-middle gap-4">
                <div className="w-2/3 font-sans text-xl font-semibold text-gray-600">
                  Commits will be pushed
                </div>
                <div className="w-3/4">
                  <select
                    disabled
                    className="git-ops--push--select"
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
              <div className="git-ops--push--unpushed">
                <div className="git-ops--push--unpushed--label">
                  Commits to be pushed
                </div>
                <div className="overflow-auto" style={{ height: "200px" }}>
                  {unpushedCommits.map((commits, index) => {
                    return (
                      <div
                        key={`unpushed-commit-${index}`}
                        className="git-ops--push--commits"
                      >
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
                <div className="git-ops--push--nochange">
                  Failed to push changes!
                </div>
              </>
            ) : null}

            {isRemoteSet && unpushedCommits.length > 0 && !loading ? (
              <div
                className="git-ops--push--btn"
                onClick={() => {
                  const remoteHost = remoteRef.current.value.trim();

                  if (remoteHost) {
                    pushHandler(remoteHost);
                  }
                }}
              >
                Push changes
              </div>
            ) : (
              <>
                {isCommitEmpty ? (
                  <div className="git-ops--push--nocommits">
                    No Commits to Push
                  </div>
                ) : null}
              </>
            )}
            <>
              {loading ? (
                <div className="git-ops--push--loader">
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
        <div className="git-ops--push--success">
          <div className="git-ops--commit--alert--success">
            Changes have been pushed to remote
          </div>
        </div>
      )}
    </>
  );
}