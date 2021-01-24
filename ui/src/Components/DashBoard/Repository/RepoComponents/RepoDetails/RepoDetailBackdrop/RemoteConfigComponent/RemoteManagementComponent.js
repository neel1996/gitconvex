import React, { useEffect, useState } from "react";
import "@fortawesome/react-fontawesome";
import { faCodeBranch } from "@fortawesome/free-solid-svg-icons";

// import axios from "axios";
// import { globalAPIEndpoint } from "../../../../../../util/env_config";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import AddRemoteRepoFormComponent from "./AddRemoteRepoFormComponent";
import RemoteCard from "./RemoteCard";

export default function AddRemoteRepoComponent({ repoId }) {
  const [reloadView, setReloadView] = useState(false);
  const [fieldMissing, setFieldMissing] = useState(false);
  const [invalidUrl, setInvalidUrl] = useState(false);
  const [remoteForm, setRemoteForm] = useState(false);
  const [deleteFailed, setDeleteFailed] = useState(false);
  const [addNewRemote, setAddNewRemote] = useState(true);
  const [addRemoteStatus, setAddRemoteStatus] = useState(false);
  const [remoteDetails, setRemoteDetails] = useState([]);

  useEffect(() => {
    let remoteArray = [];
    //TODO: Add API call here and load it into the state

    for (let i = 0; i < localStorage.length; i++) {
      var key = localStorage.key(i);
      var value = JSON.parse(localStorage.getItem(key));
      remoteArray.push(value);
    }

    setRemoteDetails([...remoteArray]);

    return () => {
      setRemoteDetails([]);
      setReloadView(false);
    };
  }, [reloadView]);

  const statusPillComponent = (border, bgColor, textColor, message) => {
    return (
      <div
        className={`${border} ${bgColor} ${textColor} border-b-2 font-sans text-xl border-dashed text-center rounded-b-none rounded-t-lg w-full py-6`}
      >
        {message}
      </div>
    );
  };

  return (
    <div
      className="xl:w-3/4 lg:w-3/4 md:w-11/12 sm:w-11/12 w-11/12 m-auto rounded-lg"
      style={{ backgroundColor: "#edf2f7" }}
    >
      {addRemoteStatus
        ? statusPillComponent(
            "border-red-800",
            "bg-red-200",
            "text-red-900",
            "Remote name already exist!"
          )
        : null}
      {fieldMissing
        ? statusPillComponent(
            "border-indigo-800",
            "bg-indigo-200",
            "text-indigo-900",
            "One or more required parameters are empty!"
          )
        : null}
      {invalidUrl
        ? statusPillComponent(
            "border-yellow-800",
            "bg-yellow-200",
            "text-yellow-900",
            "URL with whitespace is invalid!"
          )
        : null}
      {deleteFailed
        ? statusPillComponent(
            "border-red-800",
            "bg-red-200",
            "text-red-900",
            "Failed to delete"
          )
        : null}
      <div className="w-full p-2 pb-8 pt-6">
        <div className="xl:text-3xl lg:text-3xl md:text-2xl sm:text-xl text-xl m-6 font-sans text-gray-800 font-semibold flex items-center">
          <FontAwesomeIcon
            icon={faCodeBranch}
            className="xl:text-3xl lg:text-3xl md:text-2xl sm:text-xl text-xl mx-2"
          ></FontAwesomeIcon>
          <div className="border-b-4 pb-2 border-dashed border-blue-400">
            Remote details
          </div>
          {addNewRemote && remoteDetails.length > 0 ? (
            <div
              className="mx-6 px-3 py-2 font-sans rounded xl:text-lg lg:text-lg md:text-base text-base cursor-pointer bg-blue-200 text-gray-800 hover:bg-blue-300 hover:text-gray-900"
              onClick={() => {
                setAddNewRemote(false);
                setRemoteForm(true);
              }}
            >
              Add new remote
            </div>
          ) : null}
        </div>
        <div className="w-11/12 mx-auto">
          {remoteDetails.length > 0 ? (
            <>
              <div className="flex items-center w-full">
                <div className="font-sans w-1/4 xl:text-2xl lg:text-2xl md:text-xl text-lg mx-auto text-center font-semibold text-gray-600">
                  Remote name
                </div>
                <div className="font-sans xl:text-2xl lg:text-2xl md:text-xl text-lg w-7/12 mx-auto text-center font-semibold text-gray-600">
                  Remote URL
                </div>
                <div
                  className="font-sans xl:text-2xl lg:text-2xl md:text-xl text-lg mx-auto text-center font-semibold text-gray-600"
                  style={{ width: "22%" }}
                >
                  Actions
                </div>
              </div>
              {remoteForm ? (
                <AddRemoteRepoFormComponent
                  setReloadView={setReloadView}
                  setRemoteForm={setRemoteForm}
                  setFieldMissing={setFieldMissing}
                  setInvalidUrl={setInvalidUrl}
                  setAddNewRemote={setAddNewRemote}
                  setAddRemoteStatus={setAddRemoteStatus}
                ></AddRemoteRepoFormComponent>
              ) : null}
              <div
                className="mt-3 w-full mb-4 overflow-auto flex flex-col items-center"
                style={{ maxHeight: "350px" }}
              >
                {remoteDetails.map((items) => {
                  console.log(remoteDetails.length);
                  const { remoteName, remoteUrl } = items;
                  return (
                    <RemoteCard
                      key={remoteName}
                      remoteName={remoteName}
                      remoteUrl={remoteUrl}
                      remoteDetails={remoteDetails}
                      setFieldMissing={setFieldMissing}
                      setInvalidUrl={setInvalidUrl}
                      setAddRemoteStatus={setAddRemoteStatus}
                      setDeleteFailed={setDeleteFailed}
                      setReloadView={setReloadView}
                    ></RemoteCard>
                  );
                })}
              </div>
            </>
          ) : (
            <AddRemoteRepoFormComponent
              setReloadView={setReloadView}
              setRemoteForm={setRemoteForm}
              setInvalidUrl={setInvalidUrl}
              setFieldMissing={setFieldMissing}
              setAddNewRemote={setAddNewRemote}
              setAddRemoteStatus={setAddRemoteStatus}
            ></AddRemoteRepoFormComponent>
          )}
        </div>
      </div>
    </div>
  );
}
