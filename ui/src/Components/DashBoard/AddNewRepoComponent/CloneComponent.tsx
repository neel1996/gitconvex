import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useContext } from "react";
import { AddRepoActionTypes } from "./add-new-repo-state/actions";
import { AddRepoContext } from "./add-new-repo-state/addRepoContext";

export default function CloneComponent() {
  const { state, dispatch } = useContext(AddRepoContext);

  return (
    <div className="w-11/12 shadow rounded-md border flex items-center justify-between mx-auto text-indigo-800 focus-within:ring-2 focus-within:ring-opacity-20 focus-within:ring-indigo-400">
      <div className="border py-3 px-6 text-center">
        <FontAwesomeIcon icon={["fas", "link"]}></FontAwesomeIcon>
      </div>
      <div className="w-full">
        <input
          value={state.cloneURL}
          className="border-0 outline-none w-full p-2"
          placeholder="Enter the remote repo URL"
          onClick={() => {
            dispatch({
              type: AddRepoActionTypes.SET_ALERT_STATUS,
              payload: "",
            });
          }}
          onChange={(event) => {
            dispatch({
              type: AddRepoActionTypes.SET_CLONE_URL,
              payload: event.currentTarget.value,
            });
          }}
        ></input>
      </div>
    </div>
  );
}
