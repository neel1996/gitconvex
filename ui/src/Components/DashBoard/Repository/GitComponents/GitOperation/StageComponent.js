import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../../util/env_config";
import "../../../../styles/GitOperations.css";

export default function StageComponent(props) {
  const { stageComponents, repoId } = props;

  const [allStaged, setAllStaged] = useState(false);
  const [loading, setLodaing] = useState(false);
  const [errorInd, setErrorInd] = useState(false);

  useEffect(() => {
    if (!props) {
      return;
    }
  }, [props]);

  function stageAllChanges() {
    setLodaing(true);
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
          mutation GitConvexMutation{
            stageAllItems(repoId: "${repoId}")
          }
        `,
      },
    })
      .then((res) => {
        setLodaing(false);
        if (res.data.data && !res.data.error) {
          if (res.data.data.stageAllItems === "ALL_STAGED") {
            setAllStaged(true);
          }
        }
      })
      .catch((err) => {
        console.log(err);
        setLodaing(false);
        setErrorInd(true);
      });
  }

  return (
    <>
      <div className="git-ops--stage">
        <>
          {stageComponents.length > 0 && !allStaged ? (
            <>
              <div className="git-ops--stage--header">
                All these changes will be staged:
              </div>
              <div className="overflow-y-auto" style={{ height: "400px" }}>
                {stageComponents &&
                  stageComponents.map((item) => {
                    return (
                      <div className="git-ops--stage--item" key={item}>
                        {item}
                      </div>
                    );
                  })}
              </div>

              {errorInd ? (
                <div className="git-ops--stage--alert--failure">
                  Staging Failed!
                </div>
              ) : null}
              {loading ? (
                <div className="git-ops--stage--alert--progress">
                  Staging in prgoress...
                </div>
              ) : (
                <div
                  className="git-ops--stage--btn"
                  onClick={() => {
                    stageAllChanges();
                  }}
                >
                  Confirm Staging
                </div>
              )}
            </>
          ) : (
            <div className="git-ops--stage--alert--nochange">
              No Changes for staging...
            </div>
          )}
          {allStaged ? (
            <div className="git-ops--stage--alert--success">
              All changes Staged!
            </div>
          ) : null}
        </>
      </div>
    </>
  );
}
