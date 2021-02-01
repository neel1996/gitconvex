import React, { useContext, useEffect } from "react";
import InfiniteLoader from "../../Animations/InfiniteLoader";
import { AddRepoContext } from "./add-new-repo-state/addRepoContext";
import AddRepoFormComponent from "./AddRepoFormComponent";

export default function AddNewRepoComponent(props: {
  formEnable: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  const { state } = useContext(AddRepoContext);
  const loading = state.isLoading;

  useEffect(() => {
    if (state.closeForm) {
      props.formEnable(false);
    }
  }, [props, state.closeForm]);

  return (
    <div
      className={`border-gray-200 rounded-lg border-2 block justify-center my-20 mx-auto p-6 text-center shadow xl:w-7/12 lg:w-2/3 md:w-3/4 sm:w-11/12 w-11/12 ${
        loading ? "border-dashed border-2" : ""
      }`}
    >
      {loading ? (
        <>
          <div className="font-sans font-semibold text-xl text-center text-gray-600">
            Repo setup in progress...
          </div>
          <div className="flex justify-center my-6 mx-auto text-center">
            <InfiniteLoader loadAnimation={loading}></InfiniteLoader>
          </div>
        </>
      ) : (
        <AddRepoFormComponent></AddRepoFormComponent>
      )}
    </div>
  );
}
