import axios from "axios";
import React, { useContext } from "react";
import { globalAPIEndpoint } from "../../../util/env_config";
import { AddRepoActionTypes } from "./add-new-repo-state/actions";
import { AddRepoContext } from "./add-new-repo-state/addRepoContext";

export default function AddRepoActionButtonsComponent() {
  const { state, dispatch } = useContext(AddRepoContext);

  function setAPIStatusAsFailed() {
    dispatch({
      type: AddRepoActionTypes.SET_ALERT_STATUS,
      payload: "failed",
    });
  }

  function storeRepoAPI() {
    let {
      repoName,
      repoPath,
      cloneSwitch,
      initSwitch,
      cloneURL,
      authMethod,
      httpsAuthInputs,
      sshKeyPath,
    } = state;

    if (repoName && repoPath) {
      if (repoName.match(/[^a-zA-Z0-9-_.\s]/gi)) {
        dispatch({
          type: AddRepoActionTypes.SET_INPUT_INVALID,
          payload: true,
        });
        return;
      }

      let initCheck = false;
      let cloneCheck = false;
      repoPath = repoPath.replace(/\\/gi, "\\\\");

      if (cloneSwitch && !cloneURL) {
        console.log(state);
        setAPIStatusAsFailed();
        return false;
      }

      if (cloneURL.match(/[^a-zA-Z0-9-_.~@#$%:/]/gi)) {
        dispatch({
          type: AddRepoActionTypes.SET_INPUT_INVALID,
          payload: true,
        });
        return;
      }

      let userName = "";
      let password = "";

      if (initSwitch) {
        initCheck = true;
      } else if (cloneSwitch && cloneURL) {
        cloneCheck = true;
        if (authMethod === "https") {
          userName = httpsAuthInputs.userName;
          password = httpsAuthInputs.password;
        }

        if (authMethod === "ssh") {
          if (sshKeyPath) {
            sshKeyPath = sshKeyPath.replaceAll("\\", "\\\\");
          } else {
            dispatch({
              type: AddRepoActionTypes.SET_INPUT_INVALID,
              payload: true,
            });
            return;
          }
        }
      }

      dispatch({
        type: AddRepoActionTypes.SET_LOADING_STATUS,
        payload: true,
      });

      axios({
        url: globalAPIEndpoint,
        method: "POST",
        data: {
          query: `
                  mutation {
                    addRepo(repoName: "${repoName}", repoPath: "${repoPath}", initSwitch: ${initCheck}, cloneSwitch: ${cloneCheck}, repoURL: "${cloneURL}", authOption: "${authMethod}", userName: "${userName}", password: "${password}", sshKeyPath: "${sshKeyPath}"){
                      repoId
                      status
                    }
                  }
                `,
        },
      })
        .then((res) => {
          dispatch({
            type: AddRepoActionTypes.SET_LOADING_STATUS,
            payload: false,
          });

          if (res.data.data && !res.data.error) {
            const { repoId } = res.data.data.addRepo;

            if (repoId && repoId.length > 0) {
              dispatch({
                type: AddRepoActionTypes.SET_ALERT_STATUS,
                payload: "success",
              });
              dispatch({
                type: AddRepoActionTypes.RESET_STATE_VALUES,
                payload: "",
              });
            } else {
              setAPIStatusAsFailed();
            }
          } else {
            setAPIStatusAsFailed();
          }
        })
        .catch((err) => {
          dispatch({
            type: AddRepoActionTypes.SET_LOADING_STATUS,
            payload: false,
          });
          setAPIStatusAsFailed();
          console.log(err);
        });
    } else {
      setAPIStatusAsFailed();
    }
  }

  return (
    <div className="flex justify-between my-5 mx-auto w-11/12">
      <div
        className="cursor-pointer rounded-md block my-2 mx-3 p-3 w-1/2 font-sans font-semibold text-xl bg-red-400 hover:bg-red-500 text-white"
        id="addRepoClose"
        onClick={() => {
          dispatch({
            type: AddRepoActionTypes.CLOSE_FORM,
            payload: true,
          });
        }}
      >
        CLOSE
      </div>
      <div
        className="cursor-pointer rounded-md block my-2 mx-3 p-3 w-1/2 font-sans font-semibold text-xl bg-green-400 hover:bg-green-500 text-white"
        id="addRepoSubmit"
        onClick={() => {
          storeRepoAPI();
        }}
      >
        SUBMIT
      </div>
    </div>
  );
}
