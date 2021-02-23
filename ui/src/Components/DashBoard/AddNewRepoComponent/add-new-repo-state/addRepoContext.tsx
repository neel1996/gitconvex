import React from "react";

export type authMethodType = "noauth" | "ssh" | "https";
export type statusType = "success" | "failed" | "";

export interface AddRepoStateType {
  repoName: string;
  repoPath: string;
  isLoading: boolean;
  cloneURL?: string;
  cloneSwitch: boolean;
  initSwitch: boolean;
  alertStatus: statusType;
  isInputInvalid: boolean;
  authMethod: authMethodType;
  sshKeyPath?: string;
  httpsAuthInputs?: {
    userName: string;
    password: string;
  };
  closeForm: boolean;
}

export const AddRepoState: AddRepoStateType = {
  repoName: "",
  repoPath: "",
  isLoading: false,
  cloneURL: "",
  cloneSwitch: false,
  initSwitch: false,
  alertStatus: "",
  isInputInvalid: false,
  authMethod: "noauth",
  sshKeyPath: "",
  httpsAuthInputs: {
    userName: "",
    password: "",
  },
  closeForm: false,
};

const AddRepoContext = React.createContext<AddRepoStateType | any>(
  AddRepoState
);
export { AddRepoContext };
