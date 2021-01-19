import { library } from "@fortawesome/fontawesome-svg-core";
import { far } from "@fortawesome/free-regular-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useEffect, useMemo, useState } from "react";
import BranchCompareComponent from "./BranchCompareComponent/BranchCompareComponent";
import CommitCompareComponent from "./CommitCompareComponent/CommitCompareComponent";
import CompareActionButtons from "./CompareActionButtons";
import CompareActiveRepoPane from "./CompareActiveRepoPane";
import CompareSelectionHint from "./CompareSelectionHint";
import RepoSearchBar from "./RepoSearchBar";

export default function CompareComponent() {
  library.add(fas, far);

  const [selectedRepo, setSelectedRepo] = useState("");
  const [compareAction, setCompareAction] = useState("");

  useEffect(() => {
    setCompareAction("");
  }, [selectedRepo.id]);

  function activateCompare(repo) {
    setSelectedRepo(repo);
  }

  const memoizedBranchCompareComponent = useMemo(() => {
    return (
      <BranchCompareComponent repoId={selectedRepo.id}></BranchCompareComponent>
    );
  }, [selectedRepo.id]);

  const memoizedCommitCompareComponent = useMemo(() => {
    return (
      <CommitCompareComponent repoId={selectedRepo.id}></CommitCompareComponent>
    );
  }, [selectedRepo.id]);

  const memoizedCompareActionButtons = useMemo(() => {
    return (
      <CompareActionButtons
        selectedRepo={selectedRepo.id}
        compareAction={compareAction}
        setCompareAction={(action) => {
          setCompareAction(action);
        }}
      ></CompareActionButtons>
    );
  }, [compareAction, selectedRepo.id]);

  function noSelectedRepobanner() {
    return (
      <div className="block m-auto text-center w-full">
        <FontAwesomeIcon
          icon={["far", "object-group"]}
          className="font-sans text-center text-gray-300 my-20"
          size="10x"
        ></FontAwesomeIcon>
        <div className="text-6xl text-gray-200">Select a Repo to compare</div>
      </div>
    );
  }

  return (
    <div className="h-full py-10 w-full overflow-auto">
      <div className="font-sans font-light text-3xl mx-10 text-gray-800">
        Compare Branches and Commits
      </div>
      <RepoSearchBar activateCompare={activateCompare}></RepoSearchBar>
      {selectedRepo ? (
        <>
          <CompareActiveRepoPane
            repoName={selectedRepo.repoName}
          ></CompareActiveRepoPane>
          {memoizedCompareActionButtons}
          {compareAction ? (
            compareAction === "branch-compare" ? (
              memoizedBranchCompareComponent
            ) : (
              <>
                {compareAction === "commit-compare"
                  ? memoizedCommitCompareComponent
                  : null}
              </>
            )
          ) : (
            <CompareSelectionHint></CompareSelectionHint>
          )}
        </>
      ) : (
        noSelectedRepobanner()
      )}
    </div>
  );
}
