import axios from "axios";
import { format } from "date-fns";
import React, { useContext, useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../util/env_config";
import { SettingsContext } from "../Settings";
import SettingsRepoListCard from "./SettingsRepoListCard";

export default function SettingsRepoListComponent() {
  const [repoDetails, setRepoDetails] = useState([]);
  const { state } = useContext(SettingsContext);

  useEffect(() => {
    const token = axios.CancelToken;
    const source = token.source();

    axios({
      url: globalAPIEndpoint,
      cancelToken: source.token,
      method: "POST",
      data: {
        query: `
            query {
                fetchRepo{
                    repoId
                    repoName
                    repoPath
                    timeStamp
                }
            }
        `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          const repoDetails = res.data.data.fetchRepo;
          setRepoDetails(repoDetails);
        }
      })
      .catch((err) => {
        console.log(err);
      });

    return () => {
      return source.cancel;
    };
  }, [state]);

  return (
    <div className="repo-data my-10">
      <div className="text-xl text-gray-700 font-sans font-semibold">
        Saved Repositories
      </div>
      <>
        {repoDetails && repoDetails.repoId && repoDetails.repoId.length ? (
          <>
            <div className="flex my-4 bg-indigo-500 w-full rounded text-white shadow p-3 font-sand text-xl font-semibold">
              <div className="w-1/4 border-r text-center">Repo ID</div>
              <div className="w-1/2 border-r text-center">Repo Name</div>
              <div className="w-1/2 border-r text-center">Repo Path</div>
              <div className="w-1/2 border-r text-center">Timestamp</div>
              <div className="w-1/2 border-r text-center">Action</div>
            </div>
            {repoDetails.repoId.map((repoId, idx) => {
              return (
                <SettingsRepoListCard
                  key={repoId}
                  repoId={repoId}
                  repoName={repoDetails.repoName[idx]}
                  repoPath={repoDetails.repoPath[idx]}
                  timeStamp={format(
                    new Date(repoDetails.timeStamp[idx]),
                    "MMMM dd, yyyy"
                  )}
                ></SettingsRepoListCard>
              );
            })}
          </>
        ) : (
          <div className="my-10 mx-auto bg-gray-200 text-center p-10 rounded shadow w-3/4">
            No repos are being managed by Gitconvex. You can add one from the
            dashboard
          </div>
        )}
      </>
    </div>
  );
}
