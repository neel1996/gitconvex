import React, { useRef } from "react";
import { faCheck, faTimes } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export default function AddRemoteRepoFormComponent(props) {
  const remoteNameRef = useRef();
  const remoteUrlRef = useRef();

  const formAddRemote = (formId, placeholder) => {
    return (
      <input
        type="text"
        autoComplete="off"
        id={formId}
        className={`rounded w-full py-2 border-2 text-center xl:text-lg lg:text-lg md:text-base text-base items-center text-gray-800 bg-white`}
        style={{ borderColor: "rgb(113 166 196 / 33%)" }}
        placeholder={placeholder}
        ref={formId === "remoteName" ? remoteNameRef : remoteUrlRef}
        onChange={(event) => {
          props.setFieldMissing(false);
          props.setAddRemoteStatus(false);
          props.setInvalidUrl(false);
          const remoteNameVal = event.target.value;
          if (
            event.target.id === "remoteName" &&
            remoteNameVal.match(/[\s\\//*]/gi)
          ) {
            event.target.value = remoteNameVal.replace(/[\s\\//*]/gi, "-");
          }
        }}
      ></input>
    );
  };

  const addRemote = () => {
    let remoteName = remoteNameRef.current.value.trim();
    let remoteUrl = remoteUrlRef.current.value.trim();

    if (remoteName && remoteUrl && remoteUrl.match(/[^ ]*/g)) {
      if (remoteUrl.match(/(\s)/g)) {
        props.setInvalidUrl(true);
      } else {
        //TODO: Add axios
        let status = "success";
        if (status === "success") {
          localStorage.setItem(
            remoteName,
            JSON.stringify({
              remoteName: remoteName,
              remoteUrl: remoteUrl,
            })
          );
          props.setRemoteForm(false);
          props.setAddNewRemote(true);
          props.setReloadView(true);
          remoteNameRef.current.value = "";
          remoteUrlRef.current.value = "";
        } else {
          props.setAddRemoteStatus(true);
        }
      }
    } else {
      props.setAddNewRemote(false);
      props.setInvalidUrl(false);
      props.setFieldMissing(true);
    }
  };

  return (
    <form className="form--data flex w-full items-center my-6">
      <div className="mx-auto w-1/4">
        {formAddRemote("remoteName", "Remote name")}
      </div>
      <div className="mx-auto w-1/2">
        {formAddRemote("remoteURL", "Remote URL")}
      </div>
      <div
        className="text-center flex items-center justify-evenly"
        style={{ outline: "none", width: "22%" }}
      >
        <div
          className="xl:text-lg lg:text-lg md:text-base text-base items-center p-1 py-2 rounded w-5/12 mx-auto cursor-pointer bg-blue-500 hover:bg-blue-700 font-semibold"
          onClick={() => {
            addRemote();
          }}
        >
          <FontAwesomeIcon
            icon={faCheck}
            className="text-white"
          ></FontAwesomeIcon>
        </div>
        <div
          className="xl:text-lg lg:text-lg md:text-base text-base items-center p-1 py-2 rounded w-5/12 mx-auto cursor-pointer bg-red-500 hover:bg-red-600 font-semibold"
          onClick={() => {
            props.setAddNewRemote(true);
            props.setRemoteForm(false);
            props.setFieldMissing(false);
            props.setInvalidUrl(false);
            props.setAddRemoteStatus(false);
          }}
        >
          <FontAwesomeIcon
            icon={faTimes}
            className="text-white"
          ></FontAwesomeIcon>
        </div>
      </div>
    </form>
  );
}
