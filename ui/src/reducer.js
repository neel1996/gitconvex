import {
  HC_PARAM_ACTION,
  HC_DONE_SWITCH,
  PRESENT_REPO,
  GIT_TRACKED_FILES,
  GIT_GLOBAL_REPOID,
  GIT_ACTION_TRACKED_FILES,
  GIT_ACTION_UNTRACKED_FILES,
  DELETE_PRESENT_REPO,
  ADD_FORM_CLOSE,
} from "./actionStore";

export default function reducer(state, action) {
  switch (action.type) {
    case HC_DONE_SWITCH:
      return {
        ...state,
        hcDone: action.payload,
      };
    case HC_PARAM_ACTION:
      const { osCheck, gitCheck, nodeCheck } = action.payload;

      const parseValue = (payload) => {
        return JSON.parse(JSON.parse(JSON.stringify(payload))).message;
      };

      localStorage.setItem("OS_TYPE", parseValue(osCheck));
      localStorage.setItem("GIT_VERSION", parseValue(gitCheck));
      localStorage.setItem("NODE_VERSION", parseValue(nodeCheck));

      return {
        ...state,
        hcParams: {
          osCheck: parseValue(osCheck),
          gitCheck: parseValue(gitCheck),
          nodeCheck: parseValue(nodeCheck),
        },
      };
    case PRESENT_REPO:
      return {
        ...state,
        presentRepo: [...state.presentRepo, action.payload],
      };
    case ADD_FORM_CLOSE:
      return {
        ...state,
        shouldAddFormClose: action.payload,
      };
    case DELETE_PRESENT_REPO:
      return {
        ...state,
        presentRepo: [...action.payload],
      };
    case GIT_TRACKED_FILES:
      state.modifiedGitFiles = [];
      return {
        ...state,
        modifiedGitFiles: [...state.modifiedGitFiles, action.payload],
      };
    case GIT_GLOBAL_REPOID:
      state.globalRepoId = "";
      return {
        ...state,
        globalRepoId: action.payload,
      };
    case GIT_ACTION_TRACKED_FILES:
      state.gitTrackedFiles = [];
      return {
        ...state,
        gitTrackedFiles: [...state.gitTrackedFiles, action.payload],
      };
    case GIT_ACTION_UNTRACKED_FILES:
      state.gitUntrackedFiles = [];
      return {
        ...state,
        gitUntrackedFiles: [...state.gitUntrackedFiles, action.payload],
      };
    default:
      return {
        ...state,
      };
  }
}
