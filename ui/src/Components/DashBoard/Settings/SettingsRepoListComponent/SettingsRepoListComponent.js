import axios from "axios";
import { format } from "date-fns";
import React, { useContext, useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../util/env_config";
import { SettingsContext } from "../Settings";
import SettingsRepoListCard from "./SettingsRepoListCard";

export default function SettingsRepoListComponent() {
  const [repoDetails, setRepoDetails] = useState([]);
  const [loading, setLoading] = useState(true);
  const { state } = useContext(SettingsContext);

  useEffect(() => {
    const token = axios.CancelToken;
    const source = token.source();
    setLoading(true);

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
        setLoading(false);
        if (res.data.data && !res.data.error) {
          const repoDetails = res.data.data.fetchRepo;
          setRepoDetails(repoDetails);
        }
      })
      .catch((err) => {
        setLoading(false);
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
        {!loading &&
        repoDetails &&
        repoDetails.repoId &&
        repoDetails.repoId.length ? (
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
          <div className="my-10 mx-auto font-sans font-medium text-gray-700 bg-gray-200 text-center p-10 rounded shadow w-3/4">
            {loading
              ? "Fetching repo list from the server..."
              : "No repos are being managed by Gitconvex. You can add one from the dashboard"}
          </div>
        )}
      </>
    </div>
  );
}
