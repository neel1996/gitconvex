import React from "react";
import { render } from "@testing-library/react";
import App from "../App";

test("renders learn react link", () => {
  const RenderedApp = render(<App />);
  expect(RenderedApp).toBeTruthy();
});
