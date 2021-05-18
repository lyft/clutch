import React from "react";
import { shallow } from "enzyme";

import StartRedisExperiment from "../start-experiment";

jest.mock("react-router-dom", () => {
  return {
    ...jest.requireActual("react-router-dom"),
    useNavigate: jest.fn(),
  };
});

describe("Start Experiment Run workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(<StartRedisExperiment heading="testing" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
