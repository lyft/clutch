import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import ScaleResources from "../scale-resources";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ScaleResources resolverType="clutch.k8s.v1.Deployment" />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
