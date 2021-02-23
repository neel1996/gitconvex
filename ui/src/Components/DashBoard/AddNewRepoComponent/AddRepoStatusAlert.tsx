import { useContext } from "react";
import {
  AddRepoContext,
  statusType,
} from "./add-new-repo-state/addRepoContext";

export default function AddRepoStatusAlert(props: {
  status: statusType;
  message?: string;
}) {
  const { state } = useContext(AddRepoContext);

  return (
    <>
      {props.status === "success" || state.alertStatus === "success" ? (
        <div className="w-3/4 font-sans bg-green-50 border-green-300 rounded-lg border-dotted border-4 block font-light text-lg my-6 mx-auto p-2 shadow-sm text-center text-green-700">
          {props.message
            ? props.message
            : "New repo has been added for tracking"}
        </div>
      ) : null}
      {props.status === "failed" || state.alertStatus === "failed" ? (
        <div className="w-3/4 font-sans bg-red-50 border-red-300 rounded-lg border-dotted border-4 block text-lg my-6 mx-auto p-2 shadow-sm text-center text-red-700 font-semibold">
          {props.message ? props.message : "Process failed! Please try again"}
        </div>
      ) : null}
    </>
  );
}
