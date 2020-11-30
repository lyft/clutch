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
  let links;
  let columns;

  beforeAll(() => {
    links = [
      {
        displayName: "button_1",
        path: "/path1",
      },
    ];
    columns = [
      {
        id: "column_1",
        header: "column 1",
      },
      {
        id: "column_2",
        header: "column 2",
      },
    ];
  });

  it("renders correctly", () => {
    const component = shallow(<ListExperiments heading="List Experiments" columns={columns} links={links} />);
    expect(component.debug()).toMatchSnapshot();
  });
});
