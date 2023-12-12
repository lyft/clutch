import React from "react";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@clutch-sh/core/src/Theme";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import ResizeHPA from "../resize-hpa";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ThemeProvider>
        <ResizeHPA resolverType="clutch.aws.k8s.v1.HPA" />
      </ThemeProvider>
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
