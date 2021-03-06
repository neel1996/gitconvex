import { AddRepoStateType } from "./addRepoContext";

export enum AddRepoActionTypes {
  SET_REPO_NAME = "SET_REPO_NAME",
  SET_REPO_PATH = "SET_REPO_PATH",
  SET_CLONE_URL = "SET_CLONE_URL",
  SET_CLONE_SWITCH = "SET_CLONE_SWITCH",
  SET_INIT_SWITCH = "SET_INIT_SWITCH",
  SET_AUTH_OPTION = "SET_AUTH_OPTION",
  SET_ALERT_STATUS = "SET_ALERT_STATUS",
  SET_INPUT_INVALID = "SET_INPUT_INVALID",
  SET_LOADING_STATUS = "SET_LOADING_STATUS",
  SET_SSH_KEY_PATH = "SET_SSH_KEY_PATH",
  SET_HTTPS_AUTH_INPUTS = "SET_HTTPS_AUTH_INPUTS",
  RESET_STATE_VALUES = "RESET_STATE_VALUES",
  CLOSE_FORM = "CLOSE_FORM",
}

export type AddRepoAction = {
  type: AddRepoActionTypes;
  payload: AddRepoStateType;
};
