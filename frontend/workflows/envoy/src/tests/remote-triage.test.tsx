import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import RemoteTriage from "../remote-triage";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <RemoteTriage />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
