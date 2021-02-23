import {
  faAws,
  faBitbucket,
  faGithub,
  faGitlab,
  faGitSquare,
} from "@fortawesome/free-brands-svg-icons";
import {
  faPencilAlt,
  faSave,
  faTimes,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useRef, useState } from "react";
import { globalAPIEndpoint } from "../../../../../../../util/env_config";

export default function RemoteCard(props) {
  const {
    remoteName,
    remoteUrl,
    setFieldMissing,
    setInvalidUrl,
    setAddRemoteStatus,
    setDeleteFailed,
    setReloadView,
    repoId,
    setStatusCheck,
    setRemoteOperation,
  } = props;

  const remoteFormUrl = useRef();

  const [remoteUrlState, setRemoteUrlState] = useState(remoteUrl);
  const [editRemote, setEditRemote] = useState(false);
  const [deleteRemote, setDeleteRemote] = useState(false);

  var globalUrl = remoteUrl;

  const changeState = (remoteName, url) => {
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            mutation {
              editRemote(repoId: "${repoId}", remoteName: "${remoteName}", remoteUrl: "${url}"){
                status
              }
            }
        `,
      },
    })
      .then((res) => {
        const { status } = res.data.data.editRemote;
        setStatusCheck(false);
        setRemoteOperation(" ");

        if (status === "REMOTE_EDIT_SUCCESS") {
          setReloadView(true);
        } else {
          setAddRemoteStatus(true);
        }
      })
      .catch(() => {
        setStatusCheck(true);
        setRemoteOperation("edit");
      });

    setRemoteUrlState(url);
    setEditRemote(false);
    setFieldMissing(false);
    setInvalidUrl(false);
    setAddRemoteStatus(false);
  };
  const getRemoteLogo = (gitRemoteHost) => {
    let remoteLogo = "";
    if (gitRemoteHost.match(/github/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={faGithub}
          className="w-2/12 mr-2 text-xl text-pink-500 xl:text-3xl lg:text-3xl md:text-2xl"
        ></FontAwesomeIcon>
      );
    } else if (gitRemoteHost.match(/gitlab/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={faGitlab}
          className="w-2/12 mr-2 text-xl text-pink-500 xl:text-3xl lg:text-3xl md:text-2xl"
        ></FontAwesomeIcon>
      );
    } else if (gitRemoteHost.match(/bitbucket/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={faBitbucket}
          className="w-2/12 mr-2 text-xl text-pink-500 xl:text-3xl lg:text-3xl md:text-2xl"
        ></FontAwesomeIcon>
      );
    } else if (gitRemoteHost.match(/codecommit/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={faAws}
          className="w-2/12 mr-2 text-xl text-pink-500 xl:text-3xl lg:text-3xl md:text-2xl"
        ></FontAwesomeIcon>
      );
    } else {
      remoteLogo = (
        <FontAwesomeIcon
          icon={faGitSquare}
          className="w-2/12 mr-2 text-xl text-pink-500 xl:text-3xl lg:text-3xl md:text-2xl"
        ></FontAwesomeIcon>
      );
    }

    return remoteLogo;
  };

  const remoteUrlHandler = (gitRemoteData) => {
    let remoteData = "";
    if (gitRemoteData) {
      if (gitRemoteData.match(/(^https)/gi)) {
        remoteData = (
          <a
            href={gitRemoteData}
            className="text-blue-400 hover:text-blue-500"
            target="_blank"
            rel="noopener noreferrer"
          >
            {gitRemoteData}
          </a>
        );
      } else {
        remoteData = <>{gitRemoteData}</>;
      }
    } else {
      remoteData = " ";
    }
    return remoteData;
  };

  return (
    <div className="w-full">
      {editRemote ? (
        <div className="flex items-center w-full py-6 mx-auto my-1 align-middle rounded-md shadow bg-gray-50">
          <div className="flex items-center justify-center w-1/5 mx-auto text-base text-gray-700 text-sans xl:text-lg lg:text-lg md:text-base">
            {getRemoteLogo(remoteUrlState)}
            <div className="w-1/2">{remoteName}</div>
          </div>
          <div className="flex items-center justify-center w-1/2 mx-auto text-base text-center text-gray-700 text-sans xl:text-lg lg:text-lg md:text-base">
            <input
              type="text"
              autoComplete="off"
              className={`rounded shadow-md w-full py-2 border-2 text-center xl:text-lg lg:text-lg md:text-base text-base items-center text-gray-800 bg-white`}
              style={{ borderColor: "rgb(113 166 196 / 33%)" }}
              placeholder={remoteUrlState}
              value={remoteUrlState}
              ref={remoteFormUrl}
              onChange={(event) => {
                setRemoteUrlState(event.target.value);
                setAddRemoteStatus(false);
                setFieldMissing(false);
                setInvalidUrl(false);
              }}
            ></input>
          </div>
          <div
            className="flex items-center text-center"
            style={{ width: "22%" }}
          >
            <div
              className="items-center w-5/12 p-1 py-2 mx-auto text-base font-semibold bg-blue-500 rounded cursor-pointer xl:text-lg lg:text-lg md:text-base hover:bg-blue-700"
              onClick={() => {
                let url = !remoteFormUrl.current.value
                  ? remoteUrlState.trim()
                  : remoteFormUrl.current.value.trim();
                if (url.match(/(\s)/g) || url.length === 0) {
                  setInvalidUrl(true);
                } else {
                  changeState(remoteName.trim(), url);
                }
              }}
            >
              <FontAwesomeIcon
                icon={faSave}
                className="text-white"
              ></FontAwesomeIcon>
            </div>
            <div
              className="items-center w-5/12 p-1 py-2 mx-auto text-base font-semibold bg-gray-500 rounded cursor-pointer xl:text-lg lg:text-lg md:text-base hover:bg-gray-700"
              onClick={() => {
                setRemoteUrlState(globalUrl);
                setEditRemote(false);
                setAddRemoteStatus(false);
                setFieldMissing(false);
                setInvalidUrl(false);
              }}
            >
              <FontAwesomeIcon
                icon={faTimes}
                className="text-white"
              ></FontAwesomeIcon>
            </div>
          </div>
        </div>
      ) : (
        <>
          {deleteRemote ? (
            " "
          ) : (
            <div className="flex items-center w-full py-6 mx-auto my-1 align-middle rounded-md shadow bg-gray-50">
              <div className="flex items-center justify-center w-1/4 mx-auto text-base text-gray-700 text-sans xl:text-lg lg:text-lg md:text-base">
                {getRemoteLogo(remoteUrlState)}
                <div className="w-1/2">{remoteName}</div>
              </div>
              <div className="flex items-center justify-center w-7/12 mx-auto text-base text-center text-gray-700 text-sans xl:text-lg lg:text-lg md:text-base">
                {remoteUrlHandler(remoteUrlState)}
              </div>

              <div
                className="flex items-center text-center"
                style={{ width: "22%" }}
              >
                <div
                  className="items-center w-5/12 p-1 py-2 mx-auto text-base font-semibold bg-blue-500 rounded cursor-pointer xl:text-lg lg:text-lg md:text-base hover:bg-blue-700"
                  onClick={() => {
                    setEditRemote(true);
                    setAddRemoteStatus(false);
                    setFieldMissing(false);
                    setInvalidUrl(false);
                  }}
                >
                  <FontAwesomeIcon
                    icon={faPencilAlt}
                    className="text-white"
                  ></FontAwesomeIcon>
                </div>
                <div
                  className="items-center w-5/12 p-1 py-2 mx-auto text-base font-semibold bg-red-500 rounded cursor-pointer xl:text-lg lg:text-lg md:text-base hover:bg-red-600"
                  onClick={() => {
                    axios({
                      url: globalAPIEndpoint,
                      method: "POST",
                      data: {
                        query: `
                                mutation {
                                  deleteRemote(repoId: "${repoId}", remoteName: "${remoteName}"){
                                    status
                                  }
                                }
                            `,
                      },
                    })
                      .then((res) => {
                        const { status } = res.data.data.deleteRemote;

                        setStatusCheck(false);
                        setRemoteOperation(" ");

                        if (status === "REMOTE_DELETE_SUCCESS") {
                          setReloadView(true);
                          setDeleteFailed(false);
                          setDeleteRemote(true);
                        } else {
                          setDeleteFailed(true);
                        }
                      })
                      .catch(() => {
                        setStatusCheck(true);
                        setRemoteOperation("delete");
                      });
                  }}
                >
                  <FontAwesomeIcon
                    icon={faTrashAlt}
                    className="text-white"
                  ></FontAwesomeIcon>
                </div>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}
