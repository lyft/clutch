import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import TerminateInstance from "../terminate-instance";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <TerminateInstance resolverType="clutch.aws.ec2.v1.Instance" />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
