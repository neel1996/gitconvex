import React, { useRef } from "react";
import axios from "axios";
import { globalAPIEndpoint } from "../../../../../../../util/env_config";
import { faCheck, faTimes } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

export default function AddRemoteRepoFormComponent(props) {
  const {
    setReloadView,
    setRemoteForm,
    setFieldMissing,
    setInvalidUrl,
    setAddNewRemote,
    setAddRemoteStatus,
    repoId,
    setStatusCheck,
    setRemoteOperation,
  } = props;

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
          setFieldMissing(false);
          setAddRemoteStatus(false);
          setInvalidUrl(false);
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
        setInvalidUrl(true);
      } else {
        axios({
          url: globalAPIEndpoint,
          method: "POST",
          data: {
            query: `
                  mutation {
                    addRemote(repoId: "${repoId}", remoteName: "${remoteName}", remoteUrl: "${remoteUrl}"){
                      status
                    }
                  }
                `,
          },
        })
          .then((res) => {
            const { status } = res.data.data.addRemote;
            if (status === "REMOTE_ADD_SUCCESS") {
              remoteNameRef.current.value = "";
              remoteUrlRef.current.value = "";

              setRemoteForm(false);
              setAddNewRemote(true);
              setReloadView(true);
            } else {
              setAddRemoteStatus(true);
            }
            setStatusCheck(false);
            setRemoteOperation(" ");
          })
          .catch((err) => {
            console.log(err);
            setStatusCheck(true);
            setRemoteOperation("add");
          });
      }
    } else {
      setAddNewRemote(false);
      setInvalidUrl(false);
      setFieldMissing(true);
    }
  };

  return (
    <form className="flex items-center w-full my-6 form--data">
      <div className="w-1/4 mx-auto">
        {formAddRemote("remoteName", "Remote name")}
      </div>
      <div className="w-1/2 mx-auto">
        {formAddRemote("remoteURL", "Remote URL")}
      </div>
      <div
        className="flex items-center text-center justify-evenly"
        style={{ outline: "none", width: "22%" }}
      >
        <div
          className="items-center w-5/12 p-1 py-2 mx-auto text-base font-semibold bg-blue-500 rounded cursor-pointer xl:text-lg lg:text-lg md:text-base hover:bg-blue-700"
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
          className="items-center w-5/12 p-1 py-2 mx-auto text-base font-semibold bg-red-500 rounded cursor-pointer xl:text-lg lg:text-lg md:text-base hover:bg-red-600"
          onClick={() => {
            setAddNewRemote(true);
            setRemoteForm(false);
            setFieldMissing(false);
            setInvalidUrl(false);
            setAddRemoteStatus(false);
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
