import React from "react";
import { mount } from "enzyme";

import { Theme } from "../AppProvider/themes";
import NotFound from "../not-found";

describe("Not Found component", () => {
  let component;

  beforeAll(() => {
    component = mount(
      <Theme>
        <NotFound />
      </Theme>
    );
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
