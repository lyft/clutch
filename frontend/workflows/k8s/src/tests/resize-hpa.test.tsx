import React from "react";
import { shallow } from "enzyme";

import ResizeHPA from "../resize-hpa";

describe("Resize Autoscaling Group workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(<ResizeHPA resolverType="clutch.aws.k8s.v1.HPA" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
