import React from "react";
import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export default function CommitLogCardComponent(props) {
  library.add(fas);
  const { item, setCommitType } = props;

  return (
    <div
      className="flex p-3 justify-between items-center mx-auto border-b w-full cursor-pointer hover:bg-gray-100"
      key={item.hash}
      onClick={() => {
        setCommitType(item.hash);
      }}
    >
      <div className="block p-2 font-sans font-light text-gray-800">
        <div className="my-2 font-sans text-xl font-light text-blue-600 border-b border-dashed">
          {item.commitTime.split(" ")[0]}
        </div>
        <div className="w-3/4">{item.commitMessage}</div>
        <div className="flex items-center gap-4 my-2 align-middle">
          <div>
            <FontAwesomeIcon
              icon={["fas", "user-alt"]}
              className="text-indigo-500"
            ></FontAwesomeIcon>
          </div>
          <div className="font-semibold font-sans">{item.author}</div>
        </div>
      </div>
      <div className="block">
        <div className="shadow border rounded text-sm p-2 bg-indigo-100 font-mono font-semibold text-indigo-800">
          #{item.hash.substring(0, 7)}
        </div>
        <div className="my-2 shadow border rounded text-sm p-2 bg-orange-100 font-sans font-semibold">
          {item.commitRelativeTime}
        </div>
      </div>
    </div>
  );
}
