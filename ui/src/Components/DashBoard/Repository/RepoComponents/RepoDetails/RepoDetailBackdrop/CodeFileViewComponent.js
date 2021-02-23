import { LangLine } from "@itassistors/langline";
import axios from "axios";
import * as Prism from "prismjs";
import React, { useEffect, useState } from "react";
import "../../../../../../prism.css";
import InfiniteLoader from "../../../../../Animations/InfiniteLoader";
import { globalAPIEndpoint } from "../../../../../../util/env_config";
import "prismjs/components/prism-markdown";

export default function CodeFileViewComponent(props) {
  const [languageState, setLanguageState] = useState("");
  const [numberOfLines, setNumberOfLines] = useState(0);
  const [latestCommit, setLatestCommit] = useState("");
  const [prismIndicator, setPrismIndicator] = useState("");
  const [highlightedCode, setHighlightedCode] = useState([]);
  const [fileDataState, setFileDataState] = useState([]);
  const [isInvalidFile, setIsInvalidFile] = useState(false);
  const [loading, setLoading] = useState(false);

  const repoId = props.repoId;
  const fileItem = props.fileItem;

  useEffect(() => {
    setLoading(true);
    axios({
      url: globalAPIEndpoint,
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      data: {
        query: `
          query {
            codeFileDetails(repoId: "${repoId}", fileName: "${fileItem}"){
                fileData
              }
          }
        `,
      },
    })
      .then(async (res) => {
        setLoading(false);
        if (res.data.data) {
          const { fileData } = res.data.data.codeFileDetails;

          if (fileData.length === 0) {
            setIsInvalidFile(true);
          }

          let l = new LangLine();
          let lang = l.withFileName(fileItem);
          let prism;

          if (lang.prismIndicator === "" || !lang.prismIndicator) {
            prism = "markdown";
          } else {
            prism = lang.prismIndicator;
          }

          setLatestCommit(props.commitMessage);
          setNumberOfLines(fileData.length);
          setFileDataState(fileData);
          setLanguageState(lang.name);
          setPrismIndicator(prism);

          if (prism) {
            await import("prismjs/components/prism-" + prism + ".js")
              .then(() => {
                const codeHighlight = fileData.map((line) => {
                  return Prism.highlight(line, Prism.languages[prism], prism);
                });
                setHighlightedCode([...codeHighlight]);
              })
              .catch((err) => {
                console.log(err);
                const codeHighlight = fileData.map((line) => {
                  return Prism.highlight(
                    line,
                    Prism.languages["markdown"],
                    "markdown"
                  );
                });
                setHighlightedCode([...codeHighlight]);
              });
          }
        }
      })
      .catch((err) => {
        console.log(err);
        setIsInvalidFile(true);
      });
  }, [repoId, fileItem, props.commitMessage]);

  function topPanePills(label, content, accent) {
    const bg = accent.bg;
    const textColor = accent.text;
    const border = accent.border;

    return (
      <div
        className={`flex justify-between ${
          label === "Language" ? "w-3/4" : "w-1/2"
        } gap-10 items-center align-middle`}
      >
        <div className="text-center w-1/6 font-light font-sans text-xl border-b-4 border-dashed text-gray-700">
          {label}
        </div>
        <div
          className={`${border} mx-auto rounded p-2 text-center ${bg} ${textColor} ${
            content.length > 10 ? "text-sm" : ""
          } font-semibold`}
        >
          {content}
        </div>
      </div>
    );
  }

  function invalidFileAlert() {
    return (
      <div className="w-3/4 mx-auto my-auto p-6 rounded bg-red-200 text-red-600 font-sans text-2xl font-light text-center border-b-8 border-red-400 border-dashed">
        {isInvalidFile ? "File cannot be opened!" : "Loading..."}
      </div>
    );
  }

  return (
    <>
      {loading ? (
        <div className="w-full h-full flex mx-auto my-auto justify-center items-center align-middle">
          <div className="block w-3/4 rounded shadow bg-white p-6 text-center mx-auto font-sans text-2xl my-auto font-light">
            <div>Loading file content...</div>
            <div className="text-center mx-auto flex justify-center my-4">
              <InfiniteLoader loadAnimation={loading}></InfiniteLoader>
            </div>
          </div>
        </div>
      ) : (
        <div className="w-5/6 mx-auto my-auto bg-white rounded-lg p-5">
          {isInvalidFile ? (
            invalidFileAlert()
          ) : (
            <div className="w-full mx-auto my-auto overflow-auto h-full">
              <div className="w-full block mx-auto">
                <div className="w-full flex justify-center items-center mx-auto py-4 border-b-2 border-dashed border-gray-100">
                  {languageState
                    ? topPanePills("Language", languageState, {
                        text: "text-pink-500",
                        bg: "bg-pink-50",
                        border: "border-2 border-dashed border-pink-400",
                      })
                    : null}
                  {numberOfLines
                    ? topPanePills("Lines", numberOfLines, {
                        text: "text-yellow-400",
                        bg: "bg-yellow-50",
                        border: "border-2 border-dashed border-yellow-400",
                      })
                    : null}
                </div>
                {latestCommit ? (
                  <div className="w-11/12 flex justify-around items-center mx-auto my-10 text-center">
                    <div className="w-1/4 font-sans font-semibold text-base text-gray-400">
                      LATEST COMMIT
                    </div>
                    <div className="w-3/4 mx-auto bg-indigo-100 font-sans font-semibold text-base text-center text-indigo-600 p-2 rounded shadow">
                      {latestCommit}
                    </div>
                  </div>
                ) : null}
              </div>

              {fileDataState && prismIndicator ? (
                <div
                  className="w-full shadow p-4 block my-4 rounded-lg bg-gray-700 mx-auto overflow-auto"
                  style={{ height: "700px" }}
                >
                  <pre className="w-11/12 text-white">
                    <code
                      dangerouslySetInnerHTML={{
                        __html: highlightedCode.join("\n"),
                      }}
                    ></code>
                  </pre>
                </div>
              ) : null}
            </div>
          )}
        </div>
      )}
    </>
  );
}
