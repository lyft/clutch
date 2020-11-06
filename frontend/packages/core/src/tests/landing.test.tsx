import React from "react";
import { MemoryRouter } from "react-router-dom";
import { shallow } from "enzyme";

import * as appContext from "../Contexts/app-context";
import Landing from "../landing";

describe("Landing component", () => {
  let component;

  beforeAll(() => {
    appContext.useAppContext = jest.fn().mockReturnValue({ workflows: [] });
    component = shallow(<MemoryRouter><Landing /></MemoryRouter>);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
