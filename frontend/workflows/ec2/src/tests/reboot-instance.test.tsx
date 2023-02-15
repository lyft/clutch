import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import RebootInstance from "../reboot-instance";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <RebootInstance resolverType="clutch.aws.ec2.v1.Instance" />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
