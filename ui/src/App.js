import React, { useReducer } from "react";
import { Route } from "react-router";
import { BrowserRouter, Switch } from "react-router-dom";
import "./App.css";
import Dashboard from "./Components/DashBoard/Dashboard";
import SplashScreen from "./Components/SplashScreen";
import { ContextProvider } from "./context";
import reducer from "./reducer";

export default function App(props) {
  const initialState = {
    hcParams: {},
    presentRepo: [],
    modifiedGitFiles: [],
    globalRepoId: "",
    gitUntrackedFiles: [],
    gitTrackedFiles: [],
  };
  const [state, dispatch] = useReducer(reducer, initialState);

  return (
    <div className="App w-full h-full overflow-hidden">
      <ContextProvider.Provider value={{ state, dispatch }}>
        <BrowserRouter>
          <Switch>
            <Route path="/" exact component={SplashScreen}></Route>
            <Route path="/dashboard" component={Dashboard}></Route>
            <Route
              path="/dashboard/repository"
              exact
              component={Dashboard}
            ></Route>
            <Route
              path="/dashboard/repository/:repoId"
              component={Dashboard}
            ></Route>
          </Switch>
        </BrowserRouter>
      </ContextProvider.Provider>
    </div>
  );
}
