import { render } from "@testing-library/react";
import React from "react";
import RightPane from "../Components/DashBoard/DashboardPaneComponents/RightPane";

test("Dashboard HC Module test", async () => {
  const hcParams = {
    platform: "Linux",
    gitVersion: "2.26",
    nodeVersion: "14.2",
  };

  const renderedRightPane = render(<RightPane params={hcParams}></RightPane>);

  const platform = renderedRightPane.container.querySelector(
    "#hc-param__Platform"
  );
  const gitVersion = renderedRightPane.container.querySelector(
    "#hc-param__Git"
  );
  const nodeVerrsion = renderedRightPane.container.querySelector(
    "#hc-param__Node"
  );

  expect(platform.innerHTML).toBe(hcParams.platform);
  expect(gitVersion.innerHTML).toBe(hcParams.gitVersion);
  expect(nodeVerrsion.innerHTML).toBe(hcParams.nodeVersion);
});
