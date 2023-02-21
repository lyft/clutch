import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import { StartExperiment } from "../start-experiment";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <StartExperiment heading="testing" />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
