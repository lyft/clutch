import React from "react";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@clutch-sh/core";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import ViewExperimentRun from "../view-experiment-run";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ThemeProvider>
        <ViewExperimentRun />
      </ThemeProvider>
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
