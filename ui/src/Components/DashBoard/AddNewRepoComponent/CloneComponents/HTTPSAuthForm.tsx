import React, { useContext } from "react";
import { AddRepoActionTypes } from "../add-new-repo-state/actions";
import { AddRepoContext } from "../add-new-repo-state/addRepoContext";

export default function HTTPSAuthForm() {
  const { state, dispatch } = useContext(AddRepoContext);

  return (
    <div className="my-4 mx-auto">
      <div className="text-sm font-sans font-light my-1 text-center mx-auto text-pink-500">
        Basic Authentication will not work if 2-Factor Authentication is enabled
      </div>
      <div className="my-2">
        <input
          id="repoNameText"
          type="text"
          placeholder="User Name"
          className="border-blue-50 rounded-md border-2 my-3 outline-none p-4 w-11/12 shadow-md font-sans text-gray-700"
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
          id="repoNameText"
          type="password"
          placeholder="Password"
          className="border-blue-50 rounded-md border-2 my-3 outline-none p-4 w-11/12 shadow-md font-sans text-gray-700"
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
    </div>
  );
}
