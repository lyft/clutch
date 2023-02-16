import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import DeletePod from "../delete-pod";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <DeletePod resolverType="clutch.k8s.v1.Pod" />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
