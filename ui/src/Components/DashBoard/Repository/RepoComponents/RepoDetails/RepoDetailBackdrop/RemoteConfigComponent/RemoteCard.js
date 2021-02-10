import React, { useState, useRef } from "react";
import axios from "axios";
import { globalAPIEndpoint } from "../../../../../../../util/env_config";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPencilAlt,
  faTrashAlt,
  faSave,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import {
  faGithub,
  faGitlab,
  faBitbucket,
  faAws,
  faGitSquare,
} from "@fortawesome/free-brands-svg-icons";

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
  const remoteFormName = useRef();
  const remoteFormUrl = useRef();

  const [remoteNameState, setRemoteNameState] = useState(remoteName);
  const [remoteUrlState, setRemoteUrlState] = useState(remoteUrl);
  const [editRemote, setEditRemote] = useState(false);
  const [deleteRemote, setDeleteRemote] = useState(false);

  var globalUrl = remoteUrl;

  const changeState = (name, url) => {
    let status = "success";
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            mutation {
              editRemote(repoId: "${repoId}", remoteName: "${name}", remoteUrl: ${url}"){
                status
              }
            }
        `,
      },
    })
      .then((res) => {
        status = res.data.data;
        setStatusCheck(false);
        setRemoteOperation(" ");

        if (status === "success") {
          setReloadView(true);
        } else {
          setAddRemoteStatus(true); //status === "failed"
        }
      })
      .catch(() => {
        setStatusCheck(true);
        setRemoteOperation("edit");

        // remoteDetails.forEach((items) => {
        //   if (items.name === name) {
        //     items.name = name;
        //   }
        // });
        // setRemoteDetails([...remoteDetails]);
        // setReloadView(true);
      });

    setRemoteNameState(name);
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
            <input
              type="text"
              autoComplete="off"
              className={`rounded w-full shadow-md py-2 border-2 text-center xl:text-lg lg:text-lg md:text-base text-base items-center text-gray-800 bg-white`}
              style={{ borderColor: "rgb(113 166 196 / 33%)" }}
              placeholder={remoteNameState}
              ref={remoteFormName}
              onChange={(event) => {
                const remoteNameVal = event.target.value;
                if (remoteNameVal.match(/[\s\\//*]/gi)) {
                  event.target.value = remoteNameVal.replace(
                    /[\s\\//*]/gi,
                    "-"
                  );
                }
                setAddRemoteStatus(false);
                setFieldMissing(false);
                setInvalidUrl(false);
              }}
            ></input>
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
                let name;
                let url = !remoteFormUrl.current.value
                  ? remoteUrlState.trim()
                  : remoteFormUrl.current.value.trim();
                if (url.match(/(\s)/g) || url.length === 0) {
                  setInvalidUrl(true);
                } else {
                  if (
                    !remoteFormName.current.value ||
                    remoteFormName.current.value === remoteNameState
                  ) {
                    name = remoteNameState.trim();
                  } else {
                    name = remoteFormName.current.value.trim();
                  }
                  changeState(name, url);
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
                <div className="w-1/2">{remoteNameState}</div>
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
                    let status = "success";
                    axios({
                      url: globalAPIEndpoint,
                      method: "POST",
                      data: {
                        query: `
                                mutation {
                                  deleteRemote(repoId: "${repoId}", remoteName: "${remoteNameState}"){
                                    status
                                  }
                                }
                            `,
                      },
                    })
                      .then((res) => {
                        status = res.data.data;

                        setStatusCheck(false);
                        setRemoteOperation(" ");

                        if (status === "success") {
                          setReloadView(true);
                          setDeleteFailed(false);
                          setDeleteRemote(true);
                        } else {
                          setDeleteFailed(true); //status === "failed"
                        }
                      })
                      .catch(() => {
                        setStatusCheck(true);
                        setRemoteOperation("delete");

                        // setRemoteDetails([
                        //   remoteDetails.filter((items) => {
                        //     return items.name !== remoteNameState;
                        //   }),
                        // ]);

                        // setReloadView(true);
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
