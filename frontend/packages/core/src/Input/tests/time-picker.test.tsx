import * as React from "react";
import { fireEvent, render, screen } from "@testing-library/react";

import "@testing-library/jest-dom";

import { ThemeProvider } from "../../Theme";
import TimePicker from "../time-picker";

afterEach(() => {
  jest.resetAllMocks();
});

const onChange = jest.fn();
test("has padding", () => {
  const { container } = render(
    <ThemeProvider>
      <TimePicker value={new Date()} onChange={onChange} />
    </ThemeProvider>
  );

  expect(container.querySelectorAll(".MuiInputBase-adornedEnd")).toHaveLength(1);
  expect(container.querySelector(".MuiInputBase-adornedEnd")).toHaveStyle({
    "padding-right": "14px",
  });
});

test("onChange is called when valid value", () => {
  render(
    <ThemeProvider>
      <TimePicker value={new Date()} onChange={onChange} />
    </ThemeProvider>
  );

  expect(screen.getByPlaceholderText("hh:mm (a|p)m")).toBeVisible();
  fireEvent.change(screen.getByPlaceholderText("hh:mm (a|p)m"), {
    target: { value: "02:55 AM" },
  });
  expect(onChange).toHaveBeenCalled();
});

test("onChange is not called when invalid value", () => {
  render(
    <ThemeProvider>
      <TimePicker value={new Date()} onChange={onChange} />
    </ThemeProvider>
  );

  expect(screen.getByPlaceholderText("hh:mm (a|p)m")).toBeVisible();
  fireEvent.change(screen.getByPlaceholderText("hh:mm (a|p)m"), {
    target: { value: "invalid" },
  });
  expect(onChange).not.toHaveBeenCalled();
});

test("sets passed value correctly", () => {
  const date = new Date();
  const formattedTime = new Intl.DateTimeFormat("en-US", {
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
  const formattedDate = `${formattedTime}`;
  render(
    <ThemeProvider>
      <TimePicker value={date} onChange={onChange} />
    </ThemeProvider>
  );

  expect(screen.getByPlaceholderText("hh:mm (a|p)m")).toHaveValue(formattedDate);
});

test("displays label correctly", () => {
  const label = "testing";
  render(
    <ThemeProvider>
      <TimePicker value={new Date()} onChange={onChange} label={label} />
    </ThemeProvider>
  );

  expect(screen.getByLabelText(label)).toBeVisible();
});

test("is disabled", () => {
  render(
    <ThemeProvider>
      <TimePicker value={new Date()} onChange={onChange} disabled />
    </ThemeProvider>
  );

  expect(screen.getByPlaceholderText("hh:mm (a|p)m")).toBeDisabled();
});
