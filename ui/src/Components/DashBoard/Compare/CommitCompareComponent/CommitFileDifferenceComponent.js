import { library } from "@fortawesome/fontawesome-svg-core";
import { far } from "@fortawesome/free-regular-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../util/env_config";

export default function CommitFileDifferenceComponent(props) {
  library.add(fas, far);
  const { repoId, baseCommit, compareCommit } = props;

  const [fileDifference, setFileDifference] = useState([]);
  const [error, setError] = useState(false);
  const [warn, setWarn] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setFileDifference([]);
    setError(false);
    setLoading(true);

    if (baseCommit === compareCommit) {
      setLoading(false);
      setFileDifference([]);
      return;
    }

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            query {
                commitCompare (repoId: "${repoId}", baseCommit: "${baseCommit}", compareCommit: "${compareCommit}"){
                    type
                    fileName
                }
            }
          `,
      },
    })
      .then((res) => {
        setLoading(false);

        const difference = res.data.data.commitCompare;
        if (difference) {
          setFileDifference([...difference]);
        } else {
          setError(true);
          setWarn("Error occurred while comparing the selected commits!");
        }
      })
      .catch((err) => {
        setLoading(false);
        setError(true);
        console.log(err);
      });
  }, [repoId, baseCommit, compareCommit]);

  return (
    <div className="my-4 w-11/12 mx-auto p-6 rounded shadow bg-blue-50">
      {baseCommit === compareCommit ? (
        <div className="text-center font-sans font-light text-2xl">
          Same commits cannot be compared
        </div>
      ) : null}
      {error ? (
        <div className="p-3 w-full rounded bg-red-100 border font-sans font-light text-xl">
          Error occurred while fetching comparison results!
          {warn ? (
            <div className="p-2 my-4 mx-auto rounded shadow">
              <div className="text-3xl my-2 font-sans font-light text-yellow-900">
                Warning
              </div>
              {warn.map((msg) => {
                const warnMsg = msg.replace("warning: ", "");
                return (
                  <div className="text-xl font-sans font-semibold text-yellow-800">
                    {warnMsg}
                    {warnMsg.includes("diff.renameLimit") ? (
                      <div className="my-4 mx-2 p-3 rounded bg-white text-green-600 font-sans font-light">
                        run
                        <span className="font-semibold font-mono mx-2">
                          `git config diff.renamelimit 0`
                        </span>
                        from command line to fix this problem
                      </div>
                    ) : null}
                  </div>
                );
              })}
            </div>
          ) : null}
        </div>
      ) : null}
      {!error && fileDifference.length > 0 && !loading ? (
        <>
          <div className="text-left font-sans font-semibold text-2xl mx-2 my-4 text-gray-800">
            Differing Files
          </div>
          {fileDifference.map((diff) => {
            const { type, fileName } = diff;
            let iconSelector = "";
            let colorSelector = "";
            let title = "";
            switch (type[0]) {
              case "M":
                iconSelector = "plus-square";
                colorSelector = "text-yellow-600";
                title = "Modified";
                break;
              case "A":
                iconSelector = "plus-square";
                colorSelector = "text-green-500";
                title = "Added";
                break;
              case "D":
                iconSelector = "minus-square";
                colorSelector = "text-red-500";
                title = "Deleted";
                break;
              case "R":
                iconSelector = "caret-square-right";
                colorSelector = "text-indigo-500";
                title = "Renamed";
                break;
              default:
                iconSelector = "stop-circle";
                colorSelector = "text-gray-500";
                title = "Unmerged / Copied / Unknown";
                break;
            }

            return (
              <div
                className="flex items-center align-middle justify-center gap-4"
                key={type + "-" + fileName}
              >
                <div
                  className={`text-2xl cursor-pointer ${colorSelector}`}
                  title={title}
                >
                  <FontAwesomeIcon
                    icon={["far", iconSelector]}
                  ></FontAwesomeIcon>
                </div>
                <div
                  className="w-3/4 font-sans font-light truncate"
                  title={fileName}
                >
                  {fileName}
                </div>
              </div>
            );
          })}
        </>
      ) : null}
      {loading ? (
        <div className="my-2 text-2xl font-sans font-semibold text-gray-500">
          Loading comparison results...
        </div>
      ) : null}
    </div>
  );
}
