import { library } from "@fortawesome/fontawesome-svg-core";
import { far } from "@fortawesome/free-regular-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../../../util/env_config";
export default function CommitLogFileCard({
  repoId,
  commitHash,
  unmountHandler,
}) {
  library.add(far, fas);
  const [commitFiles, setCommitFiles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  useEffect(() => {
    setIsLoading(true);
    const token = axios.CancelToken;
    const source = token.source();

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      cancelToken: source.token,
      data: {
        query: `
            query
            {
              gitCommitFiles(repoId: "${repoId}", commitHash: "${commitHash}"){
                    type
                    fileName
                }
              }
          `,
      },
    })
      .then((res) => {
        setIsLoading(false);
        if (res.data.data && !res.data.err) {
          setCommitFiles([...res.data.data.gitCommitFiles]);
        }
      })
      .catch((err) => {
        console.log(err);
        setIsLoading(false);
      });
  }, [repoId, commitHash]);

  return (
    <div className="w-11/12 p-6 rounded-lg shadow block mx-auto my-6 bg-blue-50">
      <div
        className="font-sans font-light float-right right-0 cursor-pointer mx-auto text-2xl text-blue-400 mb-0"
        style={{ marginTop: "-20px" }}
        onClick={() => {
          setCommitFiles([]);
          unmountHandler();
        }}
      >
        x
      </div>
      {isLoading ? (
        <div className="mx-4 text-2xl font-sans font-light text-gray-600 text-center">
          Fetching changed files...
        </div>
      ) : null}
      {!isLoading && commitFiles ? (
        <div className="mx-4 text-2xl font-sans font-light text-gray-600">{`${commitFiles.length} Files changed`}</div>
      ) : null}
      <div className="block w-3/4 mx-10 my-4">
        {commitFiles &&
          commitFiles.map(({ type, fileName }) => {
            let iconSelector = "";
            let colorSelector = "";
            switch (type) {
              case "M":
                iconSelector = "plus-square";
                colorSelector = "text-yellow-400";
                break;
              case "A":
                iconSelector = "plus-square";
                colorSelector = "text-green-500";
                break;
              case "D":
                iconSelector = "minus-square";
                colorSelector = "text-red-500";
                break;
              default:
                iconSelector = "plus-square";
                colorSelector = "text-gray-500";
                break;
            }

            return (
              <div
                className="flex justify-evenly items-center align-middle my-auto"
                key={fileName + commitHash}
              >
                <div className={`w-1/4 text-2xl ${colorSelector}`}>
                  <FontAwesomeIcon
                    icon={["far", iconSelector]}
                  ></FontAwesomeIcon>
                </div>
                <div
                  className="truncate w-3/4 font-medium text-sm text-gray-600"
                  title={fileName}
                >
                  {fileName}
                </div>
              </div>
            );
          })}
      </div>
    </div>
  );
}
