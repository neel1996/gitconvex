import axios from "axios";
import React, { useEffect, useState } from "react";
import { NavLink } from "react-router-dom";
import {
  globalAPIEndpoint,
  ROUTE_REPO_DETAILS,
} from "../../../../util/env_config";
import InfiniteLoader from "../../../Animations/InfiniteLoader";
import { library } from "@fortawesome/fontawesome-svg-core";
import { fab } from "@fortawesome/free-brands-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import "../../../styles/RepoCard.css";

export default function RepoCard(props) {
  library.add(fab, fas);
  const { repoData } = props;

  const [repoFooterData, setRepoFooterData] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    let repoId = props.repoData.id;
    const payload = JSON.stringify(JSON.stringify({ repoId: repoId }));

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
    
                query GitConvexApi
                {
                  gitConvexApi(route: "${ROUTE_REPO_DETAILS}", payload: ${payload}){
                    gitRepoStatus {
                      gitTotalCommits
                      gitTotalTrackedFiles  
                      gitCurrentBranch  
                    }
                  }
                }
              `,
      },
    })
      .then((res) => {
        setLoading(false);
        setRepoFooterData(res.data.data.gitConvexApi.gitRepoStatus);
      })
      .catch((err) => {
        setLoading(false);
      });

    return () => {
      source.cancel();
    };
  }, [props]);

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
      className="xl:w-1/3 lg:w-2/4 md:w-1/2 repo-card"
      key={repoData.repoName}
    >
      <div className="repo-card--avatar">{avatar}</div>
      <div className="repo-card--name">{repoData.repoName}</div>
      <div className="repo-card--footer">
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
            <div className="footer--items rounded-l-md">
              <FontAwesomeIcon
                className="my-auto"
                icon={["fas", "grip-lines"]}
              ></FontAwesomeIcon>
              <div className="footer--items--label">
                {repoFooterData ? (
                  <>{repoFooterData.gitTotalCommits} Commits</>
                ) : (
                  "..."
                )}
              </div>
            </div>
            <div className="footer--items">
              <FontAwesomeIcon
                className="my-auto"
                icon={["fas", "file-alt"]}
              ></FontAwesomeIcon>
              <div className="footer--items--label">
                {repoFooterData ? (
                  <>{repoFooterData.gitTotalTrackedFiles} Tracked Files</>
                ) : (
                  "..."
                )}
              </div>
            </div>
            <div className="footer--items rounded-r-md">
              <FontAwesomeIcon
                className="my-auto"
                icon={["fas", "code-branch"]}
              ></FontAwesomeIcon>
              <div className="footer--items--label font-semibold">
                {repoFooterData ? (
                  <>{repoFooterData.gitCurrentBranch}</>
                ) : (
                  "..."
                )}
              </div>
            </div>
          </>
        )}
      </div>
    </NavLink>
  );
}
