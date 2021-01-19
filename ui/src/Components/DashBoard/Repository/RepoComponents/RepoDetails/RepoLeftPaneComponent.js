import { library } from "@fortawesome/fontawesome-svg-core";
import { fab } from "@fortawesome/free-brands-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { actionType } from "./backdropActionType";
import React, { useEffect } from "react";

export default function RepoLeftPaneComponent(props) {
  library.add(fab, fas);

  let {
    gitRemoteHost,
    gitRemoteData,
    isMultiRemote,
    multiRemoteCount,
    showCommitLogsView,
    actionTrigger,
  } = props;

  useEffect(() => {}, [props]);

  const getRemoteLogo = () => {
    let remoteLogo = "";
    if (gitRemoteHost.match(/github/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={["fab", "github"]}
          className="text-4xl text-center text-pink-500"
        ></FontAwesomeIcon>
      );
    } else if (gitRemoteHost.match(/gitlab/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={["fab", "gitlab"]}
          className="text-4xl text-center text-pink-400"
        ></FontAwesomeIcon>
      );
    } else if (gitRemoteHost.match(/bitbucket/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={["fab", "bitbucket"]}
          className="text-4xl text-center text-pink-400"
        ></FontAwesomeIcon>
      );
    } else if (gitRemoteHost.match(/codecommit/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={["fab", "aws"]}
          className="text-4xl text-center text-pink-400"
        ></FontAwesomeIcon>
      );
    } else {
      remoteLogo = (
        <FontAwesomeIcon
          icon={["fab", "git-square"]}
          className="text-4xl text-center text-pink-400"
        ></FontAwesomeIcon>
      );
    }

    return remoteLogo;
  };

  return (
    <>
      {props.received ? (
        <div className="border-gray-300 rounded-md border-dotted border-2 block my-6 mx-auto p-1 shadow-sm w-11/12 xl:w-1/2 lg:w-3/4 md:w-11/12 sm:w-11/12">
          <div className="block mx-auto my-6">
            <div className="flex items-center justify-evenly">
              <div className="text-lg text-gray-600 w-1/4">Remote Host</div>
              <div className="flex justify-around items-center align-middle w-1/2">
                <div className="w-3/4 shadow rounded-md border-dashed border cursor-pointer flex items-center justify-center my-auto p-4 align-middle">
                  {gitRemoteHost ? (
                    <div className="mx-2">{getRemoteLogo()}</div>
                  ) : null}
                  <div
                    className={`${
                      gitRemoteHost !== "No Remote Host Available"
                        ? "text-xl border-gray-300 border-dashed border-b text-center text-gray-800 w-3/4"
                        : "text-base font-sans font-light text-gray-600 text-center"
                    }`}
                  >
                    {gitRemoteHost}
                  </div>
                </div>
                <div className="w-1/4">
                  <div
                    id="addRemote"
                    className="rounded-full cursor-pointer items-center h-10 text-2xl mx-auto shadow text-center text-white align-middle w-10 bg-indigo-400 hover:bg-indigo-500"
                    onMouseEnter={(event) => {
                      let popUp =
                        '<div class="tooltip" style="margin-left:-40px;">Click to add a new remote repo</div>';
                      event.target.innerHTML += popUp;
                    }}
                    onMouseLeave={(event) => {
                      event.target.innerHTML = "+";
                    }}
                    onClick={() => {
                      actionTrigger(actionType.ADD_REMOTE_REPO);
                    }}
                  >
                    +
                  </div>
                </div>
              </div>
            </div>

            <div className="remote  flex justify-evenly my-4">
              <div className="text-lg text-gray-600 w-1/4">
                {`${gitRemoteHost} URL`}
              </div>
              <div className="cursor-pointer text-blue-400 break-words w-1/2 hover:text-blue-500">
                {gitRemoteData}
              </div>
            </div>

            {isMultiRemote ? (
              <div className=" flex justify-evenly my-2">
                <div className="font-sans text-gray-800 font-semibold w-1/4 border-dotted border-b-2 border-gray-200">
                  Entry truncated!
                </div>
                <div className="w-1/2 border-dotted border-b-2 border-gray-200">
                  {`Remote repos : ${multiRemoteCount}`}
                </div>
              </div>
            ) : null}
          </div>

          <div className="block my-6 mx-auto">
            <div
              className="p-3 text-gray-600 text-center w-3/4 mx-auto rounded-md shadow-md bg-yellow-200 text-whtie font-sans font-semibold text-xl hover:bg-yellow-100 hover:shadow-sm cursor-pointer transition"
              onClick={(event) => {
                showCommitLogsView();
              }}
            >
              SHOW COMMIT LOGS
            </div>
          </div>
        </div>
      ) : null}
    </>
  );
}
