import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";
import React, { useContext, useEffect, useState } from "react";
import { animated, useSpring } from "react-spring";
import { DELETE_PRESENT_REPO } from "../../../../actionStore";
import { ContextProvider } from "../../../../context";
import { globalAPIEndpoint } from "../../../../util/env_config";
import InfiniteLoader from "../../../Animations/InfiniteLoader";
import "../../../styles/AddRepoForm.css";

export default function AddRepoForm(props) {
  library.add(fas);
  const { state, dispatch } = useContext(ContextProvider);
  const [repoNameState, setRepoName] = useState("");
  const [repoPathState, setRepoPath] = useState("");
  const [cloneUrlState, setCloneUrlState] = useState("");
  const [repoAddFailed, setRepoAddFailed] = useState(false);
  const [repoAddSuccess, setRepoAddSuccess] = useState(false);
  const [inputInvalid, setInputInvalid] = useState(false);
  const [loading, setLoading] = useState(false);
  const [authOption, setAuthOption] = useState("noauth");
  const [authInputs, setAuthInputs] = useState({
    userName: "",
    password: "",
  });

  const [cloneSwitch, setCloneSwitch] = useState(false);
  const [initSwitch, setInitSwitch] = useState(false);

  const switchAnimationEnter = useSpring({
    config: { duration: 1500, tension: 500 },
    from: {
      transform: "translate(0em, 0em)",
    },
    to: {
      transform: "translate(2em, 0em)",
    },
  });

  const switchAnimationExit = useSpring({
    config: { duration: 500, tension: 500 },
    from: {
      transform: "translate(2em, 0em)",
    },
    to: {
      transform: "translate(0em, 0em)",
    },
  });

  const authRadio = [
    {
      key: "noauth",
      label: "No Authentication",
    },
    {
      key: "ssh",
      label: "SSH Authentication",
    },
    {
      key: "https",
      label: "HTTPS Authentication",
    },
  ];

  useEffect(() => {
    setAuthOption("noauth");
    setAuthInputs({ userName: "", password: "" });
    if (state.shouldAddFormClose) {
      props.formEnable(false);
    }
  }, [state, props]);

  function storeRepoAPI(repoName, repoPath) {
    if (repoName && repoPath) {
      if (repoName.match(/[^a-zA-Z0-9-_.\s]/gi)) {
        setInputInvalid(true);
        setRepoAddFailed(true);
        return;
      }

      let initCheck = false;
      let cloneCheck = false;
      let cloneUrl = cloneUrlState;
      repoPath = repoPath.replace(/\\/gi, "\\\\");

      if (cloneSwitch && !cloneUrlState) {
        setRepoAddFailed(true);
        return false;
      }

      if (cloneUrl.match(/[^a-zA-Z0-9-_.~@#$%:/]/gi)) {
        setInputInvalid(true);
        setRepoAddFailed(true);
        return;
      }

      let userName = "";
      let password = "";

      if (initSwitch) {
        initCheck = true;
      } else if (cloneSwitch && cloneUrl) {
        cloneCheck = true;
        if (authOption === "https") {
          userName = authInputs.userName;
          password = authInputs.password;
        }
      }

      setLoading(true);

      axios({
        url: globalAPIEndpoint,
        method: "POST",
        data: {
          query: `
              mutation {
                addRepo(repoName: "${repoName}", repoPath: "${repoPath}", initSwitch: ${initCheck}, cloneSwitch: ${cloneCheck}, repoURL: "${cloneUrl}", authOption: "${authOption}", userName: "${userName}", password: "${password}"){
                  repoId
                  status
                }
              }
            `,
        },
      })
        .then((res) => {
          setLoading(false);
          setInputInvalid(false);

          if (res.data.data && !res.data.error) {
            const { status } = res.data.data.addRepo;

            if (status && !status.match(/Failed/g)) {
              setRepoAddSuccess(true);
              setRepoAddFailed(false);
              setCloneSwitch("");
              setInitSwitch("");

              setRepoName("");
              setRepoPath("");
              setCloneUrlState("");

              dispatch({
                action: DELETE_PRESENT_REPO,
                payload: [],
              });

              console.log(state.presentRepo);
            } else {
              setRepoAddFailed(true);
              setRepoAddSuccess(false);
            }
          } else {
            setRepoAddFailed(true);
            setRepoAddSuccess(false);
          }
        })
        .catch((err) => {
          setLoading(false);

          console.log(err);
          setRepoAddFailed(true);
          setRepoAddSuccess(false);
        });
    } else {
      setInputInvalid(false);
      setRepoAddFailed(true);
    }
  }

  function resetAlertBanner() {
    setRepoAddFailed(false);
    setRepoAddSuccess(false);
  }

  function repoAddStatusBanner() {
    if (repoAddSuccess) {
      return <div className="alert-success">New repo added</div>;
    } else if (repoAddFailed) {
      return (
        <div className="alert-failure">
          Process failed! Please try again
          {inputInvalid ? (
            <div className="font-semibold">Invalid input paremeters!</div>
          ) : null}
        </div>
      );
    } else {
      return null;
    }
  }

  function switchComponent(operation) {
    return (
      <div
        key={`switch-${operation}`}
        className={`toggle-switch ${
          operation === "clone" && cloneSwitch ? "bg-green-400" : "bg-gray-200"
        }
        ${operation === "init" && initSwitch ? "bg-blue-400" : "bg-gray-200"}`}
        onClick={(event) => {
          if (operation === "clone") {
            if (!cloneSwitch) {
              setCloneSwitch(true);
              setInitSwitch(false);
            } else {
              setCloneSwitch(false);
            }
          } else {
            if (!initSwitch) {
              setInitSwitch(true);
              setCloneSwitch(false);
            } else {
              setInitSwitch(false);
            }
          }
        }}
      >
        {operation === "clone" ? (
          <animated.div
            className="toggle-switch--pill"
            id={`${operation}-switch`}
            style={cloneSwitch ? switchAnimationEnter : switchAnimationExit}
          ></animated.div>
        ) : (
          <animated.div
            className="toggle-switch--pill"
            id={`${operation}-switch`}
            style={initSwitch ? switchAnimationEnter : switchAnimationExit}
          ></animated.div>
        )}
      </div>
    );
  }

  function addRepoFormContainer() {
    return (
      <div className="repo-form block">
        {repoAddStatusBanner()}
        <div className="repo-form--header">Enter Repo Details</div>
        <div>
          <input
            id="repoNameText"
            type="text"
            placeholder="Enter a Repository Name"
            className="repo-form--input"
            onChange={(event) => {
              setRepoName(event.target.value);
            }}
            value={repoNameState}
            onClick={() => {
              resetAlertBanner();
            }}
          ></input>
        </div>
        <div>
          <input
            id="repoPathText"
            type="text"
            placeholder={
              cloneSwitch
                ? "Enter base directory path"
                : "Enter repository path"
            }
            className="repo-form--input"
            onChange={(event) => {
              setRepoPath(event.target.value);
            }}
            value={repoPathState}
            onClick={() => {
              resetAlertBanner();
            }}
          ></input>
        </div>
        {cloneSwitch && repoPathState && repoNameState ? (
          <div className="repo-form--clone">
            The repo will be cloned to
            <span className="mx-3 text-center font-sans font-semibold border-b-2 border-dashed">
              {repoPathState}
              <>{repoPathState.includes("\\") ? "\\" : "/"}</>
              {repoNameState}
            </span>
          </div>
        ) : null}
        <div className="repo-form--options">
          <div className="options--switch">
            <div>{switchComponent("clone")}</div>
            <div className="options--label">Clone from remote</div>
          </div>

          <div className="options--switch">
            <div>{switchComponent("init")}</div>
            <div className="options--label">Initialize a new repo</div>
          </div>
        </div>
        {cloneSwitch ? (
          <>
            <div className="option--clone">
              <div className="option--clone--icon">
                <FontAwesomeIcon icon={["fas", "link"]}></FontAwesomeIcon>
              </div>
              <div className="w-5/6">
                <input
                  value={cloneUrlState}
                  className="border-0 outline-none w-full p-2"
                  placeholder="Enter the remote repo URL"
                  onClick={() => {
                    setRepoAddFailed(false);
                  }}
                  onChange={(event) => {
                    setCloneUrlState(event.target.value);
                  }}
                ></input>
              </div>
            </div>
            <div className="my-3 mx-auto text-center">
              <div className="font-sans font-light my-4 mx-auto w-11/12 text-gray-600">
                If the repo is secured / private then choose the appropriate
                authentication option
              </div>
              <div className="flex gap-4 justify-center mx-auto items-center align-middle">
                {authRadio.map((item) => {
                  return (
                    <div
                      className="flex gap-4 items-center align-middle"
                      key={item.key}
                    >
                      <input
                        type="radio"
                        name="authRadio"
                        id={item.key}
                        value={item.key}
                        onChange={(e) => {
                          setAuthOption(e.target.value);
                        }}
                      ></input>
                      <label
                        htmlFor={item.key}
                        className="font-sans text-sm font-light cursor-pointer"
                      >
                        {item.label}
                      </label>
                    </div>
                  );
                })}
              </div>
              {authOption === "https" ? (
                <div className="my-4 mx-auto">
                  <div className="text-sm font-sans font-light my-1 text-center mx-auto text-pink-500">
                    Basic Authentication will not work if 2-Factor
                    Authentication is enabled
                  </div>
                  <div className="my-2">
                    <input
                      id="repoNameText"
                      type="text"
                      placeholder="User Name"
                      className="repo-form--input"
                      onChange={(event) => {
                        setAuthInputs({
                          userName: event.currentTarget.value,
                          password: authInputs.password,
                        });
                      }}
                      onClick={() => {
                        resetAlertBanner();
                      }}
                    ></input>
                  </div>
                  <div className="my-2">
                    <input
                      id="repoNameText"
                      type="password"
                      placeholder="Password"
                      className="repo-form--input"
                      onChange={(event) => {
                        setAuthInputs({
                          userName: authInputs.userName,
                          password: event.currentTarget.value,
                        });
                      }}
                      onClick={() => {
                        resetAlertBanner();
                      }}
                    ></input>
                  </div>
                </div>
              ) : null}
            </div>
          </>
        ) : null}
        <div className="repo-form--action">
          <div
            className="btn btn-danger"
            id="addRepoClose"
            onClick={() => {
              props.formEnable(false);
            }}
          >
            Close
          </div>
          <div
            className="btn btn-success"
            id="addRepoSubmit"
            onClick={() => {
              storeRepoAPI(repoNameState, repoPathState);
            }}
          >
            Submit
          </div>
        </div>
      </div>
    );
  }

  return (
    <div
      className={`repo-form--status xl:w-1/2 lg:w-2/3 md:w-3/4 sm:w-11/12 w-11/12 ${
        loading ? "border-dashed border-2" : ""
      }`}
    >
      {loading ? (
        <>
          <div className="status--label">Repo setup in progress...</div>
          <div className="animation-loader">
            <InfiniteLoader loadAnimation={loading}></InfiniteLoader>
          </div>
        </>
      ) : (
        addRepoFormContainer()
      )}
    </div>
  );
}
