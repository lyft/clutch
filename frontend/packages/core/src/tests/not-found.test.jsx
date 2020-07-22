import React from "react";
import { mount } from "enzyme";
import { ThemeProvider as StyledThemeProvider } from "styled-components";

import { getTheme } from "../AppProvider/themes";
import NotFound from "../not-found";

describe("Not Found component", () => {
  let component;

  beforeAll(() => {
    component = mount(
      <StyledThemeProvider theme={getTheme()}>
        <NotFound />
      </StyledThemeProvider>
    );
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
