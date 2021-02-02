import React, { useContext } from "react";
import { AddRepoActionTypes } from "../add-new-repo-state/actions";
import { AddRepoContext } from "../add-new-repo-state/addRepoContext";

export default function HTTPSAuthForm() {
  const { state, dispatch } = useContext(AddRepoContext);

  return (
    <form className="my-4 mx-auto">
      <div className="text-sm font-sans font-light my-1 text-center mx-auto text-pink-500">
        Basic Authentication will not work if 2-Factor Authentication is enabled
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
    </form>
  );
}
