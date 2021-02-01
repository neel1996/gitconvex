import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useContext } from "react";
import { AddRepoActionTypes } from "../add-new-repo-state/actions";
import {
  AddRepoContext,
  authMethodType
} from "../add-new-repo-state/addRepoContext";
import HTTPSAuthForm from "./HTTPSAuthForm";

export default function CloneComponent() {
  const { state, dispatch } = useContext(AddRepoContext);

  const authRadio: { key: authMethodType; label: string }[] = [
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

  return (
    <>
      <div className="w-11/12 shadow rounded-md border flex items-center justify-between mx-auto text-indigo-800">
        <div className="border py-3 px-6 text-center">
          <FontAwesomeIcon icon={["fas", "link"]}></FontAwesomeIcon>
        </div>
        <div className="w-5/6">
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
                    dispatch({
                      type: AddRepoActionTypes.SET_AUTH_OPTION,
                      payload: e.currentTarget.value,
                    });
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
        {state.authOption === "https" ? <HTTPSAuthForm></HTTPSAuthForm> : null}
      </div>
    </>
  );
}
