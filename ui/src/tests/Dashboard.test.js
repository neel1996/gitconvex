import { render } from "@testing-library/react";
import React from "react";
import RightPane from "../Components/DashBoard/DashboardPaneComponents/RightPane";

test("Dashboard HC Module test", async () => {
  const hcParams = {
    platform: "Linux",
    gitVersion: "2.26",
  };

  const renderedRightPane = render(<RightPane params={hcParams}></RightPane>);

  const platform = renderedRightPane.container.querySelector(
    "#hc-param__Platform"
  );
  const gitVersion = renderedRightPane.container.querySelector(
    "#hc-param__Git"
  );

  expect(platform.innerHTML).toBe(hcParams.platform);
  expect(gitVersion.innerHTML).toBe(hcParams.gitVersion);
});
