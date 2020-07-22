import React from "react";
import { shallow } from "enzyme";

import AppLayout from "..";

describe("AppLayout component", () => {
  let component;

  beforeAll(() => {
    component = shallow(<AppLayout />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
