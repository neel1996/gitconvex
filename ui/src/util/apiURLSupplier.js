export const getAPIURL = (httpMode, endpoint, port) => {
    return `${httpMode}://${window.location.hostname}:${port}/${endpoint}`
}