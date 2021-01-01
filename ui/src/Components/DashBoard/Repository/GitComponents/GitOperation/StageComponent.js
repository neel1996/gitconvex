import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../../util/env_config";
import "../../../../styles/GitOperations.css";

export default function StageComponent(props) {
  const { stageComponents, repoId } = props;

  const [allStaged, setAllStaged] = useState(false);
  const [loading, setLoading] = useState(false);
  const [errorInd, setErrorInd] = useState(false);

  useEffect(() => {
    if (!props) {
      return;
    }
  }, [props]);

  function stageAllChanges() {
    setLoading(true);
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          mutation {
            stageAllItems(repoId: "${repoId}")
          }
        `,
      },
    })
      .then((res) => {
        setLoading(false);
        if (res.data.data && !res.data.error) {
          if (res.data.data.stageAllItems === "ALL_STAGED") {
            setAllStaged(true);
          }
        }
      })
      .catch((err) => {
        console.log(err);
        setLoading(false);
        setErrorInd(true);
      });
  }

  return (
    <>
      <div className="w-11/12 xl:w-3/4 lg:w-3/4 mx-auto my-auto rounded-lg shadow p-10 bg-gray-50">
        <>
          {stageComponents && stageComponents.length > 0 && !allStaged ? (
            <>
              <div className="font-sans font-bold text-2xl text-gray-700 my-2 mx-4">
                {stageComponents.length === 1 ? (
                  <span>One change will be staged</span>
                ) : (
                  <span>{stageComponents.length} changes will be staged:</span>
                )}
              </div>
              <div
                className="mx-6 my-4 overflow-y-auto"
                style={{ height: "400px" }}
              >
                {stageComponents &&
                  stageComponents.map((item) => {
                    return (
                      <div
                        className="font-sans my-2 text-gray-600 border-b"
                        key={item}
                      >
                        {item}
                      </div>
                    );
                  })}
              </div>

              {errorInd ? (
                <div className="w-full mx-auto my-2 bg-red-200 border-b-2 border-red-400 rounded rounded-b-lg p-3 text-center text-xl text-red-600 font-semibold font-sans">
                  Staging Failed!
                </div>
              ) : null}
              {loading ? (
                <div className="git-ops--stage--alert--progress">
                  Staging in prgoress...
                </div>
              ) : (
                <div
                  className="w-full my-2 p-3 rounded shadow bg-green-500 font-sans font-semibold text-center text-white text-xl cursor-pointer hover:bg-green-600"
                  onClick={() => {
                    stageAllChanges();
                  }}
                >
                  STAGE ALL COMPONENTS
                </div>
              )}
            </>
          ) : null}
          {allStaged ? (
            <div className="shadow rounded text-center text-green-700 p-6 bg-green-200 font-sans font-semibold text-2xl border-b-4 border-dashed border-green-400 rounded-b-xl">
              All changes Staged!
            </div>
          ) : null}
        </>
      </div>
    </>
  );
}
