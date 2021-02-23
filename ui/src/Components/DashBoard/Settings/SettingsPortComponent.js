import React from "react";
import { useState, useEffect } from "react";
import { globalAPIEndpoint } from "../../../util/env_config";
import axios from "axios";

export default function SettingsPortComponent() {
  const [port, setPort] = useState(0);
  const [portUpdateFailed, setPortUpdateFailed] = useState(false);

  useEffect(() => {
    const token = axios.CancelToken;
    const source = token.source();

    axios({
      url: globalAPIEndpoint,
      cancelToken: source.token,
      method: "POST",
      data: {
        query: `
            query {
              settingsData{
                settingsPortDetails
              }
            }
          `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          const { settingsPortDetails } = res.data.data.settingsData;
          setPort(settingsPortDetails);
        }
      })
      .catch((err) => {
        console.log(err);
      });

    return () => {
      return source.cancel;
    };
  }, []);

  function portUpdateHandler() {
    if (port) {
      axios({
        url: globalAPIEndpoint,
        method: "POST",
        data: {
          query: `
            mutation {
              settingsEditPort(newPort: "${port}")
            }
          `,
        },
      })
        .then((res) => {
          const { settingsEditPort } = res.data.data;
          if (settingsEditPort === "PORT_UPDATED") {
            window.location.reload();
          } else {
            portUpdateFailed(true);
          }
        })
        .catch((err) => {
          console.log(err);
          setPortUpdateFailed(true);
        });
    }
  }

  return (
    <div className="my-2 mx-auto">
      <div className="text-xl font-sans text-gray-800 my-2">
        Active Gitconvex port
      </div>
      <div className="flex my-4">
        <input
          type="text"
          className="p-2 rounded border border-gray-500 bg-gray-200 text-gray-800 xl:w-1/2 lg:w-1/3 md:w-1/2 sm:w-1/2 w-1/2"
          value={port}
          onChange={(event) => {
            setPort(event.target.value);
          }}
        ></input>
        <div
          className="p-2 text-center mx-4 rounded border text-white bg-indigo-500 xl:w-1/6 lg:w-1/6 md:w-1/5 sm:w-1/4 w-1/4 hover:bg-indigo-600 cursor-pointer"
          onClick={() => {
            portUpdateHandler();
          }}
        >
          Update Port
        </div>
      </div>
      <div className="text-justify font-sand font-light text-sm my-4 text-gray-500 italic w-2/3">
        Make sure to restart the app and to change the port in the URL after
        updating it
      </div>
      {portUpdateFailed ? (
        <div className="my-2 p-2 rounded border border-red-300 text-red-700 font-sans font-semibold w-1/2 text-center">
          Port update failed
        </div>
      ) : null}
    </div>
  );
}
