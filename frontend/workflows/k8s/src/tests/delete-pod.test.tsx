import React from "react";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@clutch-sh/core/src/Theme";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import DeletePod from "../delete-pod";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ThemeProvider>
        <DeletePod resolverType="clutch.k8s.v1.Pod" />
      </ThemeProvider>
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
