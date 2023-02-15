import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import ResizeAutoscalingGroup from "../resize-asg";

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ResizeAutoscalingGroup resolverType="clutch.aws.ec2.v1.AutoscalingGroup" />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
