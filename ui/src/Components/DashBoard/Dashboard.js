import axios from "axios";
import React, { useContext, useEffect, useMemo, useState } from "react";
import { BrowserRouter, Route } from "react-router-dom";
import { ContextProvider } from "../../context";
import { globalAPIEndpoint } from "../../util/env_config";
import LeftPane from "./DashboardPaneComponents/LeftPane";
import RightPane from "./DashboardPaneComponents/RightPane";
import Help from "./Help/Help";
import RepositoryAction from "./Repository/RepoComponents/RepositoryAction";
import RepositoryDetails from "./Repository/RepoComponents/RepoDetails/RepositoryDetails";
import Settings from "./Settings/Settings";
import CompareComponent from "./Compare/CompareComponent";

export default function Dashboard(props) {
  const { state } = useContext(ContextProvider);

  const [platform, setPlatform] = useState("");
  const [gitVersion, setGitVersion] = useState("");

  const memoizedRepoDetails = useMemo(() => {
    return <RepositoryDetails parentProps={props}></RepositoryDetails>;
  }, [props]);

  const memoizedRepoAction = useMemo(() => {
    return <RepositoryAction></RepositoryAction>;
  }, []);

  const memoizedCompare = useMemo(() => {
    return <CompareComponent></CompareComponent>;
  }, []);

  const memoizedSettings = useMemo(() => {
    return <Settings></Settings>;
  }, []);

  const memoizedHelp = useMemo(() => {
    return <Help></Help>;
  }, []);

  useEffect(() => {
    const { osCheck, gitCheck } = state.hcParams;

    const localStorageItems = ["OS_TYPE", "GIT_VERSION"];

    const token = axios.CancelToken;
    const source = token.source();

    if (osCheck && gitCheck) {
      setPlatform(osCheck);
      setGitVersion(gitCheck);
    } else {
      let checkArray = localStorageItems.filter((item) => {
        return localStorage.getItem(item) ? true : false;
      });

      if (checkArray.length === 2) {
        setPlatform(localStorage.getItem("OS_TYPE"));
        setGitVersion(localStorage.getItem("GIT_VERSION"));
      } else {
        axios({
          url: globalAPIEndpoint,
          method: "POST",
          cancelToken: source.token,
          data: {
            query: `
              query{
                healthCheck{
                  os
                  gitconvex
                }
              }
            `,
          },
        })
          .then((res) => {
            if (res.data.data && !res.data.error) {
              const { os, gitconvex } = res.data.data.gitConvexApi.healthCheck;

              setPlatform(JSON.parse(os).message);
              setGitVersion(JSON.parse(gitconvex).message);
            }
          })
          .catch((err) => {
            console.log(err);
          });
      }
    }
    return () => {
      source.cancel();
    };
  }, [state.hcParams]);

  const params = {
    platform,
    gitVersion,
  };

  const renderRightPaneComponent = () => {
    switch (props.location.pathname) {
      case "/dashboard":
        return <RightPane params={params}></RightPane>;
      case "/dashboard/repository":
        return memoizedRepoAction;
      case "/dashboard/compare":
        return memoizedCompare;
      case "/dashboard/settings":
        return memoizedSettings;
      case "/dashboard/help":
        return memoizedHelp;
      default:
        return (
          <BrowserRouter>
            <Route exact path="/dashboard/repository/:repoId">
              {memoizedRepoDetails}
            </Route>
          </BrowserRouter>
        );
    }
  };

  return (
    <>
      <div className="flex w-full h-screen">
        <LeftPane parentProps={props}></LeftPane>
        {renderRightPaneComponent()}
      </div>
    </>
  );
}
