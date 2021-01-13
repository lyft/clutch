import React from "react";
import { shallow } from "enzyme";

import StartExperiment from "../start-experiment";

jest.mock("react-router-dom", () => {
  return {
    ...jest.requireActual("react-router-dom"),
    useNavigate: jest.fn(),
  };
});

describe("Start Experiment workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(<StartExperiment upstreamClusterTypeSelectionEnabled />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
