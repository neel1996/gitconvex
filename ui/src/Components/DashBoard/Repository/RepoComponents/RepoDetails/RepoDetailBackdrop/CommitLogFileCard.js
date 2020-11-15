import React, { useState, useEffect } from "react";
import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { far } from "@fortawesome/free-regular-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  ROUTE_COMMIT_FILES,
  globalAPIEndpoint,
} from "../../../../../../util/env_config";
import axios from "axios";
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

    const payload = JSON.stringify(
      JSON.stringify({ repoId: repoId, commitHash: commitHash })
    );

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      cancelToken: source.token,
      data: {
        query: `
                query GitConvexApi
                {
                    gitConvexApi(route: "${ROUTE_COMMIT_FILES}", payload: ${payload}){
                        gitCommitFiles {
                            type
                            fileName
                        }
                    }
                }
            `,
      },
    })
      .then((res) => {
        setIsLoading(false);
        if (res.data.data && !res.data.err) {
          setCommitFiles([...res.data.data.gitConvexApi.gitCommitFiles]);
        }
      })
      .catch((err) => {
        console.log(err);
        setIsLoading(false);
      });
  }, [repoId, commitHash]);

  return (
    <div className="commitlogs--files">
      <div
        className="commitlogs--files--close"
        style={{ marginTop: "-20px" }}
        onClick={() => {
          setCommitFiles([]);
          unmountHandler();
        }}
      >
        x
      </div>
      {isLoading ? (
        <div className="commitlogs--files--header text-center">
          Fetching changed files...
        </div>
      ) : null}
      {!isLoading && commitFiles ? (
        <div className="commitlogs--files--header">{`${commitFiles.length} Files changed`}</div>
      ) : null}
      <div className="commitlogs--files--list">
        {commitFiles &&
          commitFiles.map(({ type, fileName }) => {
            let iconSelector = "";
            let colorSelector = "";
            switch (type) {
              case "M":
                iconSelector = "plus-square";
                colorSelector = "text-yellow-600";
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
                className="commitlogs--files--list--item"
                key={fileName + commitHash}
              >
                <div className={`w-1/4 text-2xl ${colorSelector}`}>
                  <FontAwesomeIcon
                    icon={["far", iconSelector]}
                  ></FontAwesomeIcon>
                </div>
                <div className="commitlogs--files--list--filename">
                  {fileName}
                </div>
              </div>
            );
          })}
      </div>
    </div>
  );
}
