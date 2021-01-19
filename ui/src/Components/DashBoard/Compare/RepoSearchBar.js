import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import debounce from "lodash.debounce";
import React, { useRef, useState } from "react";
import SearchRepoCards from "./SearchRepoCards";

export default function RepoSearchBar(props) {
  const debounceRef = useRef(
    debounce(
      function () {
        setSelectedRepo("");
        setToggleSearchResult(true);
      },
      500,
      { maxWait: 1500 }
    )
  ).current;

  const [toggleSearchResult, setToggleSearchResult] = useState(false);
  const [searchQueryState, setSearchQueryState] = useState("");
  const [selectedRepo, setSelectedRepo] = useState("");

  const searchTextRef = useRef();

  function setSelectedRepoHandler(repo) {
    setSelectedRepo(repo);
    searchTextRef.current.value = "";
  }

  return (
    <>
      <div className="rounded-md border flex items-center justify-between mx-auto mt-4 shadow w-11/12">
        <div className="w-11/12 rounded-r-md">
          <input
            type="text"
            ref={searchTextRef}
            className="w-full p-3 outline-none text-lg font-light font-sans"
            placeholder="Enter repo name to search"
            onChange={(event) => {
              setSearchQueryState(event.target.value);
              debounceRef();
            }}
          />
        </div>
        <div
          className="bg-gray-200 rounded-r-md cursor-pointer py-4 px-6 text-center hover:bg-gray-400"
          onClick={() => {
            debounceRef();
          }}
        >
          <FontAwesomeIcon
            icon={["fas", "search"]}
            className="text-3xl text-gray-600"
          ></FontAwesomeIcon>
        </div>
      </div>
      {toggleSearchResult && searchQueryState && !selectedRepo ? (
        <div className="w-11/12 mx-auto rounded-b-md p-3 border">
          <SearchRepoCards
            searchQuery={searchQueryState}
            setSelectedRepoHandler={setSelectedRepoHandler}
          ></SearchRepoCards>
        </div>
      ) : null}

      {selectedRepo ? props.activateCompare(selectedRepo) : null}
    </>
  );
}
