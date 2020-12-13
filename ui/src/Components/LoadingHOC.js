import React from "react";
import InfiniteLoader from "./Animations/InfiniteLoader";

export default function LoadingHOC(props) {
  return (
    <div className="flex mx-auto my-auto w-1/2 h-full bg-white">
      <div className="flex w-full mx-auto my-auto justify-center">
        <div className="block justify-center w-full h-full mx-auto my-auto">
          <div className="p-6 my-auto mx-auto font-sans font-light text-3xl shadow-lg rounded-lg border-b-4 border-dashed border-gray-300 text-center">
            {props.message}
          </div>
          <div className="flex mx-auto my-6 text-center justify-center">
            <InfiniteLoader loadAnimation={props.loading}></InfiniteLoader>
          </div>
        </div>
      </div>
    </div>
  );
}
