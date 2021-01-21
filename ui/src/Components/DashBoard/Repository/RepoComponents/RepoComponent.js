import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useContext, useEffect, useState } from "react";
import {
  ADD_FORM_CLOSE,
  DELETE_PRESENT_REPO,
  PRESENT_REPO,
} from "../../../../actionStore";
import { ContextProvider } from "../../../../context";
import { globalAPIEndpoint } from "../../../../util/env_config";
import InfiniteLoader from "../../../Animations/InfiniteLoader";
import AddRepoFormComponent from "./AddRepoForm";
import RepoCard from "./RepoCard";

export default function RepoComponent(props) {
  library.add(fas);

  const [repo, setRepo] = useState([]);
  const [repoFormEnable, setRepoFormEnable] = useState(false);
  const [loading, setLoading] = useState(false);

  const { dispatch } = useContext(ContextProvider);

  useEffect(() => {
    setLoading(true);
    const fetchRepoURL = globalAPIEndpoint;
    const token = axios.CancelToken;
    const source = token.source();

    axios({
      url: fetchRepoURL,
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
        setLoading(false);

        if (apiResponse) {
          const { repoId, repoName } = apiResponse;
          let repoContent = [];

          repoId.forEach((entry, index) => {
            repoContent.push({ id: entry, repoName: repoName[index] });
          });

          setRepo(repoContent);

          dispatch({
            type: DELETE_PRESENT_REPO,
            payload: [],
          });

          dispatch({
            action: PRESENT_REPO,
            payload: [...repoContent],
          });
        }
      })
      .catch((err) => {
        console.log(err);
        setLoading(false);
      });

    return () => {
      source.cancel();
    };
  }, [repoFormEnable, dispatch]);

  const showAvailableRepo = () => {
    const repoArray = repo;

    return (
      <>
        <div className="flex flex-wrap justify-center mx-auto pb-40 w-full gap-16">
          <>
            {repoArray.length > 0 ? (
              <>
                {repoArray.map((entry) => {
                  return <RepoCard key={entry.id} repoData={entry}></RepoCard>;
                })}
                {repoArray.length % 2 !== 0 && repoArray.length !== 1 ? (
                  <div className="xl:w-1/3 lg:w-2/4 md:w-1/2 block p-6 my-6 text-center"></div>
                ) : null}
              </>
            ) : (
              <div className="w-3/4 bg-gray-100 rounded-md font-sans text-xl my-10 mx-auto p-10 shadow text-center text-gray-800">
                {loading ? (
                  <div className="block justify-center my-6 mx-auto text-center">
                    <div className="flex justify-center my-6 mx-auto text-center">
                      <InfiniteLoader loadAnimation={loading}></InfiniteLoader>
                    </div>
                    <div>Loading available repos...</div>
                  </div>
                ) : (
                  <div className="font-sans text-2xl text-gray-400 text-center">
                    No repos present. Press + to add a new repo
                  </div>
                )}
              </div>
            )}
          </>
        </div>
        <div className="fixed bottom-0 right-0 mb-10 mr-16 cursor-pointer justify-center">
          <div
            id="addRepoButton"
            className="border-8 border-indigo-100 shadow-lg bg-indigo-300 hover:bg-indigo-400 rounded-full h-20 w-20 flex justify-center text-white font-sans font-black"
            onClick={() => {
              setRepoFormEnable(true);
              dispatch({ type: ADD_FORM_CLOSE, payload: false });
            }}
            onMouseEnter={(event) => {
              event.stopPropagation();
              event.preventDefault();
              document.getElementById("pop-up").classList.remove("hidden");
            }}
            onMouseLeave={(event) => {
              document.getElementById("pop-up").classList.add("hidden");
            }}
          >
            <div className="flex w-full h-full justify-center items-center text-center">
              <div>
                <FontAwesomeIcon
                  icon={["fas", "plus"]}
                  size="2x"
                  className="text-indigo-100"
                ></FontAwesomeIcon>
              </div>
            </div>
            <div
              id="pop-up"
              className="fixed p-2 rounded bg-white w-48 text-gray-600 text-center font-sans font-medium text-sm shadow border mx-auto hidden"
              style={{ marginTop: "-75px", width: "130px" }}
            >
              Click to add a new repo
            </div>
          </div>
        </div>
      </>
    );
  };

  const addFormRemove = (param) => {
    setRepoFormEnable(param);
  };

  return (
    <div className="flex flex-wrap justify-center mx-auto text-center align-middle">
      {!repoFormEnable ? (
        showAvailableRepo()
      ) : (
        <AddRepoFormComponent formEnable={addFormRemove}></AddRepoFormComponent>
      )}
    </div>
  );
}
