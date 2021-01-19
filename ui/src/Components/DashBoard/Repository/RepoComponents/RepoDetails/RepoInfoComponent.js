import React from "react";

export default function RepoInfoComponent(props) {
  const { gitRepoName, gitCurrentBranch } = props;
  return (
    <div className="border-gray-100 rounded-md border-2 flex items-center justify-evenly mx-auto p-4 align-middle w-full shadow">
      <div className="font-semibold text-lg mx-2 p-2 text-gray-800 font-sans">
        Repo Name
      </div>
      <div className="bg-blue-100 text-blue-900 border-blue-200 repo-info--data rounded shadow">
        {gitRepoName}
      </div>
      <div className="font-semibold text-lg mx-2 p-2 text-gray-800 font-sans">
        Active Branch
      </div>
      <div className="bg-yellow-50 text-yellow-800 border-yellow-300 repo-info--data rounded shadow">
        {gitCurrentBranch === "Repo HEAD is nil" ? "---" : gitCurrentBranch}
      </div>
    </div>
  );
}
