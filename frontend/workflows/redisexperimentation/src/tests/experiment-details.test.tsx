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
  const { asFragment } = setup();

  expect(asFragment()).toMatchSnapshot();
});

test("renders correctly with upstream cluster type selection enabled", () => {
  const { utils } = setup();

  expect(utils.getByLabelText(/Upstream Redis Cluster/i)).toBeVisible();
  expect(utils.getByRole("textbox", { name: "Upstream Redis Cluster" })).toBeEnabled();
});

test("renders correctly with environments selection enabled", () => {
  const { utils, asFragment } = setup({
    downstreamClusterTemplate: "",
    environments: [
      {
        value: "staging",
      },
    ],
  });

  expect(utils.getByText(/Environment/i, { selector: "label" })).toBeVisible();
  expect(asFragment().querySelector("#environmentValue-select")).toBeEnabled();
});
