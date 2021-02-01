import { AddRepoAction, AddRepoActionTypes } from "./actions";
import { AddRepoStateType } from "./addRepoContext";

export const addRepoReducer = (
  state: AddRepoStateType | any,
  action: AddRepoAction | any
) => {
  const { type, payload } = action;

  switch (type) {
    case AddRepoActionTypes.SET_REPO_NAME: {
      return { ...state, repoName: payload };
    }
    case AddRepoActionTypes.SET_REPO_PATH: {
      return { ...state, repoPath: payload };
    }
    case AddRepoActionTypes.SET_CLONE_URL: {
      return { ...state, cloneURL: payload };
    }
    case AddRepoActionTypes.SET_CLONE_SWITCH: {
      return { ...state, cloneSwitch: payload };
    }
    case AddRepoActionTypes.SET_INIT_SWITCH: {
      return { ...state, initSwitch: payload };
    }
    case AddRepoActionTypes.SET_ALERT_STATUS: {
      return { ...state, alertStatus: payload };
    }
    case AddRepoActionTypes.SET_INPUT_INVALID: {
      return { ...state, alertStatus: "failed", isInputInvalid: true };
    }
    case AddRepoActionTypes.SET_LOADING_STATUS: {
      return { ...state, isLoading: payload };
    }
    case AddRepoActionTypes.SET_AUTH_OPTION: {
      return { ...state, authMethod: payload };
    }
    case AddRepoActionTypes.SET_HTTPS_AUTH_INPUTS: {
      return {
        ...state,
        httpsAuthInputs: {
          ...payload,
        },
      };
    }
    case AddRepoActionTypes.RESET_STATE_VALUES: {
      return {
        ...state,
        repoName: "",
        repoPath: "",
        isLoading: false,
        cloneURL: "",
        cloneSwitch: false,
        initSwitch: false,
        isInputInvalid: false,
        authMethod: "noauth",
        httpsAuthInputs: {
          userName: "",
          password: "",
        },
        closeForm: false,
      };
    }
    case AddRepoActionTypes.CLOSE_FORM: {
      return { ...state, closeForm: true };
    }
    default: {
      return { ...state };
    }
  }
};
