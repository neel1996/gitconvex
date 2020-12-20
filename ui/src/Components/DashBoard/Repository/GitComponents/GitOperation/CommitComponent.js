import axios from "axios";
import React, { useEffect, useRef, useState } from "react";
import { globalAPIEndpoint } from "../../../../../util/env_config";
import "../../../../styles/GitOperations.css";

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
            <div className="git-ops--commit--wrapper">
              <div className="git-ops--commit--header">
                {stagedCount} {stagedCount > 1 ? "Changes" : "Change"} to
                commit...
              </div>
              <div className="overflow-auto" style={{ height: "300px" }}>
                {stagedFilesState.map((stagedFile) => {
                  return (
                    <div
                      className="git-ops--commit--stagedfiles"
                      key={stagedFile}
                    >
                      {stagedFile}
                    </div>
                  );
                })}
              </div>
              <div className="text-xl my-4">Commit Message</div>
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
                className="git-ops--commit--message"
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
                <div className="git-ops--commit--alert--failed">
                  Commit Failed!
                </div>
              ) : null}
              {loadingCommit ? (
                <div className="git-ops--commit--alert--progress">
                  Committing Changes...
                </div>
              ) : (
                <div
                  className="git-ops--commit--btn"
                  onClick={(event) => {
                    const commitMsg = commitRef.current.value;

                    if (commitMsg) {
                      commitHandler(commitMsg);
                    } else {
                      alert("Commit message can't be empty");
                    }
                  }}
                >
                  Commit Changes
                </div>
              )}
            </div>
          ) : (
            <div className="git-ops--commit--alert--success">
              All changes have been committed!
            </div>
          )}
        </>
      );
    }
  }

  return (
    <div className="git-ops--commit">
      {stagedFilesState && stagedFilesState.length > 0 ? (
        commitComponent()
      ) : (
        <div className="git-ops--commit--alert">
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
