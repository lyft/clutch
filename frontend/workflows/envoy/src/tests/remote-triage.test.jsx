import React from "react";
import { shallow } from "enzyme";

import RemoteTriage from "../remote-triage";

describe("Remote Triage workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(<RemoteTriage />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
