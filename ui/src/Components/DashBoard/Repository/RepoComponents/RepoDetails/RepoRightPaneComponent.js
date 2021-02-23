import { faPlus } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useEffect, useState } from "react";
import { v4 as uuid } from "uuid";
import { actionType } from "./backdropActionType";

export default function RepoRightPaneComponent(props) {
  const {
    gitBranchList,
    gitCurrentBranch,
    gitLatestCommit,
    gitTotalCommits,
    actionTrigger,
    switchBranchHandler,
  } = props;

  const [branchValid, setBranchValid] = useState(false);

  useEffect(() => {
    const invalidBranchList = "No Branches available";
    const invalidCurrentBranch = "Repo HEAD is nil";

    if (gitBranchList && gitCurrentBranch && gitBranchList.length > 0) {
      if (
        gitBranchList[0] !== invalidBranchList &&
        gitCurrentBranch !== invalidCurrentBranch
      ) {
        setBranchValid(true);
      }
    }
  }, [gitBranchList, gitCurrentBranch]);

  function addBranchButton() {
    return (
      <div
        id="addBranch"
        className="rounded-full flex cursor-pointer items-center h-10 text-xl mx-auto shadow text-center text-white align-middle w-10 bg-green-400 hover:bg-green-500"
        title="Add a new branch"
        onClick={() => {
          actionTrigger(actionType.ADD_BRANCH);
        }}
      >
        <FontAwesomeIcon
          className="mx-auto text-center"
          icon={faPlus}
        ></FontAwesomeIcon>
      </div>
    );
  }

  return (
    <>
      {props.received ? (
        <div className="border-gray-300 rounded-md border-dotted border-2 block my-6 mx-auto p-1 shadow-sm xl:w-1/2 lg:w-3/4 md:w-11/12 sm:w-11/12 w-11/12 ">
          <div className="block mx-auto my-2">
            <div className="flex justify-around my-3">
              <div className="text-lg text-gray-500 w-1/4">Total Commits</div>
              <div className="font-bold text-left text-gray-800 w-1/2">
                {`${gitTotalCommits} Commits`}
              </div>
            </div>

            <div className="flex justify-around my-3">
              <div className="text-lg text-gray-500 w-1/4">Latest Commit</div>
              <div
                className="font-bold text-sm text-left text-gray-900 truncate w-1/2"
                title={gitLatestCommit}
              >
                {gitLatestCommit}
              </div>
            </div>

            <div className="flex justify-around mx-auto my-2 align-middle items-center">
              <div className="text-lg text-gray-500 w-1/4">
                Available Branches
              </div>

              {branchValid ? (
                <div className="flex items-center justify-evenly align-middle w-1/2">
                  <div className="w-3/4 my-auto">
                    <div
                      className="border-dotted border-b cursor-pointer font-semibold text-lg my-1 text-indigo-500 hover:text-indigo-600"
                      key={`${gitCurrentBranch}-${uuid()}`}
                    >
                      {gitCurrentBranch}
                    </div>
                    {gitBranchList &&
                      gitCurrentBranch &&
                      gitBranchList
                        .slice(0, 2)
                        .map((entry) => {
                          if (entry) {
                            if (entry !== gitCurrentBranch) {
                              return (
                                <div
                                  className="border-dotted border-b cursor-pointer font-semibold my-2 hover:text-indigo-400 font-sans"
                                  key={`entry-key-${uuid()}`}
                                  onClick={() => {
                                    switchBranchHandler(entry);
                                    actionTrigger(actionType.SWITCH_BRANCH);
                                  }}
                                >
                                  {entry}
                                </div>
                              );
                            } else {
                              return null;
                            }
                          }
                          return null;
                        })
                        .filter((item) => {
                          if (item) {
                            return item;
                          }
                          return false;
                        })}
                    <div
                      className="border-dashed border-b cursor-pointer my-auto text-center text-blue-500 hover:text-blue-800 font-sans"
                      onClick={() => {
                        actionTrigger(actionType.LIST_BRANCH);
                      }}
                    >
                      List all branches
                    </div>
                  </div>
                  {addBranchButton()}
                </div>
              ) : (
                <div className="flex w-1/2 justify-between items-center">
                  <div className="w-1/2 font-sans font-light text-gray-400 border-b-2 border-dashed border-gray-600 text-center">
                    No branches available
                  </div>
                  {addBranchButton()}
                </div>
              )}
            </div>

            <div className="flex justify-center mx-auto mt-6 gap-4">
              <div
                className="w-5/12 font-semibold rounded text-center cursor-pointer p-3 bg-indigo-400 hover:bg-indigo-500 text-white font-sans truncate"
                onClick={() => {
                  actionTrigger(actionType.FETCH);
                }}
              >
                FETCH FROM REMOTE
              </div>
              <div
                className="w-5/12 font-semibold text-center cursor-pointer rounded text-white p-3 bg-blue-400 hover:bg-blue-500 font-sans truncate"
                onClick={() => {
                  actionTrigger(actionType.PULL);
                }}
              >
                PULL FROM REMOTE
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </>
  );
}
