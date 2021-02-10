import React, { useEffect, useState } from "react";
import "@fortawesome/react-fontawesome";
import { faCodeBranch } from "@fortawesome/free-solid-svg-icons";
import axios from "axios";
import { globalAPIEndpoint } from "../../../../../../../util/env_config";
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
  const [statusCheck, setStatusCheck] = useState(false);
  const [remoteOperation, setRemoteOperation] = useState("add");
  const [fetchStatus, setFetchStatus] = useState(false);
  const [remoteDetails, setRemoteDetails] = useState([]);

  useEffect(() => {
    let remoteArray = [];
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
              query {
                getRemote(repoId: "${repoId}"){
                  remoteDetails
                }
              }
            `,
      },
    })
      .then((res) => {
        if (res.data.data.remoteDetails) {
          res.data.data.remoteDetails.forEach((items) => {
            remoteArray.push(items);
          });
          setRemoteDetails([...remoteArray]);
          setStatusCheck(false);
          setFetchStatus(false);
          setRemoteOperation("");
        }
      })
      .catch(() => {
        setFetchStatus(true);
        // setRemoteDetails(remoteDetails);
        // console.log(err);
      });

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
      className="w-11/12 m-auto rounded-lg xl:w-3/4 lg:w-3/4 md:w-11/12 sm:w-11/12"
      style={{ backgroundColor: "#edf2f7" }}
    >
      {fetchStatus ? (
        <div className="w-3/4 mx-auto">
          <div className="items-center w-full p-6 my-12 text-2xl text-center text-red-900 bg-red-200 border-2 border-red-800 rounded-md text-semibold">
            Failed to fetch remote
          </div>
        </div>
      ) : (
        <>
          {statusCheck ? (
            <div className="w-3/4 mx-auto">
              <div className="items-center w-full p-6 my-12 text-2xl text-center text-red-900 bg-red-200 border-2 border-red-800 rounded-md text-semibold">
                Failed to {remoteOperation} remote
              </div>
            </div>
          ) : (
            <>
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
              <div className="w-full p-2 pt-6 pb-8">
                <div className="flex items-center m-6 font-sans text-xl font-semibold text-gray-800 xl:text-3xl lg:text-3xl md:text-2xl sm:text-xl">
                  <FontAwesomeIcon
                    icon={faCodeBranch}
                    className="mx-2 text-xl xl:text-3xl lg:text-3xl md:text-2xl sm:text-xl"
                  ></FontAwesomeIcon>
                  <div className="pb-2 border-b-4 border-blue-400 border-dashed">
                    Remote details
                  </div>
                  {addNewRemote && remoteDetails.length > 0 ? (
                    <div
                      className="px-3 py-2 mx-6 font-sans text-base text-gray-800 bg-blue-200 rounded cursor-pointer xl:text-lg lg:text-lg md:text-base hover:bg-blue-300 hover:text-gray-900"
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
                  {remoteDetails && remoteDetails.length > 0 ? (
                    <>
                      <div className="flex items-center w-full">
                        <div className="w-1/4 mx-auto font-sans text-lg font-semibold text-center text-gray-600 xl:text-2xl lg:text-2xl md:text-xl">
                          Remote name
                        </div>
                        <div className="w-7/12 mx-auto font-sans text-lg font-semibold text-center text-gray-600 xl:text-2xl lg:text-2xl md:text-xl">
                          Remote URL
                        </div>
                        <div
                          className="mx-auto font-sans text-lg font-semibold text-center text-gray-600 xl:text-2xl lg:text-2xl md:text-xl"
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
                          repoId={repoId}
                          setStatusCheck={setStatusCheck}
                          setRemoteOperation={setRemoteOperation}
                        ></AddRemoteRepoFormComponent>
                      ) : null}
                      <div
                        className="flex flex-col items-center w-full mt-3 mb-4 overflow-auto"
                        style={{ maxHeight: "350px" }}
                      >
                        {remoteDetails.map((items) => {
                          const { remoteName, remoteUrl } = items;
                          return (
                            <RemoteCard
                              key={remoteName}
                              remoteName={remoteName}
                              remoteUrl={remoteUrl}
                              setFieldMissing={setFieldMissing}
                              setInvalidUrl={setInvalidUrl}
                              setAddRemoteStatus={setAddRemoteStatus}
                              setDeleteFailed={setDeleteFailed}
                              setReloadView={setReloadView}
                              repoId={repoId}
                              setStatusCheck={setStatusCheck}
                              setRemoteOperation={setRemoteOperation}
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
                      repoId={repoId}
                      setStatusCheck={setStatusCheck}
                    ></AddRemoteRepoFormComponent>
                  )}
                </div>
              </div>
            </>
          )}
        </>
      )}
    </div>
  );
}
