import axios from "axios";
import React, { useEffect, useRef, useState } from "react";
import {
  globalAPIEndpoint,
  ROUTE_GIT_UNPUSHED_COMMITS,
  ROUTE_REPO_DETAILS,
} from "../../../../../util/env_config";
import InfiniteLoader from "../../../../Animations/InfiniteLoader";
import "../../../../styles/GitOperations.css";

export default function PushComponent(props) {
  const { repoId } = props;

  const [remoteData, setRemoteData] = useState();
  const [isRemoteSet, setIsRemoteSet] = useState(false);
  const [isBranchSet, setIsBranchSet] = useState(false);
  const [unpushedCommits, setUnpushedCommits] = useState([]);
  const [isCommitEmpty, setIsCommitEmpty] = useState(false);

  const [pushDone, setPushDone] = useState(false);
  const [pushFailed, setPushFailed] = useState(false);
  const [loading, setLoading] = useState(false);

  const remoteRef = useRef();
  const branchRef = useRef();

  useEffect(() => {
    let payload = JSON.stringify(JSON.stringify({ repoId: props.repoId }));

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      data: {
        query: `
                query GitConvexApi
                {
                  gitConvexApi(route: "${ROUTE_REPO_DETAILS}", payload: ${payload}){
                    gitRepoStatus {
                      gitRemoteData
                      gitCurrentBranch
                      gitRemoteHost
                      gitBranchList 
                    }
                  }
                }
              `,
      },
    })
      .then((res) => {
        const repoDetails = res.data.data.gitConvexApi.gitRepoStatus;
        setRemoteData(repoDetails);
      })
      .catch((err) => {
        console.log(err);
      });
  }, [props]);

  function getUnpushedCommits() {
    let payload = JSON.stringify(JSON.stringify({ repoId: props.repoId }));

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          query GitConvexApi
          {
            gitConvexApi(route: "${ROUTE_GIT_UNPUSHED_COMMITS}", payload: ${payload}){
              gitUnpushedCommits{
                commits
              }
            }
          }
        `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          const { commits } = res.data.data.gitConvexApi.gitUnpushedCommits;
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

  function pushHandler(remote, branch) {
    setLoading(true);
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          mutation GitConvexMutation{
            pushToRemote(repoId: "${repoId}", remoteHost: "${remote}", branch: "${branch}")
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
            <div className="git-ops--push--header">Available remote repos</div>
            <div>
              <select
                className="git-ops--push--select"
                defaultValue="checked"
                onChange={() => {
                  setIsRemoteSet(true);
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

            {isRemoteSet ? (
              <div>
                <select
                  className="git-ops--push--select"
                  defaultValue="checked"
                  onChange={() => {
                    setIsBranchSet(true);
                    getUnpushedCommits();
                  }}
                  ref={branchRef}
                >
                  <option disabled hidden value="checked">
                    Select upstream branch
                  </option>
                  {remoteData ? branchListGenerator() : null}
                </select>
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

            {isRemoteSet &&
            isBranchSet &&
            unpushedCommits.length > 0 &&
            !loading ? (
              <div
                className="git-ops--push--btn"
                onClick={() => {
                  const remoteHost = remoteRef.current.value.trim();
                  const branchName = branchRef.current.value.trim();

                  if (remoteHost && branchName) {
                    pushHandler(remoteHost, branchName);
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
          <div className="git-ops--push--alert--success">
            Changes have been pushed to remote
          </div>
        </div>
      )}
    </>
  );
}
