import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../util/env_config";

export default function SearchRepoCards(props) {
  const [repo, setRepo] = useState([]);

  useEffect(() => {
    const token = axios.CancelToken;
    const source = token.source();

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      cancelToken: source.token,
      data: {
        query: `
          query {
              fetchRepo{
                repoId
                repoName
                repoPath
              }
          }
        `,
      },
    })
      .then((res) => {
        const apiResponse = res.data.data.fetchRepo;

        if (apiResponse) {
          const { repoId, repoName, repoPath } = apiResponse;
          let repoContent = [];

          repoId.forEach((entry, index) => {
            if (
              repoName[index]
                .toLowerCase()
                .match(props.searchQuery.toLowerCase())
            ) {
              repoContent.push({
                id: entry,
                repoName: repoName[index],
                repoPath: repoPath[index],
              });
            }
          });

          setRepo(repoContent);
        }
      })
      .catch((err) => {});

    return () => {
      return source.cancel;
    };
  }, [props.searchQuery]);

  return (
    <div className="w-full">
      {repo ? (
        repo.map((item) => {
          return (
            <div
              className="border-b cursor-pointer flex items-center justify-around my-4 mx-auto p-4 hover:bg-gray-100"
              key={item.id}
              onClick={(e) => {
                props.setSelectedRepoHandler(item);
              }}
            >
              <div className="w-1/2 text-2xl font-sand font-semibold">
                {item.repoName}
              </div>
              <div className="w-1/4 block justify-center">
                <div className="bg-blue-100 border shadow rounded p-2">
                  PATH
                </div>
                <div
                  className="my-2 font-light break-words text-gray-700"
                  title={item.repoPath}
                >
                  {item.repoPath}
                </div>
              </div>
            </div>
          );
        })
      ) : (
        <div className="text-center text-2xl font-sans font-light">
          Loading...
        </div>
      )}
      {!repo[0] ? (
        <div className="text-3xl my-6 font-sans font-light text-gray-600">
          There are no matching repos...
        </div>
      ) : null}
    </div>
  );
}
