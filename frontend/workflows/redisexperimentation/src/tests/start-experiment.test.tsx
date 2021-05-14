import React from "react";
import { configure, shallow } from "enzyme";
import Adapter from "enzyme-adapter-react-16";

import StartRedisExperiment from "../start-experiment";

jest.mock("react-router-dom", () => {
  return {
    ...jest.requireActual("react-router-dom"),
    useNavigate: jest.fn(),
  };
});

configure({ adapter: new Adapter() });
describe("Start Experiment Run workflow", () => {
  let component;

  beforeAll(() => {
    component = shallow(<StartRedisExperiment heading="testing" />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
