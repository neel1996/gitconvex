import axios from "axios";
import React, { useEffect, useRef, useState } from "react";
import { globalAPIEndpoint } from "../../../../../util/env_config";

export default function CommitComponent(props) {
  const { repoId } = props;

  const [loading, setLoading] = useState(true);
  const [stagedFilesState, setStagedFilesState] = useState([]);
  const [commitDone, setCommitDone] = useState(false);
  const [commitError, setCommitError] = useState(false);
  const [loadingCommit, setLoadingCommit] = useState(false);
  const [commitMessageWarning, setCommitMessageWarning] = useState(false);

  const commitRef = useRef();

  useEffect(() => {
    setLoading(true);

    const cancelToken = axios.CancelToken;
    const source = cancelToken.source();

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      cancelToken: source.token,
      data: {
        query: `
          query {
            gitChanges(repoId: "${props.repoId}"){
              gitStagedFiles
            }
          }
          `,
      },
    })
      .then((res) => {
        const { gitStagedFiles } = res.data.data.gitChanges;
        setLoading(false);

        if (gitStagedFiles && gitStagedFiles.length > 0) {
          setStagedFilesState([...gitStagedFiles]);
        }
      })
      .catch((err) => {
        setLoading(false);
        console.log(err);
      });

    return () => {
      source.cancel();
    };
  }, [props]);

  function commitHandler(commitMsg) {
    setLoadingCommit(true);
    commitMsg = commitMsg.replace(/"/gi, '"');
    if (commitMsg.split("\n") && commitMsg.split("\n").length > 0) {
      commitMsg = commitMsg.toString().split("\n").join("||");
    }

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          mutation {
            commitChanges(repoId: "${repoId}", commitMessage: "${commitMsg}")
          }
        `,
      },
    })
      .then((res) => {
        setLoadingCommit(false);

        if (
          res.data.data &&
          !res.data.error &&
          res.data.data.commitChanges === "COMMIT_DONE"
        ) {
          setCommitDone(true);
        } else {
          setCommitError(true);
        }
      })
      .catch((err) => {
        setLoadingCommit(false);
        setCommitError(true);
      });
  }

  function commitComponent() {
    if (stagedFilesState && stagedFilesState.length > 0) {
      const stagedCount = stagedFilesState.length;

      return (
        <>
          {!commitDone ? (
            <div className="w-11/12 mx-auto my-auto p-4">
              <div className="font-font-sans font-semibold text-3xl my-4 text-gray-700">
                {stagedCount} {stagedCount > 1 ? "Changes" : "Change"} to
                commit...
              </div>
              <div className="overflow-auto" style={{ height: "300px" }}>
                {stagedFilesState.map((stagedFile) => {
                  return (
                    <div
                      className="font-sans text-gray-700 my-1 mx-6 border-b"
                      key={stagedFile}
                    >
                      {stagedFile}
                    </div>
                  );
                })}
              </div>
              <div className="text-xl my-4 font-sans font-semibold text-gray-600">
                Commit Message
              </div>
              {commitMessageWarning ? (
                <div className="font-sans font-semibold italic p-2 border-b border-dotted border-yellow-400 text-yellow-500">
                  <span role="img" aria-label="suggestion">
                    💡
                  </span>
                  <span className="mx-1">
                    It is usually a good practice to limit the commit message to
                    50 characters
                  </span>
                  <div className="my-1 font-sans text-sm font-semibold text-yellow-700">
                    For additional content, include a line break and enter the
                    messages
                  </div>
                </div>
              ) : null}
              <textarea
                className="w-full outline-none rounded-lg p-4 shadow-lg border border-blue-100"
                placeholder="Enter commit message"
                cols="20"
                rows="5"
                ref={commitRef}
                onChange={(e) => {
                  const content = e.currentTarget.value;
                  const len = content.split("\n")[0].length;
                  if (len > 50) {
                    setCommitMessageWarning(true);
                  } else {
                    setCommitMessageWarning(false);
                  }
                }}
              ></textarea>
              {commitError ? (
                <div className="w-full mx-auto my-2 bg-red-200 border-b-2 border-red-400 rounded rounded-b-lg p-3 text-center text-xl text-red-600 font-semibold font-sans">
                  Commit Failed!
                </div>
              ) : null}
              {loadingCommit ? (
                <div className="font-sans font-semibold my-1 mx-auto p-2 text-center text-xl bg-gray-400 shadow-md w-full cursor-pointer rounded-lg text-white">
                  Committing Changes...
                </div>
              ) : (
                <div
                  className="w-full my-2 p-3 rounded shadow bg-green-500 font-sans font-semibold text-center text-white text-xl cursor-pointer hover:bg-green-600"
                  onClick={(event) => {
                    const commitMsg = commitRef.current.value;

                    if (commitMsg) {
                      commitHandler(commitMsg);
                    } else {
                      alert("Commit message can't be empty");
                    }
                  }}
                >
                  COMMIT CHANGES
                </div>
              )}
            </div>
          ) : (
            <div className="mx-auto my-2 p-3 bg-green-200 text-green-600 font-sans font-semibold text-2xl text-center border-b-4 border-dashed border-green-500 rounded">
              All changes have been committed!
            </div>
          )}
        </>
      );
    }
  }

  return (
    <div className="w-11/12 xl:w-3/4 lg:w-3/4 mx-auto my-auto bg-gray-100 rounded-lg shadow p-5">
      {stagedFilesState && stagedFilesState.length > 0 ? (
        commitComponent()
      ) : (
        <div className="bg-gray-200 p-3 text-center font-sans text-2xl mx-auto my-2 rounded font-light border-b-4 border-dashed border-gray-400">
          {loading ? (
            <span>Loading staged files to commit...</span>
          ) : stagedFilesState.length === 0 ? (
            <span>No Staged files to commit</span>
          ) : (
            <span>Loading...</span>
          )}
        </div>
      )}
    </div>
  );
}
