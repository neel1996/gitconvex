import React from "react";

export default function CompareActiveRepoPane(props) {
  return (
    <div className="compare--active-repo">
      <div className="rounded p-2 text-center font-sans font-semibold text-xl">
        Selected Repository
      </div>
      <div className="compare--active-repo--name">{props.repoName}</div>
    </div>
  );
}
