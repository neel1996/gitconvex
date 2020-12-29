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
        <div className="repo-leftpane xl:w-1/2 lg:w-3/4 md:w-11/12 sm:w-11/12">
          <div className="block mx-auto my-6">
            <div className="repo-leftpane--remote">
              <div className="text-lg text-gray-600 w-1/4">Remote Host</div>
              <div className="flex justify-around items-center align-middle w-1/2">
                <div className="repo-leftpane--remote--host">
                  {gitRemoteHost ? (
                    <div className="mx-2">{getRemoteLogo()}</div>
                  ) : null}
                  <div
                    className={`${
                      gitRemoteHost !== "No Remote Host Available"
                        ? "text-xl repo-leftpane--remote--name"
                        : "text-base font-sans font-light text-gray-600 text-center"
                    }`}
                  >
                    {gitRemoteHost}
                  </div>
                </div>
                <div className="w-1/4">
                  <div
                    id="addRemote"
                    className="add-btn bg-indigo-400 hover:bg-indigo-500"
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

            <div className="remote flex-even my-4">
              <div className="text-lg text-gray-600 w-1/4">
                {`${gitRemoteHost} URL`}
              </div>
              <div className="remote--url">{gitRemoteData}</div>
            </div>

            {isMultiRemote ? (
              <div className="flex-even my-2">
                <div className="font-sans text-gray-800 font-semibold w-1/4 border-dotted border-b-2 border-gray-200">
                  Entry truncated!
                </div>
                <div className="w-1/2 border-dotted border-b-2 border-gray-200">
                  {`Remote repos : ${multiRemoteCount}`}
                </div>
              </div>
            ) : null}
          </div>

          <div className="commitlogs">
            <div className="flex-even my-3">
              <div className="commitlogs--label">Commit Logs</div>
              <div
                className="commitlogs--content"
                onClick={(event) => {
                  showCommitLogsView();
                }}
              >
                Show Commit Logs
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </>
  );
}
