import React from "react";
import { shallow } from "enzyme";

import DeletePod from "../delete-pod";

describe("Terminate Instance workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(<DeletePod resolverType="clutch.k8s.v1.Pod" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
