import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import ResizeHPA from "../resize-hpa";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ResizeHPA resolverType="clutch.aws.k8s.v1.HPA" />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
