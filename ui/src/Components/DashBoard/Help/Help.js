import { library } from "@fortawesome/fontawesome-svg-core";
import { fab } from "@fortawesome/free-brands-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useEffect, useState } from "react";
import GoLogo from "../../../assets/Go-Logo_White.svg";
import { CURRENT_VERSION } from "../../../util/env_config";

export default function Help() {
  library.add(fas, fab);

  const [currentVersion, setCurrentVersion] = useState("");
  const [availableUpdate, setAvailableUpdate] = useState("");
  const [showUpdatePane, setShowUpdatePane] = useState(false);
  const [isLatest, setIsLatest] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);

  useEffect(() => {
    setCurrentVersion(CURRENT_VERSION);
  }, []);

  function resetStates() {
    setShowUpdatePane(false);
    setLoading(false);
    setError(false);
    setAvailableUpdate("");
    setIsLatest(false);
  }

  const supportData = [
    {
      label: "reach via discord",
      link: "https://discord.gg/PSd2Cq9",
      icon: ["fab", "discord"],
      color: ["bg-indigo-400", "bg-indigo-500"],
    },
    {
      label: "reach via github",
      link: "https://github.com/neel1996/gitconvex-package/issues",
      icon: ["fab", "github"],
      color: ["bg-gray-800", "bg-gray-700"],
    },
    {
      label: "reach via twitter",
      link: "https://twitter.com/neeldev96",
      icon: ["fab", "twitter"],
      color: ["bg-blue-400", "bg-blue-600"],
    },
  ];

  const contributionData = [
    {
      label: "Gitconvex react project",
      link: "https://github.com/neel1996/gitconvex-ui",
      icon: ["fab", "react"],
      color: ["bg-blue-400"],
      ind: "",
    },
    {
      label: "Gitconvex Go project",
      link: "https://github.com/neel1996/gitconvex-server",
      icon: GoLogo,
      color: ["bg-gray-700"],
      ind: "go",
    },
  ];

  function checkUpdateHandler() {
    resetStates();
    setShowUpdatePane(true);
    setLoading(true);

    const githubEndpoint =
      "https://api.github.com/repos/neel1996/gitconvex/releases/latest";

    axios({
      url: githubEndpoint,
      method: "GET",
      headers: {
        Accept: "application/json",
      },
    })
      .then((res) => {
        const { tag_name } = res.data;
        setLoading(false);

        if (currentVersion === tag_name) {
          setIsLatest(true);
          setAvailableUpdate("");
        } else {
          setAvailableUpdate(tag_name);
        }
      })
      .catch((err) => {
        console.log(err);
      });
  }

  return (
    <div className="w-full h-auto overflow-auto">
      <div className="flex text-5xl text-gray-700 mx-6 my-auto align-middle items-center">
        <FontAwesomeIcon icon={["fas", "question-circle"]}></FontAwesomeIcon>
        <div className="my-5 mx-5 font-sans">Help and Support</div>
      </div>
      <div className="my-4 mx-10">
        <div className="text-2xl font-sans text-gray-900">
          Facing an issue or need any help?
        </div>
        <span className="text-xl text-gray-500 font-medium my-2">
          Before raising an issue, please make sure you have gone through the
        </span>
        <span className="font-mono text-xl mx-2 text-indigo-500 hover:text-indigo-600 hover:font-semibold cursor-pointer">
          <a
            href="https://github.com/neel1996/gitconvex/blob/master/DOCUMENTATION.md"
            target="_blank"
            rel="noopener noreferrer"
          >
            Documentation
          </a>
        </span>
      </div>

      <div className="support-feedback my-10 mx-10">
        <div className="text-2xl font-sans font-semibold">
          Support and Feedback
        </div>

        <div className="my-20 flex justify-center gap-10">
          {supportData.map((data) => {
            return (
              <div key={data.label}>
                <a href={data.link} target="_blank" rel="noopener noreferrer">
                  <div
                    className={`block mx-auto p-6 rounded-lg shadow-md ${data.color[0]} text-white text-center hover:${data.color[1]} hover:shadow-lg`}
                  >
                    <FontAwesomeIcon
                      size="4x"
                      icon={data.icon}
                    ></FontAwesomeIcon>
                    <div className="mx-3 font-sans font-semibold">
                      {data.label}
                    </div>
                  </div>
                </a>
              </div>
            );
          })}
        </div>
      </div>

      <div className="support-feedback my-10 mx-10">
        <div className="text-2xl font-sans font-semibold">Contribution</div>

        <div className="font-sans font-light border-b border-dashed text-center text-gray-600 my-5 break-normal">
          gitconvex is open source and please visit the repo if you are
          interested in contributing to the platform
        </div>
      </div>

      <div className="my-10 flex justify-center gap-10">
        {contributionData.map((data) => {
          return (
            <div key={data.label}>
              <a href={data.link} target="_blank" rel="noopener noreferrer">
                <div
                  className={`block mx-auto p-6 rounded-lg shadow-md ${data.color[0]} text-white text-center hover:shadow-xl`}
                  style={{
                    width: "250px",
                    height: "150px",
                  }}
                >
                  {data.ind !== "go" ? (
                    <FontAwesomeIcon
                      size="4x"
                      icon={data.icon}
                    ></FontAwesomeIcon>
                  ) : (
                    <img
                      src={GoLogo}
                      alt="go-logo"
                      className="text-center mx-auto items-center flex w-20"
                    ></img>
                  )}
                  <div className="mx-3 font-sans font-semibold">
                    {data.label}
                  </div>
                </div>
              </a>
            </div>
          );
        })}
      </div>

      <div className="fixed bottom-0 right-0 mr-6 mb-6 ">
        {showUpdatePane ? (
          <div
            className="bg-white shadow-md block border rounded-lg p-10 text-xl font-sans font-semibold right-0 bottom-0 fixed mr-5"
            style={{
              marginBottom: "120px",
            }}
          >
            <div
              className="relative text-gray-400 text-xl float-right cursor-pointer"
              style={{
                marginTop: "-40px",
                marginRight: "-25px",
                paddingRight: "-15px",
              }}
              onClick={() => {
                resetStates();
              }}
            >
              x
            </div>
            <div className="text-xl font-sans font-semibold text-gray-800 text-center">
              Current Version : {currentVersion}
            </div>
            {loading ? (
              <div className="text-gray-400 font-sans font-light">
                checking for updates...
              </div>
            ) : null}

            {isLatest ? (
              <div className="text-center text-sm font-sans font-light text-gray-600">
                You are using the latest version of gitconvex
              </div>
            ) : null}

            {error ? (
              <div className="text-red-600 font-sans">
                Cannot reach update validation server
              </div>
            ) : null}

            {availableUpdate ? (
              <div className="my-2">
                <div className="text-xl text-gray-800 text-left">
                  New update v{availableUpdate} available.
                </div>
                <div className="p-4 rounded-lg my-2 text-lg font-sans text-center bg-pink-500 text-white">
                  <a
                    href={`https://github.com/neel1996/gitconvex/releases/tag/${availableUpdate}`}
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    Get Update
                  </a>
                </div>
              </div>
            ) : null}
          </div>
        ) : null}

        <div
          className="bg-indigo-400 shadow-lg w-20 h-20 rounded-full text-center cursor-pointer hover:bg-indigo-300 border-8 border-indigo-200"
          onClick={() => {
            checkUpdateHandler();
          }}
        >
          <div className="flex justify-center items-center align-middle my-auto">
            <FontAwesomeIcon
              icon={["fas", "sync-alt"]}
              className="text-white mt-3"
            ></FontAwesomeIcon>
          </div>
          <div className="font-sans font-light text-indigo-100">Update</div>
        </div>
      </div>
    </div>
  );
}
