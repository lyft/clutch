import React from "react";
import createRouterContext from "react-router-test-context";
import { shallow } from "enzyme";

import * as appContext from "../Contexts/app-context";
import Landing from "../landing";

describe("Landing component", () => {
  let component;
  const context = createRouterContext();

  beforeAll(() => {
    appContext.useAppContext = jest.fn().mockReturnValue({ workflows: [] });
    component = shallow(<Landing workflows={[]} />, { context });
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
