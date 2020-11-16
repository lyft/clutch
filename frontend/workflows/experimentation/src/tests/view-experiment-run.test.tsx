import React from "react";
import { shallow } from "enzyme";

import ViewExperimentRun from "../view-experiment-run";

jest.mock("react-router-dom", () => {
  return {
    ...jest.requireActual("react-router-dom"),
    useNavigate: jest.fn(),
  };
});

describe("View Experiment Run workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(<ViewExperimentRun />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
