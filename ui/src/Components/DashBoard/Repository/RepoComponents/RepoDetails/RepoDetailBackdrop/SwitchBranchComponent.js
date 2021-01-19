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
          const status = res.data.data.checkoutBranch;
          console.log(status);
          if (status === "CHECKOUT_FAILED") {
            setBranchError(true);
          } else {
            switchReloadView();
            closeBackdrop(true);
          }
        } else {
          setBranchError(true);
        }
      })
      .catch((err) => {
        if (err) {
          setBranchError(true);
        }
      });
  }, [branchName, closeBackdrop, repoId, switchReloadView]);

  return (
    <div className="xl:w-3/4 lg:w-3/4 md:w-11/12 sm:w-11/12 w-11/12 mx-auto my-auto p-6 rounded-md bg-gray-100">
      <div className="bg-blue-100 p-2 border-indigo-400 rounded shadow">
        Switching to branch -
        <span className="font-sans text-xl font-semibold">{branchName}</span>...
      </div>
      {branchError ? (
        <div className="bg-red-100 p-2 border-red-400 rounded shadow">
          Switching to branch -
          <span className="font-sans text-xl font-semibold">{branchName}</span>
          Failed!
        </div>
      ) : null}
    </div>
  );
}
