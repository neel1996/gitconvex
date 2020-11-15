import React from "react";
import "../../../../styles/RepositoryDetails.css";

export default function RepoInfoComponent(props) {
  const { gitRepoName, gitCurrentBranch } = props;
  return (
    <div className="repo-info">
      <div className="repo-info--label">Repo Name</div>
      <div className="bg-blue-100 text-blue-900 border-blue-200 repo-info--data">
        {gitRepoName}
      </div>
      <div className="repo-info--label">Active Branch</div>
      <div className="bg-orange-200 text-orange-900 border-orange-400 repo-info--data">
        {gitCurrentBranch}
      </div>
    </div>
  );
}
