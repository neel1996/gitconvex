import { faCogs } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useReducer } from "react";
import SettingsDataFileComponent from "./SettingsDataFileComponent";
import SettingsPortComponent from "./SettingsPortComponent";
import SettingsRepoListComponent from "./SettingsRepoListComponent/SettingsRepoListComponent";

const settingsInitialState = {
  viewReloadCount: 0,
};
export const SettingsContext = React.createContext(settingsInitialState);

export default function Settings() {
  const settingsReducer = (state, payload) => {
    if (payload > 0) {
      return {
        viewReloadCount: state.viewReloadCount + 1,
      };
    } else {
      return state;
    }
  };
  const [state, dispatch] = useReducer(settingsReducer, settingsInitialState);

  return (
    <>
      <div className="block w-full h-full pt-5 pb-10 overflow-auto">
        <div className="font-sans text-6xl my-4 mx-10 text-gray-700 block items-center align-middle">
          <FontAwesomeIcon className="text-5xl" icon={faCogs}></FontAwesomeIcon>
          <span className="mx-10">Settings</span>
        </div>
        <div className="block my-10 justify-center mx-auto w-11/12">
          <SettingsContext.Provider value={{ state, dispatch }}>
            <SettingsDataFileComponent></SettingsDataFileComponent>
            <SettingsRepoListComponent></SettingsRepoListComponent>
            <SettingsPortComponent></SettingsPortComponent>
          </SettingsContext.Provider>
        </div>
      </div>
    </>
  );
}
