import * as React from "react";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import { ThemeProvider } from "../../Theme";
import { MultiSelect, Select } from "../select";

test("select has lower bound", () => {
  const { container } = render(
    <ThemeProvider>
      <Select name="foobar" defaultOption={-1} options={[{ label: "foo" }, { label: "bar" }]} />
    </ThemeProvider>
  );

  expect(container.querySelector("#foobar-select")).toBeInTheDocument();
  expect(container.querySelector("#foobar-select")).toHaveTextContent("foo");
});

test("select has upper bound", () => {
  const { container } = render(
    <ThemeProvider>
      <Select name="foobar" defaultOption={2} options={[{ label: "foo" }]} />
    </ThemeProvider>
  );

  expect(container.querySelector("#foobar-select")).toBeInTheDocument();
  expect(container.querySelector("#foobar-select")).toHaveTextContent("foo");
});

test("multi select handles multiple", () => {
  const { container } = render(
    <ThemeProvider>
      {" "}
      <MultiSelect
        defaultOptions={[0, 1]}
        name="foobar"
        selectOptions={[{ label: "foo" }, { label: "bar" }]}
      />
    </ThemeProvider>
  );

  expect(container.querySelector("#foobar-multi-select")).toBeInTheDocument();
  expect(container.querySelector("#foobar-multi-select")).toHaveTextContent("foo, bar");
});
