import React from "react";
import { shallow } from "enzyme";

import TerminateInstance from "../terminate-instance";

describe("Terminate Instance workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(<TerminateInstance resolverType="clutch.aws.ec2.v1.Instance" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
