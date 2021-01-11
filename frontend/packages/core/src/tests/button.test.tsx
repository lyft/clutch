import React from "react";
import { shallow } from "enzyme";

import { Button } from "../button";

describe("Primary Button component", () => {
  let component;

  beforeAll(() => {
    component = shallow(<Button text="test" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});

describe("Neutral Button component", () => {
  let component;

  beforeAll(() => {
    component = shallow(<Button variant="neutral" text="test" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});

describe("Destructive Button component", () => {
  let component;

  beforeAll(() => {
    component = shallow(<Button variant="destructive" text="test" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
