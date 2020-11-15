// Global ENV file for all environment variables

export const CURRENT_VERSION = "2.0.0";

//PORTS LIST

export const PORT_GLOBAL_API =
  process.env.NODE_ENV === "development" ? 9001 : window.location.port;

//CONFIG LIST

export const CONFIG_HTTP_MODE = "http";

//API LIST

export const API_GLOBAL_GQL = "gitconvexapi";
export const globalAPIEndpoint = `${CONFIG_HTTP_MODE}://${window.location.hostname}:${PORT_GLOBAL_API}/${API_GLOBAL_GQL}`;

// ROUTED FOR GLOBAL API

export const ROUTE_HEALTH_CHECK = "HEALTH_CHECK";
export const ROUTE_FETCH_REPO = "FETCH_REPO";
export const ROUTE_ADD_REPO = "ADD_REPO";
export const ROUTE_REPO_DETAILS = "REPO_DETAILS";
export const ROUTE_REPO_TRACKED_DIFF = "REPO_TRACKED_DIFF";
export const ROUTE_REPO_FILE_DIFF = "REPO_FILE_DIFF";
export const ROUTE_REPO_COMMIT_LOGS = "COMMIT_LOGS";
export const ROUTE_COMMIT_FILES = "COMMIT_FILES";
export const ROUTE_GIT_STAGED_FILES = "GIT_STAGED_FILES";
export const ROUTE_GIT_UNPUSHED_COMMITS = "GIT_UNPUSHED_COMMITS";
export const ROUTE_SETTINGS_DBPATH = "SETTINGS_DBPATH";
export const ROUTE_SETTINGS_REPODETAILS = "SETTINGS_REPODETAILS";
export const ROUTE_SETTINGS_PORT = "SETTINGS_PORT";
export const GIT_FOLDER_CONTENT = "GIT_FOLDER_CONTENT";
export const CODE_FILE_VIEW = "CODE_FILE_VIEW";
export const BRANCH_COMPARE = "BRANCH_COMPARE";
export const COMMIT_COMPARE = "COMMIT_COMPARE";
