import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React from "react";

export default function HTTPSAuthHintComponent(props: {
  hideHint: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  return (
    <div
      className="w-full h-full top-0 left-0 right-0 bottom-0 z-40 fixed mx-auto my-auto bg-opacity-80 bg-gray-400"
      id="hintContainer"
    >
      <div
        className="w-14 h-14 mr-5 mt-6 rounded-full bg-red-500 text-white flex justify-center items-center shadow cursor-pointer fixed right-0 top-0"
        onClick={() => {
          props.hideHint(!true);
        }}
      >
        <FontAwesomeIcon
          className="flex text-center text-3xl my-auto"
          icon={["fas", "times"]}
        ></FontAwesomeIcon>
      </div>
      <div className="w-full h-full flex mx-auto my-auto">
        <div className="mx-auto my-auto w-3/4 p-10 bg-white rounded-xl shadow">
          <div className="text-2xl border-b border-dashed font-sans font-semibold mb-6 text-gray-500">
            A FEW THINGS TO CONSIDER
          </div>
          <div className="w-11/12 mx-auto text-base font-sans font-medium my-1 text-center bg-yellow-100 text-yellow-500 p-6 rounded">
            Basic Authentication will not work with a password if 2-Factor
            Authentication is enabled. You may go for
            <span className="text-yellow-700 font-sans font-semibold mx-2 hover:underline">
              <a
                href="https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token"
                target="_blank"
                rel="noreferrer"
              >
                personal access tokens
              </a>
            </span>
            in such cases
          </div>
          <div className="w-11/12 mx-auto my-4 p-6 rounded text-center font-sans font-semibold bg-red-100 text-red-400 text-lg">
            Your password will be stored as an encrypted text in a plain JSON
            file.So go for HTTPS auth with caution!
          </div>
        </div>
      </div>
    </div>
  );
}
