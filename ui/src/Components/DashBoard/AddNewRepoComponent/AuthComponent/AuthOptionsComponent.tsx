import { useContext } from "react";
import { AddRepoActionTypes } from "../add-new-repo-state/actions";
import {
  AddRepoContext,
  authMethodType,
} from "../add-new-repo-state/addRepoContext";

export default function AuthOptionsComponent() {
  const { dispatch } = useContext(AddRepoContext);

  const authRadio: { key: authMethodType; label: string }[] = [
    {
      key: "noauth",
      label: "No Authentication",
    },
    {
      key: "ssh",
      label: "SSH Authentication",
    },
    {
      key: "https",
      label: "HTTPS Authentication",
    },
  ];

  return (
    <div className="flex gap-4 justify-center mx-auto items-center align-middle">
      {authRadio.map((item) => {
        return (
          <div className="flex gap-4 items-center align-middle" key={item.key}>
            <input
              type="radio"
              name="authRadio"
              id={item.key}
              value={item.key}
              onChange={(e) => {
                dispatch({
                  type: AddRepoActionTypes.SET_AUTH_OPTION,
                  payload: e.currentTarget.value,
                });
              }}
            ></input>
            <label
              htmlFor={item.key}
              className="font-sans text-sm font-light cursor-pointer"
            >
              {item.label}
            </label>
          </div>
        );
      })}
    </div>
  );
}
