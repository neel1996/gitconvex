import React from "react";
import RepoComponent from "../Repository/RepoComponents/RepoComponent";

export default function RightPane(props) {
  const { platform, gitVersion } = props.params;

  const hcParams = [
    {
      label: "Platform",
      value: platform,
    },
    {
      label: "Git",
      value: gitVersion,
    },
  ];

  return (
    <>
      <div className="w-full mx-auto overflow-auto">
        <div className="w-11/12 p-3 my-6 rounded-lg shadow-md justify-center mx-auto bg-blue-50 border border-blue-200 border-dashed block xl:flex lg:flex md:block sm:block">
          {hcParams.map((entry) => {
            return (
              <div
                key={entry.label}
                className="my-2 flex mx-auto gap-10 justify-around items-center align-middle"
              >
                <div className="w-1/2 border-b-2 border-dashed text-center font-sans xl:font-bold lg:font-semibold md:font-medium xl:text-2xl lg:text-xl md:text-md">
                  {entry.label}
                </div>
                {entry.value !== "" ? (
                  <div
                    className="w-2/3 mx-2 bg-green-100 border border-green-200 text-center p-2 rounded-lg"
                    id={`hc-param__${entry.label}`}
                  >
                    {entry.value}
                  </div>
                ) : (
                  <div className="rounded-md bg-red-200 text-red-900 font-bold p-2">
                    Invalid
                  </div>
                )}
              </div>
            );
          })}
        </div>
        {repoEntry()}
      </div>
    </>
  );

  function repoEntry() {
    if (platform && gitVersion) {
      return <RepoComponent parentProps={props}></RepoComponent>;
    }
  }
}
