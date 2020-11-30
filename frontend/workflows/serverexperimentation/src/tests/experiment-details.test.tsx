import React from "react";
import { shallow } from "enzyme";

import FormFields from "../form-fields";
import { ExperimentDetails } from "../start-experiment";

jest.mock("react-router-dom", () => {
  return {
    ...jest.requireActual("react-router-dom"),
    useNavigate: jest.fn(),
  };
});

describe("Start Experiment workflow", () => {
  it("renders correctly", () => {
    const component = shallow(
      <ExperimentDetails
        upstreamClusterTypeSelectionEnabled={false}
        hostsPercentageBasedTargeting={false}
        onStart={() => {}}
      />
    );
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });

  it("renders correctly with upstream cluster type selection enabled", () => {
    const component = shallow(
      <ExperimentDetails
        upstreamClusterTypeSelectionEnabled
        hostsPercentageBasedTargeting={false}
        onStart={() => {}}
      />
    );
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });

  it("renders correctly with host percentage based faults enabled", () => {
    const component = shallow(
      <ExperimentDetails
        upstreamClusterTypeSelectionEnabled={false}
        hostsPercentageBasedTargeting={false}
        onStart={() => {}}
      />
    );
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });

  it("renders correctly with host percentage based faults and upstream cluster type selecion enabled", () => {
    const component = shallow(
      <ExperimentDetails
        upstreamClusterTypeSelectionEnabled
        hostsPercentageBasedTargeting
        onStart={() => {}}
      />
    );
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });
});
