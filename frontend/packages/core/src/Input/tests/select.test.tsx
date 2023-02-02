import * as React from "react";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import { Select } from "../select";

test("has lower bound", () => {
  const { container } = render(
    <Select name="foobar" defaultOption={-1} options={[{ label: "foo" }, { label: "bar" }]} />
  );

  expect(container.querySelector("#foobar-select")).toBeInTheDocument();
  expect(container.querySelector("#foobar-select")).toHaveTextContent("foo");
});

test("has upper bound", () => {
  const { container } = render(
    <Select name="foobar" defaultOption={2} options={[{ label: "foo" }]} />
  );

  expect(container.querySelector("#foobar-select")).toBeInTheDocument();
  expect(container.querySelector("#foobar-select")).toHaveTextContent("foo");
});
