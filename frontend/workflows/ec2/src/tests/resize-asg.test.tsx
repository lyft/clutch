import React from "react";
import { shallow } from "enzyme";

import ResizeAutoscalingGroup from "../resize-asg";

describe("Resize Autoscaling Group workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(
      <ResizeAutoscalingGroup resolverType="clutch.aws.ec2.v1.AutoscalingGroup" />
    );
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
