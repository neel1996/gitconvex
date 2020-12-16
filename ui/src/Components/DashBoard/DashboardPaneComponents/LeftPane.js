import { library } from "@fortawesome/fontawesome-svg-core";
import { far } from "@fortawesome/free-regular-svg-icons";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useContext } from "react";
import { NavLink } from "react-router-dom";
import { ADD_FORM_CLOSE } from "../../../actionStore";
import { ContextProvider } from "../../../context";
import "../../styles/LeftPane.css";

export default function LeftPane(props) {
  library.add(far, fas);
  const { dispatch } = useContext(ContextProvider);
  const menuItemParams = [
    {
      link: "/dashboard/repository",
      icon: (
        <FontAwesomeIcon
          icon={["far", "folder"]}
          className="text-3xl text-gray-500"
        ></FontAwesomeIcon>
      ),
      label: "Repositories",
    },
    {
      link: "/dashboard/compare",
      icon: (
        <FontAwesomeIcon
          icon={["fas", "object-group"]}
          className="text-3xl text-gray-500"
        ></FontAwesomeIcon>
      ),
      label: "Compare",
    },
    {
      link: "/dashboard/settings",
      icon: (
        <FontAwesomeIcon
          icon={["fas", "cog"]}
          className="text-3xl text-gray-500"
        ></FontAwesomeIcon>
      ),
      label: "Settings",
    },
    {
      link: "/dashboard/help",
      icon: (
        <FontAwesomeIcon
          icon={["far", "question-circle"]}
          className="text-3xl text-gray-500"
        ></FontAwesomeIcon>
      ),
      label: "Help",
    },
  ];

  return (
    <div className="dashboard--leftpane block xl:w-1/4 lg:w-1/3 md:w-1/6 sm:w-1/6 w-1/6">
      <div
        className="leftpane--logo"
        onClick={(event) => {
          dispatch({ type: ADD_FORM_CLOSE, payload: true });
          props.parentProps.history.push("/dashboard");
        }}
      >
        <div className="dashboard-leftpane__logo"></div>
        <div className="font-mono xl:text-3xl lg:text-2xl md:text-3xl sm:text-2xl p-4 xl:block lg:block md:hidden sm:hidden hidden">
          <span className="font-bold mx-2 border-b-4 border-pink-400">Git</span>
          Convex
        </div>
      </div>
      <div className="menu xl:mt-32 lg:mt-24 md:mt-48 sm:mt-56 mt-56 cursor-pointer block items-center align-middle">
        {menuItemParams.map((entry) => {
          return (
            <NavLink
              to={`${entry.link}`}
              exact
              activeClassName="bg-gray-200"
              className="menu--link xl:justify-between lg:justify-between md:justify-center sm:justify-center justify-center xl:my-0 lg:my-0 md:my-6 sm:my-6 my-6"
              key={entry.label}
            >
              <div className="menu--items sm:text-center">
                <div className="menu--items__icon text-sm w-1/6">
                  {entry.icon}
                </div>
                <div className="menu--items__label w-5/6 xl:text-2xl lg:text-2xl md:text-xl block xl:block lg:block md:hidden sm:hidden">
                  {entry.label}
                </div>
              </div>
            </NavLink>
          );
        })}
      </div>
    </div>
  );
}
