import React from "react";
import { render, screen } from "@testing-library/react";

import "@testing-library/jest-dom";

import { Tab, Tabs } from "../tab";

beforeEach(() => {
  render(
    <Tabs value={1}>
      <Tab label="meow" />
      <Tab label="mix" />
    </Tabs>
  );
});

test("renders correctly", () => {
  expect(screen.getByTestId("styled-tabs")).toBeVisible();
});

test("renders children tabs", () => {
  expect(screen.getAllByRole("tab")).toHaveLength(2);
});

test("has the mix tab selected", () => {
  expect(screen.getAllByRole("tab")[1]).toHaveTextContent("mix");
  expect(screen.getAllByRole("tab")[1]).toHaveAttribute("tabindex", "0");
});
