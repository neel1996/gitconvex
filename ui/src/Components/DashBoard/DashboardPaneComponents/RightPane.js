import React from "react";
import RepoComponent from "../Repository/RepoComponents/RepoComponent";
import "../../styles/RightPane.css";

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
      <div className="dashboard--rightpane overflow-auto">
        <div className="rightpane--toparea xl:flex lg:flex md:block sm:block">
          {hcParams.map((entry) => {
            return (
              <div key={entry.label} className="rightpane--toparea--hc">
                <div className="rightpane--toparea--hclabel xl:font-bold lg:font-semibold md:font-medium xl:text-2xl lg:text-xl md:text-md">
                  {entry.label}
                </div>
                {entry.value !== "" ? (
                  <div
                    className="rightpane--toparea--hcpills"
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
