import React from "react";
import { shallow } from "enzyme";

import { StartExperiment } from "../start-experiment";

jest.mock("@clutch-sh/core", () => {
  return {
    ...(jest.requireActual("@clutch-sh/core") as any),
    useNavigate: jest.fn(),
  };
});

describe("Start Experiment Run workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(<StartExperiment heading="testing" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
