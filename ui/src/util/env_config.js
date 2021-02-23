// Global ENV file for all environment variables

export const CURRENT_VERSION = "2.0.2";

//PORTS LIST
export const PORT_GLOBAL_API =
  process.env.NODE_ENV === "development" ? 9001 : window.location.port;

//CONFIG LIST
export const CONFIG_HTTP_MODE = "http";
export const API_GLOBAL_GQL = "gitconvexapi";
export const globalAPIEndpoint = `${CONFIG_HTTP_MODE}://${window.location.hostname}:${PORT_GLOBAL_API}/${API_GLOBAL_GQL}`;