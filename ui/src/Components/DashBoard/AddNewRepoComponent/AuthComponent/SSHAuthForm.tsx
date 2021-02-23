import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useContext } from "react";
import { AddRepoActionTypes } from "../add-new-repo-state/actions";
import { AddRepoContext } from "../add-new-repo-state/addRepoContext";

export default function SSHAuthForm() {
  const { state, dispatch } = useContext(AddRepoContext);

  const hintData: { key: string; path: string }[] = [
    {
      key: "Windows",
      path: "C:\\Users\\name\\.ssh\\id_rsa",
    },
    {
      key: "Linux",
      path: "~/.ssh/id_rsa",
    },
  ];

  return (
    <>
      <div className="w-11/12 shadow rounded-md border flex items-center justify-between my-4 mx-auto text-indigo-800 focus-within:ring-2 focus-within:ring-opacity-20 focus-within:ring-indigo-400">
        <div className="border py-3 px-6 text-center">
          <FontAwesomeIcon icon={["fas", "key"]}></FontAwesomeIcon>
        </div>
        <div className="w-full">
          <input
            value={state.sshKeyPath}
            className="border-0 outline-none w-full p-2"
            placeholder="Enter the full path for the SSH private key"
            onClick={() => {
              dispatch({
                type: AddRepoActionTypes.SET_ALERT_STATUS,
                payload: "",
              });
            }}
            onChange={(event) => {
              dispatch({
                type: AddRepoActionTypes.SET_SSH_KEY_PATH,
                payload: event.currentTarget.value,
              });
            }}
          ></input>
        </div>
      </div>
      <div className="my-2 font-sans font-light text-gray-600 mx-8">
        <span className="mx-1">E.g: </span>
        {hintData.map((item) => {
          return (
            <span key={item.key}>
              <span className="font-semibold text-left mx-2 text-gray-800">
                {item.key}
              </span>
              {item.path}
            </span>
          );
        })}
      </div>
    </>
  );
}
