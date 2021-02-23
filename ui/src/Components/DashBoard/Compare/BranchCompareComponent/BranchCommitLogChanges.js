import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../util/env_config";
import { format } from "date-fns";

export default function BranchCommitLogChanges(props) {
  library.add(fas);
  const { repoId, baseBranch, compareBranch } = props;

  const [commitLogs, setCommitLogs] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setCommitLogs([]);
    setLoading(true);

    if (compareBranch === baseBranch) {
      setLoading(false);
      return;
    }

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          query GitConvexApi {
              branchCompare ( repoId: "${repoId}", baseBranch: "${baseBranch}", compareBranch: "${compareBranch}"){
                date
                commits{
                  hash
                  author
                  commitMessage
                }
            }
          }
        `,
      },
    })
      .then((res) => {
        setLoading(false);
        if (res.data.data) {
          const { branchCompare } = res.data.data;
          setCommitLogs([...branchCompare]);
        }
      })
      .catch((err) => {
        console.log(err);
        setLoading(false);
      });
  }, [repoId, baseBranch, compareBranch]);

  return (
    <div className="mx-auto mt-6 mb-10 p-6 block w-11/12 overflow-auto">
      {!loading &&
        commitLogs &&
        commitLogs.map((commit) => {
          return (
            <div
              className="my-4 border-b p-3 border border-gray-200 shadow-md rounded"
              key={commit.date}
            >
              <div className="flex items-center gap-10 text-xl font-sans font-semibold text-gray-800 border-b">
                <div>
                  <FontAwesomeIcon
                    icon={["fas", "calendar-day"]}
                    className="text-gray-700"
                  ></FontAwesomeIcon>
                </div>
                <div>Committed on - {format(new Date(commit.date), "d MMM yyyy")}</div>
              </div>
              <div>
                {commit.commits.map((item) => {
                  return (
                    <div
                      className="flex p-3 justify-between items-center mx-auto border-b w-full"
                      key={item.hash}
                    >
                      <div className="block p-2 font-sans font-light text-gray-800">
                        <div>{item.commitMessage}</div>
                        <div className="flex items-center gap-4 my-2 align-middle">
                          <div>
                            <FontAwesomeIcon
                              icon={["fas", "user-alt"]}
                              className="text-indigo-500"
                            ></FontAwesomeIcon>
                          </div>
                          <div className="text-xl font-semibold font-sans">
                            {item.author}
                          </div>
                        </div>
                      </div>
                      <div className="shadow border rounded p-2 bg-indigo-100 font-mono font-semibold text-indigo-800">
                        #{<>{item.hash ? item.hash.substring(0, 7) : null}</>}
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          );
        })}
      {compareBranch === baseBranch ? (
        <div className="p-2 rounded shadow text-center font-sans text-2xl font-light bg-gray-100">
          Same branches cannot be compared
        </div>
      ) : !loading && commitLogs.length === 0 ? (
        <div className="w-full mx-auto my-auto p-6 shadow border rounded-lg font-sans text-center text-lg">
          No changes found after comparing
          <span className="mx-2 font-sans font-semibold text-blue-600">
            {baseBranch}
          </span>
          with
          <span className="mx-2 font-sans font-semibold text-geay-800">
            {compareBranch}
          </span>
          <div className="my-2 text-xm text-red-400 font-sans font-semibold border-b border-dashed">
            If this is not right, try the other way around
          </div>
        </div>
      ) : null}
    </div>
  );
}
