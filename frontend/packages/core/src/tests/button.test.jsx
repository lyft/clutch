import React from "react";
import { shallow } from "enzyme";

import { AdvanceButton, Button, DestructiveButton } from "../button";

describe("Advance Button component", () => {
  let component;

  beforeAll(() => {
    component = shallow(<AdvanceButton text="test" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});

describe("Button component", () => {
  let component;

  beforeAll(() => {
    component = shallow(<Button text="test" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});

describe("Destructive Button component", () => {
  let component;

  beforeAll(() => {
    component = shallow(<DestructiveButton text="test" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
