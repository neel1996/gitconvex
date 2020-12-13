import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useMemo, useState } from "react";
import { globalAPIEndpoint } from "../../../../util/env_config";
import BranchCommitLogChanges from "./BranchCommitLogChanges";

export default function BranchCompareComponent(props) {
  library.add(fas);
  const [branchList, setBranchList] = useState([]);
  const [currentBranch, setCurrentBranch] = useState("");
  const [compareBranch, setCompareBranch] = useState("");
  const [baseBranch, setBaseBranch] = useState("");
  const [errState, setErrState] = useState(false);

  const memoizedBranchCommitLogChangesComponent = useMemo(() => {
    return (
      <BranchCommitLogChanges
        repoId={props.repoId}
        baseBranch={baseBranch}
        compareBranch={compareBranch}
      ></BranchCommitLogChanges>
    );
  }, [props.repoId, baseBranch, compareBranch]);

  useEffect(() => {
    const token = axios.CancelToken;
    const source = token.source();
    setErrState(false);

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      cancelToken: source.token,
      data: {
        query: `
            query 
            {
                gitRepoStatus(repoId: "${props.repoId}") {
                    gitBranchList  
                    gitCurrentBranch
                }
            }
          `,
      },
    })
      .then((res) => {
        let { gitBranchList, gitCurrentBranch } = res.data.data.gitRepoStatus;

        if (gitBranchList.length <= 0 || gitCurrentBranch === "") {
          setErrState(true);
        }

        gitBranchList =
          gitBranchList &&
          gitBranchList.map((branch) => {
            return branch.trim();
          });

        if (gitBranchList.length > 1) {
          setCompareBranch(gitBranchList[1].trim());
        }

        setBranchList(gitBranchList);
        setCurrentBranch(gitCurrentBranch.trim());
        setBaseBranch(gitCurrentBranch.trim());
      })
      .catch((err) => {
        console.log(err);
        setErrState(true);
      });

    return () => {
      return source.cancel;
    };
  }, [props.repoId]);

  function noBranchToCompare() {
    return (
      <div className="w-full mx-auto my-auto text-center block">
        <FontAwesomeIcon
          icon={["fas", "puzzle-piece"]}
          className="font-sans text-center text-gray-300 my-20"
          size="10x"
        ></FontAwesomeIcon>
        <div className="text-2xl text-gray-300">
          Only one branch is available, hence can't be set for comparison
        </div>
      </div>
    );
  }

  function compareBranchSelectPane() {
    return (
      <div className="w-11/12 p-3 flex mx-auto items-center align-middle rounded-lg shadow-md border-2 justify-around">
        <div className="flex gap-6 justify-between items-center">
          <div className="text-xl text-center font-sans font-semibold border-b-2 border-dashed border-gray-400">
            Base branch
          </div>
          <div>
            <select
              defaultValue={currentBranch}
              className="outline-none p-2 shadow border-2 bg-white rounded-lg"
              onChange={(e) => {
                setBaseBranch(e.currentTarget.value.trim());
              }}
            >
              <option value={currentBranch}>{currentBranch}</option>
              {branchList.map((branch) => {
                if (branch !== currentBranch) {
                  return (
                    <option value={branch} key={branch}>
                      {branch}
                    </option>
                  );
                }
                return null;
              })}
            </select>
          </div>
        </div>
        <div className="flex gap-6 justify-between items-center">
          <div className="text-xl text-center font-sans font-semibold border-b-2 border-dashed border-gray-400">
            Compare branch
          </div>
          <div>
            <select
              className="outline-none p-2 shadow border-2 bg-white rounded-lg"
              onChange={(e) => {
                setCompareBranch(e.currentTarget.value);
              }}
            >
              {branchList.map((branch, index) => {
                if (baseBranch && baseBranch !== branch) {
                  return (
                    <option value={branch} key={branch}>
                      {branch}
                    </option>
                  );
                } else {
                  return index === 0 ? null : (
                    <option value={branch} key={branch}>
                      {branch}
                    </option>
                  );
                }
              })}
            </select>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div>
      {branchList.length === 1 ? (
        noBranchToCompare()
      ) : branchList.length === 0 && !errState ? (
        <div className="mx-auto my-20 text-center flex justify-center text-4xl font-sans text-gray-300">
          Loading Branch Info...
        </div>
      ) : !errState ? (
        compareBranchSelectPane()
      ) : null}
      {baseBranch && compareBranch && !errState
        ? memoizedBranchCommitLogChangesComponent
        : null}

      {errState ? (
        <div className="mx-auto text-center text-2xl text-gray-500 font-sans font-semibold p-4 border-b border-dashed border-gray-400">
          Error occurred while fetching results. Please verify if the repo has
          valid branches
        </div>
      ) : null}
    </div>
  );
}
