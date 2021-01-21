import { library } from "@fortawesome/fontawesome-svg-core";
import { fab } from "@fortawesome/free-brands-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useState } from "react";
import { NavLink } from "react-router-dom";
import { globalAPIEndpoint } from "../../../../util/env_config";
import InfiniteLoader from "../../../Animations/InfiniteLoader";

export default function RepoCard(props) {
  library.add(fab, fas);
  const { repoData } = props;

  const [repoFooterData, setRepoFooterData] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    let repoId = props.repoData.id;
    const token = axios.CancelToken;
    const source = token.source();

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      cancelToken: source.token,
      data: {
        query: `
          query 
          {
            gitRepoStatus(repoId:"${repoId}"){
              gitCurrentBranch
              gitTotalCommits
              gitTotalTrackedFiles
            }
          }
        `,
      },
    })
      .then((res) => {
        setLoading(false);
        setRepoFooterData(res.data.data.gitRepoStatus);
      })
      .catch((err) => {
        setLoading(false);
      });

    return () => {
      source.cancel();
    };
  }, [props.repoData.id]);

  const repoName = repoData.repoName;
  var avatar = "";

  if (repoName) {
    if (repoName.split(" ").length > 1) {
      let tempName = repoName.split(" ");
      avatar = tempName[0].substring(0, 1) + tempName[1].substring(0, 1);
      avatar = avatar.toUpperCase();
    } else {
      avatar = repoName.substring(0, 1).toUpperCase();
    }
  }

  return (
    <NavLink
      to={`/dashboard/repository/${repoData.id}`}
      className="xl:w-96 lg:w-2/4 md:w-1/2 bg-indigo-400 border-gray-300 rounded-lg border cursor-pointer block my-6 p-6 shadow-md text-center hover:shadow-xl"
      key={repoData.repoName}
    >
      <div className="bg-indigo-300 rounded text-5xl my-2 py-5 px-10 shadow text-center text-white">
        {avatar}
      </div>
      <div className="border-indigo-300 border-dashed border-b-2 font-sans text-2xl my-4 pb-2 text-white">
        {repoData.repoName}
      </div>
      <div className="rounded-md flex justify-center my-2 mx-auto shadow-sm text-center align-middle w-full">
        {loading || !repoFooterData ? (
          <div className="block mx-auto w-full bg-white rounded">
            <div className="flex mx-auto my-6 text-center justify-center">
              <InfiniteLoader
                loadAnimation={loading || !repoFooterData}
              ></InfiniteLoader>
            </div>
          </div>
        ) : (
          <>
            <div className="bg-white border-indigo-300 flex items-center my-2 p-2 shadow-lg w-1/2 rounded-l-md">
              <FontAwesomeIcon
                className="my-auto"
                icon={["fas", "grip-lines"]}
              ></FontAwesomeIcon>
              <div className="items-center font-sans text-sm mx-2 text-center">
                {repoFooterData && repoFooterData.gitTotalCommits ? (
                  <>{repoFooterData.gitTotalCommits} Commits</>
                ) : (
                  <>0 Commits</>
                )}
              </div>
            </div>
            <div className="bg-white border-indigo-300 flex items-center my-2 p-2 shadow-lg w-1/2">
              <FontAwesomeIcon
                className="my-auto"
                icon={["fas", "file-alt"]}
              ></FontAwesomeIcon>
              <div className="items-center font-sans text-sm mx-2 text-center">
                {repoFooterData && repoFooterData.gitTotalTrackedFiles ? (
                  <>{repoFooterData.gitTotalTrackedFiles} Tracked Files</>
                ) : (
                  <>0 Tracked Files</>
                )}
              </div>
            </div>
            <div className="bg-white border-indigo-300 flex items-center my-2 p-2 shadow-lg w-1/2 rounded-r-md">
              <FontAwesomeIcon
                className="my-auto"
                icon={["fas", "code-branch"]}
              ></FontAwesomeIcon>
              <div className="items-center font-sans text-sm mx-2 text-center font-semibold">
                {repoFooterData && repoFooterData.gitCurrentBranch ? (
                  <>{repoFooterData.gitCurrentBranch}</>
                ) : (
                  <>No Branches Available</>
                )}
              </div>
            </div>
          </>
        )}
      </div>
    </NavLink>
  );
}
