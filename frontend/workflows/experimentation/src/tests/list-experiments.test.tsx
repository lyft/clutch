import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import ListExperiments from "../list-experiments";

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

test("renders correctly", () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ListExperiments heading="List Experiments" columns={columns} links={links} />
    </BrowserRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
