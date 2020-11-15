import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useMemo, useState } from "react";
import {
  globalAPIEndpoint,
  ROUTE_REPO_COMMIT_LOGS,
} from "../../../../util/env_config";
import CommitFileDifferenceComponent from "./CommitFileDifferenceComponent";
import CommitLogCardComponent from "./CommitLogCardComponent";

export default function CommitCompareComponent(props) {
  library.add(fas);

  const [skipCount, setSkipCount] = useState(0);
  const [totalCommitCount, setTotalCommitCount] = useState(0);
  const [commitData, setCommitData] = useState([]);
  const [baseCommit, setBaseCommit] = useState("");
  const [compareCommit, setCompareCommit] = useState("");

  const memoizedCommitFileDifference = useMemo(() => {
    return (
      <CommitFileDifferenceComponent
        repoId={props.repoId}
        baseCommit={baseCommit}
        compareCommit={compareCommit}
      ></CommitFileDifferenceComponent>
    );
  }, [props.repoId, baseCommit, compareCommit]);

  useEffect(() => {
    const payload = JSON.stringify(
      JSON.stringify({ repoId: props.repoId, skipLimit: skipCount })
    );

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            query GitConvexApi
            {
                gitConvexApi(route: "${ROUTE_REPO_COMMIT_LOGS}", payload: ${payload}){
                    gitCommitLogs {
                        totalCommits
                        commits{
                            commitTime
                            hash
                            author
                            commitMessage
                            commitRelativeTime
                            commitFilesCount
                        }  
                    }
                }
            }
            `,
      },
    })
      .then((res) => {
        if (res.data.data) {
          const {
            commits,
            totalCommits,
          } = res.data.data.gitConvexApi.gitCommitLogs;
          setTotalCommitCount(totalCommits);

          setCommitData((data) => {
            if (data) {
              return [...data, ...commits];
            } else {
              return [...commits];
            }
          });
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }, [props.repoId, skipCount]);

  function commitCardComponent(setCommitType) {
    return (
      <>
        {commitData &&
          commitData.map((item) => {
            return (
              <CommitLogCardComponent
                item={item}
                setCommitType={setCommitType}
                key={item.hash}
              ></CommitLogCardComponent>
            );
          })}
        {(skipCount >= 10 || skipCount === 0) &&
        skipCount <= totalCommitCount ? (
          <div
            className="p-3 border cursor-pointer hover:bg-gray-100 text-center font-sans font-semibold"
            onClick={() => {
              setSkipCount(skipCount + 10);
            }}
          >
            Load More commits
          </div>
        ) : null}
      </>
    );
  }

  function baseAndCompareCommitComponent() {
    return (
      <div className="w-11/12 mx-auto my-6 flex gap-10 justify-around">
        {baseCommit ? (
          <div className="flex my-4 gap-10 justify-center items-center">
            <div className="font-sans font-semibold text-xl border-b border-dashed">
              Base Commit
            </div>
            <div className="text-xl font-sans font-semibold p-3 rounded-lg shadow text-gray-600 border-indigo-300 border-2 border-dashed">
              {baseCommit}
            </div>
            <div
              className="p-2 rounded border-b-2 border-dashed shadow cursor-pointer hover:bg-gray-100"
              onClick={() => {
                setBaseCommit("");
              }}
            >
              <FontAwesomeIcon
                className="text-2xl text-gray-600"
                icon={["fas", "edit"]}
              ></FontAwesomeIcon>
            </div>
          </div>
        ) : null}
        {compareCommit ? (
          <div className="flex gap-10 justify-between items-center">
            <div className="font-sans font-semibold text-xl border-b border-dashed">
              Commit to Compare
            </div>
            <div className="text-xl font-sans font-semibold p-3 rounded-lg shadow text-gray-600 border-orange-400 border-2 border-dashed">
              {compareCommit}
            </div>
            <div
              className="p-2 rounded border-b-2 border-dashed shadow cursor-pointer hover:bg-gray-100"
              onClick={() => {
                setCompareCommit("");
              }}
            >
              <FontAwesomeIcon
                className="text-2xl text-gray-600"
                icon={["fas", "edit"]}
              ></FontAwesomeIcon>
            </div>
          </div>
        ) : null}
      </div>
    );
  }

  return (
    <>
      {baseAndCompareCommitComponent()}
      {commitData.length === 0 ? (
        <div className="text-3xl text-center font-sans text-gray-300">
          Loading Commits...
        </div>
      ) : (
        <div className="w-11/12 mx-auto flex gap-10 justify-around">
          {!baseCommit ? (
            <div className="w-1/2 p-2 shadow border rounded">
              <div className="p-2 font-sans font-semibold text-xl font-semibold">
                Select the base Commit
              </div>
              {commitCardComponent(setBaseCommit)}
            </div>
          ) : null}
          {!compareCommit ? (
            <div className="w-1/2 p-2 shadow border rounded">
              <div className="p-2 font-sans font-semibold text-xl font-semibold">
                Select the Commit to compare
              </div>
              {commitCardComponent(setCompareCommit)}
            </div>
          ) : null}
        </div>
      )}
      {baseCommit && compareCommit ? memoizedCommitFileDifference : null}
    </>
  );
}
