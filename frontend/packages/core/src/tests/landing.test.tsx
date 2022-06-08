import React from "react";
import { MemoryRouter } from "react-router-dom";
import { shallow } from "enzyme";

import * as appContext from "../Contexts/app-context";
import Landing from "../landing";

describe("Landing component", () => {
  let component;

  beforeAll(() => {
    jest.spyOn(appContext, "useAppContext").mockReturnValue({ workflows: [] });
    component = shallow(
      <MemoryRouter>
        <Landing />
      </MemoryRouter>
    );
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
