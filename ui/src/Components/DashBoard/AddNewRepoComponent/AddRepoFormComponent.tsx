import { FormEvent, useContext } from "react";
import { AddRepoActionTypes } from "./add-new-repo-state/actions";
import { AddRepoContext } from "./add-new-repo-state/addRepoContext";
import AddRepoActionButtonsComponent from "./AddRepoActionButtonsComponent";
import AddRepoStatusAlert from "./AddRepoStatusAlert";
import CloneComponent from "./CloneComponents/CloneComponent";
import ToggleSwitchComponent from "./ToggleSwitchComponent";

export default function AddRepoFormComponent() {
  const { state, dispatch } = useContext(AddRepoContext);

  function onRepoPathChange(event: FormEvent) {
    const target = event.target as HTMLInputElement;
    dispatch({
      type: AddRepoActionTypes.SET_REPO_PATH,
      payload: target.value,
    });
  }

  function autoDetectRepoName(event: FormEvent) {
    const currentTarget = event.currentTarget as HTMLInputElement;

    if (state.repoName === "") {
      let pathValue = currentTarget.value;
      if (pathValue.includes("/")) {
        let splitPath = pathValue.split("/");
        let detectedRepoName = splitPath[splitPath.length - 1];
        dispatch({
          type: AddRepoActionTypes.SET_REPO_NAME,
          payload: detectedRepoName,
        });
      }

      if (pathValue.includes("\\")) {
        let splitPath = pathValue.split("\\");
        let detectedRepoName = splitPath[splitPath.length - 1];
        dispatch({
          type: AddRepoActionTypes.SET_REPO_NAME,
          payload: detectedRepoName,
        });
      }
    }
  }

  return (
    <div className="block my-2">
      {state.isInputInvalid ? (
        <AddRepoStatusAlert
          status={state.alertStatus}
          message="Entered inputs are invalid!"
        ></AddRepoStatusAlert>
      ) : null}
      {state.alertStatus === "failed" ? (
        <AddRepoStatusAlert status={state.alertStatus}></AddRepoStatusAlert>
      ) : null}
      {state.alertStatus === "success" ? (
        <AddRepoStatusAlert status={state.alertStatus}></AddRepoStatusAlert>
      ) : null}
      <div className="block text-3xl my-4 text-center text-gray-500 font-sans font-semibold">
        ENTER REPO DETAILS
      </div>
      <div>
        <input
          id="repoNameText"
          type="text"
          placeholder="Enter a Repository Name"
          className="border-blue-50 rounded-md border-2 my-3 outline-none p-4 w-11/12 shadow-md font-sans text-gray-700 focus:ring-2 focus:ring-blue-200 focus:ring-opacity-60"
          onChange={(event) => {
            dispatch({
              type: AddRepoActionTypes.SET_REPO_NAME,
              payload: event.target.value,
            });
          }}
          value={state.repoName}
          onClick={() => {
            dispatch({
              type: AddRepoActionTypes.SET_ALERT_STATUS,
              payload: "",
            });
          }}
        ></input>
      </div>
      <div>
        <input
          id="repoPathText"
          type="text"
          placeholder={
            state.cloneSwitch
              ? "Enter base directory path"
              : "Enter repository path"
          }
          className="border-blue-50 rounded-md border-2 my-3 outline-none p-4 w-11/12 shadow-md font-sans text-gray-700 focus:ring-2 focus:ring-blue-200 focus:ring-opacity-60"
          onChange={onRepoPathChange}
          onBlur={autoDetectRepoName}
          value={state.repoPath}
          onClick={() => {
            dispatch({
              type: AddRepoActionTypes.SET_ALERT_STATUS,
              payload: "",
            });
          }}
        ></input>
      </div>
      {state.cloneSwitch && state.repoPath && state.repoName ? (
        <div className="items-center font-light text-sm my-4 mx-auto text-center text-gray-600 font-sans">
          The repo will be cloned to
          <span className="mx-3 text-center font-sans font-semibold border-b-2 border-dashed">
            {state.repoPath}
            <>{state.repoPath.includes("\\") ? "\\" : "/"}</>
            {state.repoName}
          </span>
        </div>
      ) : null}
      <ToggleSwitchComponent></ToggleSwitchComponent>
      {state.cloneSwitch ? <CloneComponent></CloneComponent> : null}
      <AddRepoActionButtonsComponent></AddRepoActionButtonsComponent>
    </div>
  );
}
