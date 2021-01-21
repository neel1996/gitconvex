import axios from "axios";
import React, { useEffect, useState } from "react";
import { globalAPIEndpoint } from "../../../util/env_config";

export default function SettingsDataFileComponent(props) {
  const [dbPath, setDbPath] = useState("");
  const [newDbPath, setNewDbPath] = useState("");
  const [dbUpdateFailed, setDbUpdateFailed] = useState(false);

  useEffect(() => {
    const token = axios.CancelToken;
    const source = token.source();

    axios({
      url: globalAPIEndpoint,
      method: "POST",
      cancelToken: source.token,
      data: {
        query: `
          query {
            settingsData{
              settingsDatabasePath
            }
          }
        `,
      },
    })
      .then((res) => {
        if (res.data.data && !res.data.error) {
          const { settingsDatabasePath } = res.data.data.settingsData;

          setDbPath(settingsDatabasePath);
          setNewDbPath(settingsDatabasePath);
        }
      })
      .catch((err) => {
        console.log(err);
      });

    return () => {
      return source.cancel;
    };
  }, []);

  const updateDbFileHandler = () => {
    if (newDbPath) {
      axios({
        url: globalAPIEndpoint,
        method: "POST",
        data: {
          query: `
            mutation {
              updateRepoDataFile(newDbFile: "${newDbPath.toString()}")
            }
          `,
        },
      })
        .then((res) => {
          if (res.data.data && !res.data.error) {
            const updateStatus = res.data.data.updateRepoDataFile;
            if (updateStatus === "DATAFILE_UPDATE_SUCCESS") {
              setDbUpdateFailed(false);
              setDbPath(newDbPath);
            } else {
              setDbUpdateFailed(true);
            }
          } else {
            setDbUpdateFailed(true);
          }
        })
        .catch((err) => {
          console.log("Datafile update error", err);
          setDbUpdateFailed(true);
        });
    }
  };

  function handleDataFileTextChange(e) {
    setNewDbPath(e.target.value);
    setDbUpdateFailed(false);
  }

  return (
    <div className="settings-data">
      <div className="text-xl text-gray-700 font-sans font-semibold">
        Server data file (file which stores repo details)
      </div>
      <div className="my-4">
        <input
          type="text"
          className="p-2 rounded border border-gray-500 bg-gray-200 text-gray-800 w-2/3"
          value={newDbPath}
          onChange={handleDataFileTextChange}
          onClick={() => {
            setDbUpdateFailed(false);
          }}
        ></input>
        <div className="text-justify font-sand font-light text-sm my-4 text-gray-500 italic w-2/3">
          The data file can be updated. The data file must be an accessible JSON
          file with read / write permissions set to it. Also make sure you enter
          the full path for the file
          <pre className="my-2">E.g: /opt/my_data/data-file.json</pre>
        </div>
        {dbPath !== newDbPath && newDbPath ? (
          <div
            className="my-4 text-center p-2 font-sans text-white border-green-400 border-2 bg-green-500 rounded-md cursor-pointer shadow w-1/4 hover:bg-green-600"
            onClick={() => {
              updateDbFileHandler();
              setDbUpdateFailed(false);
            }}
          >
            Update Data file
          </div>
        ) : null}
        {dbUpdateFailed ? (
          <div className="my-2 p-2 rounded border border-red-300 text-red-700 font-sans font-semibold w-2/3 text-center">
            Data file update failed
          </div>
        ) : null}
      </div>
    </div>
  );
}
