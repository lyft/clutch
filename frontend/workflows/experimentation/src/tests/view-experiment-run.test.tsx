import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import ViewExperimentRun from "../view-experiment-run";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ViewExperimentRun />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
