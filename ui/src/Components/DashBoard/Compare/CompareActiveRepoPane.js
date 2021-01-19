import React from "react";

export default function CompareActiveRepoPane(props) {
  return (
    <div className="w-11/12 bg-indigo-100 border-indigo-400 rounded-lg border-dashed border flex items-center justify-around my-4 mx-auto p-3 shadow gap-10">
      <div className="rounded p-2 text-center font-sans font-semibold text-xl">
        Selected Repository
      </div>
      <div className="bg-indigo-400 border-orange-400 rounded-lg border-dashed border font-semibold text-xl p-4 text-center text-white shadow font-sans">
        {props.repoName}
      </div>
    </div>
  );
}
