import { faInfoCircle } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useContext, useState } from "react";
import { AddRepoActionTypes } from "../../add-new-repo-state/actions";
import { AddRepoContext } from "../../add-new-repo-state/addRepoContext";
import HTTPSAuthHintComponent from "./HTTPSAuthHintComponent";

export default function HTTPSAuthForm() {
  const { state, dispatch } = useContext(AddRepoContext);
  const [showHint, setShowHint] = useState<boolean>(false);

  return (
    <form className="my-4 mx-auto">
      <div
        className="flex mx-auto transition-all my-4 justify-center items-center w-11/12 p-2 border-b border-dashed font-sans font-semibold cursor-pointer text-indigo-500 hover:text-indigo-800"
        onClick={() => {
          setShowHint(!showHint);
        }}
      >
        <FontAwesomeIcon className="mx-2" icon={faInfoCircle}></FontAwesomeIcon>
        <div>KNOW MORE ABOUT HTTPS AUTH</div>
      </div>
      <div className="my-2">
        <input
          id="userName"
          type="text"
          placeholder="User Name"
          autoComplete="username"
          className="border-blue-50 rounded-md border-2 my-3 outline-none p-4 w-11/12 shadow-md font-sans text-gray-700 focus:ring-2 focus:ring-blue-200 focus:ring-opacity-60"
          onChange={(event) => {
            dispatch({
              type: AddRepoActionTypes.SET_HTTPS_AUTH_INPUTS,
              payload: {
                userName: event.currentTarget.value,
                password: state.httpsAuthInputs.password,
              },
            });
          }}
          onClick={() => {
            dispatch({
              type: AddRepoActionTypes.SET_ALERT_STATUS,
              payload: "",
            });
          }}
        ></input>
      </div>
      <div className="my-2">
        <input
          id="password"
          type="password"
          placeholder="Password"
          autoComplete="current-password"
          className="border-blue-50 rounded-md border-2 my-3 outline-none p-4 w-11/12 shadow-md font-sans text-gray-700 focus:ring-2 focus:ring-blue-200 focus:ring-opacity-60"
          onChange={(event) => {
            dispatch({
              type: AddRepoActionTypes.SET_HTTPS_AUTH_INPUTS,
              payload: {
                userName: state.httpsAuthInputs.userName,
                password: event.currentTarget.value,
              },
            });
          }}
          onClick={() => {
            dispatch({
              type: AddRepoActionTypes.SET_ALERT_STATUS,
              payload: "",
            });
          }}
        ></input>
      </div>
      {showHint ? (
        <HTTPSAuthHintComponent hideHint={setShowHint}></HTTPSAuthHintComponent>
      ) : null}
    </form>
  );
}
