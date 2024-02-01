import React from "react";
import { BrowserRouter } from "react-router-dom";
import { ThemeProvider } from "@clutch-sh/core/src/Theme";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import TerminateInstance from "../terminate-instance";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ThemeProvider>
        <TerminateInstance resolverType="clutch.aws.ec2.v1.Instance" />
      </ThemeProvider>
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
