import React from "react";
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

  return (
    <>
      {props.received ? (
        <div className="repo-rightpane xl:w-1/2 lg:w-3/4 md:w-11/12 sm:w-11/12 w-11/12 ">
          <div className="block mx-auto my-2">
            <div className="flex justify-around my-3">
              <div className="text-lg text-gray-500 w-1/4">Total Commits</div>
              <div className="total-commits">
                {`${gitTotalCommits} Commits`}
              </div>
            </div>

            <div className="flex justify-around my-3">
              <div className="text-lg text-gray-500 w-1/4">Latest Commit</div>
              <div className="latest-commit" title={gitLatestCommit}>
                {gitLatestCommit}
              </div>
            </div>

            <div className="flex justify-around mx-auto my-2 align-middle items-center">
              <div className="text-lg text-gray-500 w-1/4">
                Available Branches
              </div>

              <div className="branch-list">
                <div className="w-3/4 my-auto">
                  <div
                    className="branch-list--current"
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
                                className="branch-list--branches"
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
                    className="branch-list__listbranch"
                    onClick={() => {
                      actionTrigger(actionType.LIST_BRANCH);
                    }}
                  >
                    List all branches
                  </div>
                </div>
                <div
                  id="addBranch"
                  className="rounded-full items-center align-middle w-10 h-10 text-white text-2xl bg-green-400 text-center mx-auto shadow hover:bg-green-500 cursor-pointer"
                  onMouseEnter={(event) => {
                    let popUp =
                      '<div class="tooltip" style="margin-left:-40px;">Click to add a new branch</div>';
                    event.target.innerHTML += popUp;
                  }}
                  onMouseLeave={(event) => {
                    event.target.innerHTML = "+";
                  }}
                  onClick={() => {
                    actionTrigger(actionType.ADD_BRANCH);
                  }}
                >
                  +
                </div>
              </div>
            </div>

            <div className="flex justify-around mx-auto mt-4">
              <div
                className="w-1/3 rounded text-center cursor-pointer p-2 bg-indigo-400 hover:bg-indigo-500 text-white font-sans nowrap"
                onClick={() => {
                  actionTrigger(actionType.FETCH);
                }}
              >
                Fetch from remote
              </div>
              <div
                className="w-1/3 text-center cursor-pointer rounded text-white p-2 bg-blue-400 hover:bg-blue-500 font-sans"
                onClick={() => {
                  actionTrigger(actionType.PULL);
                }}
              >
                Pull from remote
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </>
  );
}
