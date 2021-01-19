import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React from "react";

export default function CompareSelectionHint() {
  library.add(fas);

  const selectionHints = [
    {
      message: (
        <div className="my-10 font-sans text-xl font-semibold">
          Select
          <span className="mx-2 p-1 rounded-lg text-center text-white bg-gray-400">
            Branch Compare
          </span>
          for comparing two branches
        </div>
      ),
      icon: "code-branch",
    },
    {
      message: "",
      icon: "",
    },
    {
      message: (
        <div className="my-10 font-sans text-xl font-semibold">
          Select
          <span className="mx-2 p-1 rounded-lg text-center text-white bg-gray-400">
            Commit Compare
          </span>
          for comparing two selected commits
        </div>
      ),
      icon: "hashtag",
    },
  ];

  return (
    <div className="flex justify-around my-6 mx-auto text-gray-300 gap-10 w-11/12">
      {selectionHints.map((hint, index) => {
        if (hint.message) {
          return (
            <div className="w-1/3 block text-center" key={`hint-${index}`}>
              <div>
                <FontAwesomeIcon
                  icon={["fas", hint.icon]}
                  size="10x"
                ></FontAwesomeIcon>
              </div>
              {hint.message}
            </div>
          );
        } else {
          return (
            <div
              className="block p-1 border-r-4 boreder-dashed"
              key={`hint-${index}`}
            ></div>
          );
        }
      })}
    </div>
  );
}
