import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useMemo, useState } from "react";
import { globalAPIEndpoint } from "../../../../util/env_config";
import CommitFileDifferenceComponent from "./CommitFileDifferenceComponent";
import CommitLogCardComponent from "./CommitLogCardComponent";

export default function CommitCompareComponent(props) {
  library.add(fas);

  const [skipCount, setSkipCount] = useState(0);
  const [totalCommitCount, setTotalCommitCount] = useState(0);
  const [commitData, setCommitData] = useState([]);
  const [baseCommit, setBaseCommit] = useState("");
  const [compareCommit, setCompareCommit] = useState("");
  const [errState, setErrState] = useState(false);

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
    setErrState(false);
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            query
            {
                gitCommitLogs (repoId: "${props.repoId}", skipLimit: ${skipCount}){
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
            `,
      },
    })
      .then((res) => {
        if (res.data.data) {
          const { commits, totalCommits } = res.data.data.gitCommitLogs;

          if (totalCommits === 0 || totalCommits == null) {
            setErrState(true);
          }

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
        setErrState(true);
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
              {baseCommit.substring(0, 7)}
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
              {compareCommit.substring(0, 7)}
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
      {commitData.length === 0 && !errState ? (
        <div className="text-3xl text-center font-sans text-gray-300">
          Loading Commits...
        </div>
      ) : !errState ? (
        <div className="w-11/12 mx-auto flex gap-10 justify-around">
          {!baseCommit ? (
            <div className="w-1/2 p-2 shadow border rounded">
              <div className="p-2 font-sans font-semibold text-xl">
                Select the base Commit
              </div>
              {commitCardComponent(setBaseCommit)}
            </div>
          ) : null}
          {!compareCommit ? (
            <div className="w-1/2 p-2 shadow border rounded">
              <div className="p-2 font-sans text-xl font-semibold">
                Select the Commit to compare
              </div>
              {commitCardComponent(setCompareCommit)}
            </div>
          ) : null}
        </div>
      ) : (
        <div className="mx-auto text-center text-2xl text-gray-500 font-sans font-semibold p-4 border-b border-dashed border-gray-400">
          Error occurred while fetching results. Please verify if the repo has
          valid branches
        </div>
      )}
      {baseCommit && compareCommit ? memoizedCommitFileDifference : null}
    </>
  );
}
