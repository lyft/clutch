import React from "react";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@clutch-sh/core/src/Theme";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import ScaleResources from "../scale-resources";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ThemeProvider>
        <ScaleResources resolverType="clutch.k8s.v1.Deployment" />
      </ThemeProvider>
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
