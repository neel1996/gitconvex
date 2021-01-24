import React, { useState, useRef } from "react";
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
  } = props;
  const remoteFormName = useRef();
  const remoteFormUrl = useRef();

  const [remoteNameState, setRemoteNameState] = useState(remoteName);
  const [remoteUrlState, setRemoteUrlState] = useState(remoteUrl);
  const [editRemote, setEditRemote] = useState(false);
  const [deleteRemote, setDeleteRemote] = useState(false);

  var globalUrl = remoteUrl;

  const changeState = (name, url) => {
    //TODO:Axios

    let status = "success";

    if (status === "success") {
      localStorage.removeItem(remoteNameState);
      localStorage.setItem(
        name,
        JSON.stringify({
          remoteName: name,
          remoteUrl: url,
        })
      );
      setReloadView(true);
    } else {
      setAddRemoteStatus(true);
    }

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
          className="xl:text-3xl lg:text-3xl md:text-2xl text-xl text-pink-500 w-2/12 mr-2"
        ></FontAwesomeIcon>
      );
    } else if (gitRemoteHost.match(/gitlab/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={faGitlab}
          className="xl:text-3xl lg:text-3xl md:text-2xl text-xl text-pink-500 w-2/12 mr-2"
        ></FontAwesomeIcon>
      );
    } else if (gitRemoteHost.match(/bitbucket/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={faBitbucket}
          className="xl:text-3xl lg:text-3xl md:text-2xl text-xl text-pink-500 w-2/12 mr-2"
        ></FontAwesomeIcon>
      );
    } else if (gitRemoteHost.match(/codecommit/i)) {
      remoteLogo = (
        <FontAwesomeIcon
          icon={faAws}
          className="xl:text-3xl lg:text-3xl md:text-2xl text-xl text-pink-500 w-2/12 mr-2"
        ></FontAwesomeIcon>
      );
    } else {
      remoteLogo = (
        <FontAwesomeIcon
          icon={faGitSquare}
          className="xl:text-3xl lg:text-3xl md:text-2xl text-xl text-pink-500 w-2/12 mr-2"
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
        <div className="flex items-center align-middle w-full mx-auto my-1 shadow rounded-md py-6 bg-gray-50">
          <div className="flex items-center w-1/5 mx-auto justify-center text-sans xl:text-lg lg:text-lg md:text-base text-base text-gray-700">
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
          <div className="text-sans mx-auto justify-center items-center text-center flex xl:text-lg lg:text-lg md:text-base text-base text-gray-700 w-1/2">
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
            className="text-center flex items-center"
            style={{ width: "22%" }}
          >
            <div
              className="xl:text-lg lg:text-lg md:text-base text-base items-center p-1 py-2 rounded w-5/12 mx-auto cursor-pointer bg-blue-500 hover:bg-blue-700 font-semibold"
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
              className="xl:text-lg lg:text-lg md:text-base text-base items-center p-1 py-2 rounded w-5/12 mx-auto cursor-pointer bg-gray-500 hover:bg-gray-700 font-semibold"
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
            <div className="flex items-center align-middle w-full mx-auto my-1 shadow rounded-md py-6 bg-gray-50">
              <div className="flex items-center w-1/4 mx-auto justify-center text-sans xl:text-lg lg:text-lg md:text-base text-base text-gray-700">
                {getRemoteLogo(remoteUrlState)}
                <div className="w-1/2">{remoteNameState}</div>
              </div>
              <div className="text-sans mx-auto justify-center items-center text-center flex xl:text-lg lg:text-lg md:text-base text-base text-gray-700 w-7/12">
                {remoteUrlHandler(remoteUrlState)}
              </div>

              <div
                className="text-center flex items-center"
                style={{ width: "22%" }}
              >
                <div
                  className="xl:text-lg lg:text-lg md:text-base text-base items-center p-1 py-2 rounded w-5/12 mx-auto cursor-pointer bg-blue-500 hover:bg-blue-700 font-semibold"
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
                  className="xl:text-lg lg:text-lg md:text-base text-base items-center p-1 py-2 rounded w-5/12 mx-auto cursor-pointer bg-red-500 hover:bg-red-600 font-semibold"
                  onClick={() => {
                    //TODO: Add axios
                    let status = "success";
                    if (status === "success") {
                      localStorage.removeItem(remoteNameState);
                      setReloadView(true);
                      setDeleteFailed(false);
                      setDeleteRemote(true);
                    } else {
                      setDeleteFailed(true);
                    }
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
