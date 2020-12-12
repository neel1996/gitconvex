import React from "react";
import "../../../../styles/RepositoryDetails.css";

export default function RepoInfoComponent(props) {
  const { gitRepoName, gitCurrentBranch } = props;
  return (
    <div className="repo-info">
      <div className="repo-info--label">Repo Name</div>
      <div className="bg-blue-100 text-blue-900 border-blue-200 repo-info--data rounded shadow">
        {gitRepoName}
      </div>
      <div className="repo-info--label">Active Branch</div>
      <div className="bg-yellow-50 text-yellow-800 border-yellow-300 repo-info--data rounded shadow">
        {gitCurrentBranch}
      </div>
    </div>
  );
}
