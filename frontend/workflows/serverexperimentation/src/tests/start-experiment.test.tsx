import React from "react";
import { shallow } from "enzyme";

import { StartExperiment } from "../start-experiment";

jest.mock("react-router-dom", () => {
  return {
    ...jest.requireActual("react-router-dom"),
    useNavigate: jest.fn(),
  };
});

describe("Start Experiment workflow", () => {
  it("renders correctly", () => {
    const component = shallow(<StartExperiment heading="Start Experiment" />);
    expect(component.debug()).toMatchSnapshot();
  });

  it("renders correctly with upstream cluster type selection enabled", () => {
    const component = shallow(
      <StartExperiment heading="Start Experiment" upstreamClusterTypeSelectionEnabled />
    );
    expect(component.debug()).toMatchSnapshot();
  });
});
