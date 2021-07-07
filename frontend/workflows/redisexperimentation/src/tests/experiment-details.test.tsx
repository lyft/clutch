import React from "react";
import { FormFields } from "@clutch-sh/experimentation";
import { shallow } from "enzyme";

import { ExperimentDetails } from "../start-experiment";

jest.mock("@clutch-sh/core", () => {
  return {
    ...(jest.requireActual("@clutch-sh/core") as any),
    useNavigate: jest.fn(),
  };
});

describe("Start Experiment workflow", () => {
  it("renders correctly", () => {
    const component = shallow(<ExperimentDetails environments={[]} onStart={() => {}} />);
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });

  it("renders correctly with upstream cluster type selection enabled", () => {
    const component = shallow(<ExperimentDetails environments={[]} onStart={() => {}} />);
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });

  it("renders correctly with environments selection enabled", () => {
    const component = shallow(
      <ExperimentDetails
        downstreamClusterTemplate=""
        environments={[{ value: "staging" }]}
        onStart={() => {}}
      />
    );
    expect(component.find(FormFields).dive().debug()).toMatchSnapshot();
  });
});
