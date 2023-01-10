import React from "react";
import TabContext from "@mui/lab/TabContext";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import { Tab, Tabs } from "../tab";

jest.mock("@mui/lab/TabContext", () => {
  return jest.fn(() => null);
});

test("TabContext is called with specified value", () => {
  render(
    <Tabs value={1}>
      <Tab label="meow" />
      <Tab label="mix" />
    </Tabs>
  );

  expect(TabContext).toHaveBeenCalledWith(
    expect.objectContaining({ value: "1" }),
    expect.anything()
  );
});
