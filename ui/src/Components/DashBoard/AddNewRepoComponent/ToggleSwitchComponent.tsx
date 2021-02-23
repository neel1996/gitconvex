import React, { useContext } from "react";
import { animated, useSpring } from "react-spring";
import { AddRepoActionTypes } from "./add-new-repo-state/actions";
import { AddRepoContext } from "./add-new-repo-state/addRepoContext";

export default function ToggleSwitchComponent() {
  const { state, dispatch } = useContext(AddRepoContext);
  const { cloneSwitch, initSwitch } = state;

  const switchAnimationEnter = useSpring({
    config: { duration: 1500, tension: 500 },
    from: {
      transform: "translate(0em, 0em)",
    },
    to: {
      transform: "translate(2em, 0em)",
    },
  });

  const switchAnimationExit = useSpring({
    config: { duration: 500, tension: 500 },
    from: {
      transform: "translate(2em, 0em)",
    },
    to: {
      transform: "translate(0em, 0em)",
    },
  });

  function switchComponent(operation: string) {
    return (
      <div
        key={`switch-${operation}`}
        className={`rounded-full cursor-pointer flex items-center h-8 py-2 pl-1 shadow-inner align-middle w-16 ${
          operation === "clone" && cloneSwitch ? "bg-green-400" : "bg-gray-200"
        }
      ${operation === "init" && initSwitch ? "bg-blue-400" : "bg-gray-200"}`}
        onClick={() => {
          if (operation === "clone") {
            if (!cloneSwitch) {
              dispatch({
                type: AddRepoActionTypes.SET_CLONE_SWITCH,
                payload: true,
              });
              dispatch({
                type: AddRepoActionTypes.SET_INIT_SWITCH,
                payload: false,
              });
            } else {
              dispatch({
                type: AddRepoActionTypes.SET_CLONE_SWITCH,
                payload: false,
              });
            }
          } else {
            if (!initSwitch) {
              dispatch({
                type: AddRepoActionTypes.SET_CLONE_SWITCH,
                payload: false,
              });
              dispatch({
                type: AddRepoActionTypes.SET_INIT_SWITCH,
                payload: true,
              });
            } else {
              dispatch({
                type: AddRepoActionTypes.SET_INIT_SWITCH,
                payload: false,
              });
            }
          }
        }}
      >
        {operation === "clone" ? (
          <animated.div
            className="bg-white rounded-full h-6 w-6 shadow-md"
            id={`${operation}-switch`}
            style={cloneSwitch ? switchAnimationEnter : switchAnimationExit}
          ></animated.div>
        ) : (
          <animated.div
            className="bg-white rounded-full h-6 w-6 shadow-md"
            id={`${operation}-switch`}
            style={initSwitch ? switchAnimationEnter : switchAnimationExit}
          ></animated.div>
        )}
      </div>
    );
  }

  return (
    <div className="w-11/12 flex items-center justify-center my-10 mx-auto">
      <div className="w-1/2 flex items-center justify-between">
        <div className="w-1/6">{switchComponent("clone")}</div>
        <div className="w-5/6 text-left xl:text-lg lg:text-lg text-base font-sans font-light mx-10">
          Clone from remote
        </div>
      </div>

      <div className="w-1/2 flex items-center justify-between">
        <div className="w-1/6">{switchComponent("init")}</div>
        <div className="w-5/6 text-left xl:text-lg lg:text-lg text-base font-sans font-light mx-10">
          Initialize a new repo
        </div>
      </div>
    </div>
  );
}
