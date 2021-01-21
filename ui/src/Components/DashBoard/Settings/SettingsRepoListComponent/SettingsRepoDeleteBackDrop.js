import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React from "react";
import SettingsDataFileComponent from "./SettingsRepoDeleteComponent";

export default function SettingsRepoDeleteBackDrop(props) {
  const { setBackdropToggle, deleteRepoData } = props;

  return (
    <div
      className="fixed w-full h-full top-0 left-0 right-0 flex xl:overflow-auto lg:overflow-auto md:overflow-none sm:overflow-none"
      id="settings-backdrop"
      style={{ background: "rgba(0,0,0,0.7)" }}
      onClick={(event) => {
        if (event.target.id === "settings-backdrop") {
          setBackdropToggle(false);
        }
      }}
    >
      <SettingsDataFileComponent
        deleteRepoData={deleteRepoData}
      ></SettingsDataFileComponent>
      <div
        className="w-14 h-14 mr-5 mt-6 rounded-full bg-red-500 text-white flex justify-center items-center shadow cursor-pointer fixed right-0 top-0"
        onClick={() => {
          setBackdropToggle(false);
        }}
      >
        <FontAwesomeIcon
          className="flex text-center text-3xl my-auto"
          icon={["fas", "times"]}
        ></FontAwesomeIcon>
      </div>
    </div>
  );
}
