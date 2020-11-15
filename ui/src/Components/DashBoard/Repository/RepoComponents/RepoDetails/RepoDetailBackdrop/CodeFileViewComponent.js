import { LangLine } from "@itassistors/langline";
import axios from "axios";
import * as Prism from "prismjs";
import React, { useEffect, useState } from "react";
import "../../../../../../prism.css";
import InfiniteLoader from "../../../../../Animations/InfiniteLoader";
import { globalAPIEndpoint } from "../../../../../../util/env_config";

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

          console.log(lang);

          if (lang.prismIndicator === "" || !lang.prismIndicator) {
            prism = "go";
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

    return (
      <div className="flex justify-between w-1/2 gap-10 items-center align-middle">
        <div className="text-center w-1/4 font-light font-sans text-xl border-b-4 border-dashed text-gray-700">
          {label}
        </div>
        <div
          className={`"w-11/12 mx-auto rounded p-2 text-center ${bg} ${textColor} ${
            content.length > 10 ? "text-sm" : ""
          } font-semibold"`}
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
        <div className="repo-backdrop--codeview">
          {isInvalidFile ? (
            invalidFileAlert()
          ) : (
            <div className="codeview">
              <div className="codeview--toppane">
                <div className="codeview--language">
                  {languageState
                    ? topPanePills("Language", languageState, {
                        text: "text-pink-500",
                        bg: "bg-pink-200",
                      })
                    : null}
                  {numberOfLines
                    ? topPanePills("Lines", numberOfLines, {
                        text: "text-orange-500",
                        bg: "bg-orange-200",
                      })
                    : null}
                </div>
                {latestCommit ? (
                  <div className="codeview--commits">
                    <div className="codeview--commits--latest">
                      <div className="codeview--commits--latest--label">
                        Latest Commit
                      </div>
                      <div className="codeview--commits--latest--data">
                        {latestCommit}
                      </div>
                    </div>
                  </div>
                ) : null}
              </div>

              {fileDataState && prismIndicator ? (
                <div className="codeview--prismview">
                  <pre className="codeview--prismview--pre">
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
