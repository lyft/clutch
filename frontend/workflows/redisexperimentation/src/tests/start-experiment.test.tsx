import React from "react";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@clutch-sh/core/src/Theme";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import { StartExperiment } from "../start-experiment";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ThemeProvider>
        <StartExperiment heading="testing" />
      </ThemeProvider>
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
