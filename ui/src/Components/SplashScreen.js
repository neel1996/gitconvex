import axios from "axios";
import React, { useContext, useEffect, useState } from "react";
import { HC_PARAM_ACTION } from "../actionStore";
import { ContextProvider } from "../context";
import { globalAPIEndpoint } from "../util/env_config";
import "./SplashScreen.css";

export default function SplashScreen(props) {
  const [showAlert, setShowAlert] = useState(false);
  const [hcCheck, setHcCheck] = useState(false);
  const { dispatch } = useContext(ContextProvider);

  useEffect(() => {
    const apiURL = globalAPIEndpoint;
    axios({
      url: apiURL,
      method: "POST",
      data: {
        query: `
          query GitConvexAPI{
            gitConvexApi(route:"HEALTH_CHECK"){
              healthCheck{
                osCheck
                gitCheck
                nodeCheck
              }
            }
          }
        `,
      },
    })
      .then((res) => {
        if (res.data.data) {
          const {
            osCheck,
            gitCheck,
            nodeCheck,
          } = res.data.data.gitConvexApi.healthCheck;

          dispatch({
            type: HC_PARAM_ACTION,
            payload: {
              osCheck,
              gitCheck,
              nodeCheck,
            },
          });
          setHcCheck(true);
        } else {
          setShowAlert(true);
        }
      })
      .catch((err) => {
        console.log(err);
        setShowAlert(true);
      });
  }, [dispatch]);

  return (
    <>
      {!hcCheck ? (
        <div className="w-64 h-full justify-center mx-auto flex my-auto align-center items-center">
          <div className="splash-logo w-64 h-full justify-center mx-auto flex my-auto align-center items-center">
            <div className="p-5 shadow-md border-l-4 border-t-4 border-blue-100 rounded-lg block">
              <div className="splash-logo__image"></div>
            </div>
            <div className="logo-label my-3 p-3 text-center block">
              <div className="logo-label__title block text-6xl border-solid border-b-4 border-pink-200">
                <span className="logo-label__title-first font-sans font-bold mx-3">
                  Git
                </span>
                <span>Convex</span>
              </div>
              <div className="block font-mono my-2">
                A visualizer for your git repo
              </div>
            </div>
          </div>

          {showAlert ? (
            <div className="fixed bottom-0 left-0 right-0 w-full p-3 rounded-lg text-center font-sans bg-red-200 border-red-900">
              <div className="h1 text-3xl p-2 m-2 text-red-800">
                Server Unreachable
              </div>
              <p>
                The server is not reachable.Please check if the server module is
                running.
              </p>
            </div>
          ) : null}
        </div>
      ) : (
        props.history.push("/dashboard")
      )}
    </>
  );
}
