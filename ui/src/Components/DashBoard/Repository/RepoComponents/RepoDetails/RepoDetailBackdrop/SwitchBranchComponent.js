import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../../../../util/env_config";
import axios from "axios";

export default function SwitchBranchComponent({
  repoId,
  branchName,
  closeBackdrop,
  switchReloadView,
}) {
  const [branchError, setBranchError] = useState(false);

  useEffect(() => {
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      data: {
        query: `
            mutation{
              checkoutBranch(repoId: "${repoId}", branchName: "${branchName}")
            }
          `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          switchReloadView();
          closeBackdrop(true);
        }
      })
      .catch((err) => {
        if (err) {
          setBranchError(true);
        }
      });
  }, [branchName, closeBackdrop, repoId, switchReloadView]);

  return (
    <div className="xl:w-3/4 lg:w-3/4 md:w-11/12 sm:w-11/12 repo-backdrop--switchbranch">
      <div className="switchbranch--alert--success">
        Switching to branch -
        <span className="switchbranch--name">{branchName}</span>...
      </div>
      {branchError ? (
        <div className="switchbranch--alert--failed">
          Switching to branch -
          <span className="switchbranch--name">{branchName}</span> Failed!
        </div>
      ) : null}
    </div>
  );
}
