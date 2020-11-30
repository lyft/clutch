import React from "react";
import { shallow } from "enzyme";

import { ExperimentDetails } from "../start-experiment";
import FormFields from '../form-fields';

jest.mock("react-router-dom", () => {
  return {
    ...jest.requireActual("react-router-dom"),
    useNavigate: jest.fn(),
  };
});

describe("Start Experiment workflow", () => {
  it("renders correctly", () => {
    const component = shallow(<ExperimentDetails upstreamClusterTypeSelectionEnabled={false} hostsPercentageBasedTargeting={false} onStart={() => {}} />);
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });

  it("renders correctly with upstream cluster type selection enabled", () => {
    const component = shallow(<ExperimentDetails upstreamClusterTypeSelectionEnabled={true} hostsPercentageBasedTargeting={false} onStart={() => {}} />);
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });

  it("renders correctly with host percentage based faults enabled", () => {
    const component = shallow(<ExperimentDetails upstreamClusterTypeSelectionEnabled={false} hostsPercentageBasedTargeting={false} onStart={() => {}} />);
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });

  it("renders correctly with host percentage based faults and upstream cluster type selecion enabled", () => {
    const component = shallow(<ExperimentDetails upstreamClusterTypeSelectionEnabled={true} hostsPercentageBasedTargeting={true} onStart={() => {}} />);
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });
});
