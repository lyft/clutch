import React from "react";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@clutch-sh/core/src/Theme";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import { ExperimentDetails } from "../start-experiment";

const setup = (props?) => {
  const utils = render(
    <BrowserRouter>
      <ThemeProvider>
        <ExperimentDetails environments={[]} onStart={() => {}} {...props} />
      </ThemeProvider>
    </BrowserRouter>
  );

  const { asFragment } = utils;

  return { utils, asFragment };
};

test("renders correctly", () => {
  const { asFragment } = setup({
    upstreamClusterTypeSelectionEnabled: false,
  });

  expect(asFragment()).toMatchSnapshot();
});

test("renders correctly with upstream cluster type selection enabled", () => {
  const { asFragment } = setup({
    upstreamClusterTypeSelectionEnabled: true,
  });

  expect(asFragment()).toMatchSnapshot();
});

test("renders correctly with environments", () => {
  const { utils, asFragment } = setup({
    upstreamClusterTypeSelectionEnabled: false,
    environments: [{ value: "staging" }],
  });
  expect(utils.getByText(/Environment/i, { selector: "label" })).toBeVisible();
  expect(asFragment().querySelector("#environmentValue-select")).toBeEnabled();
});
