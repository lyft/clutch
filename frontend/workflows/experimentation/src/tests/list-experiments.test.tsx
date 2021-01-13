import React from "react";
import { shallow } from "enzyme";

import ListExperiments from "../list-experiments";

jest.mock("react-router-dom", () => {
  return {
    ...jest.requireActual("react-router-dom"),
    useNavigate: jest.fn(),
  };
});

describe("List Experiments workflow", () => {
  let component;

  beforeAll(() => {
    const links = [
      {
        displayName: "button_1",
        path: "/path1",
      },
    ];
    const columns = [
      {
        id: "column_1",
        header: "column 1",
      },
      {
        id: "column_2",
        header: "column 2",
      },
    ];

    component = shallow(<ListExperiments columns={columns} links={links} />);
  });

  it("renders correctly", () => {
    expect(component.debug()).toMatchSnapshot();
  });
});
