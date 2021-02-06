import React, { useContext } from "react";
import { AddRepoContext } from "../add-new-repo-state/addRepoContext";
import AuthOptionsComponent from "./AuthOptionsComponent";
import HTTPSAuthForm from "./HTTPSAuthForm/HTTPSAuthForm";
import SSHAuthForm from "./SSHAuthForm";

export default function AuthComponent() {
  const { state } = useContext(AddRepoContext);

  return (
    <div>
      <div className="my-3 mx-auto text-center">
        <div className="font-sans font-light my-4 mx-auto w-11/12 text-gray-600">
          If the repo is secured / private then choose the appropriate
          authentication option
        </div>
        <AuthOptionsComponent></AuthOptionsComponent>
        {state.authMethod === "https" ? <HTTPSAuthForm></HTTPSAuthForm> : null}
        {state.authMethod === "ssh" ? <SSHAuthForm></SSHAuthForm> : null}
      </div>
    </div>
  );
}
