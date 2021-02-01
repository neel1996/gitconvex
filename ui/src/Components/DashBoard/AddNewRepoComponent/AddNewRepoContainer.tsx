import React, { useReducer } from "react";
import { addRepoReducer } from "./add-new-repo-state/addRepoReducer";
import {
  AddRepoContext,
  AddRepoState,
} from "./add-new-repo-state/addRepoContext";
import AddNewRepoComponent from "./AddNewRepoComponent";

export default function AddNewRepoContainer(props: {
  formEnable: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const [state, dispatch] = useReducer(addRepoReducer, AddRepoState);

  return (
    <AddRepoContext.Provider value={{ state, dispatch }}>
      <AddNewRepoComponent formEnable={props.formEnable}></AddNewRepoComponent>
    </AddRepoContext.Provider>
  );
}
