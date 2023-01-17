import * as React from "react";
import { fireEvent, render, screen } from "@testing-library/react";

import "@testing-library/jest-dom";

import DateTimePicker from "../date-time";

afterEach(() => {
  jest.resetAllMocks();
});

const onChange = jest.fn();
test("has padding", () => {
  const { container } = render(<DateTimePicker value={new Date()} onChange={onChange} />);

  expect(container.querySelectorAll(".MuiInputBase-adornedEnd")).toHaveLength(1);
  expect(container.querySelector(".MuiInputBase-adornedEnd")).toHaveStyle({
    "padding-right": "14px",
  });
});

test("onChange is called when valid value", () => {
  render(<DateTimePicker value={new Date()} onChange={onChange} />);

  expect(screen.getByPlaceholderText("mm/dd/yyyy hh:mm (a|p)m")).toBeVisible();
  fireEvent.change(screen.getByPlaceholderText("mm/dd/yyyy hh:mm (a|p)m"), {
    target: { value: "11/16/2023 02:55 AM" },
  });
  expect(onChange).toHaveBeenCalled();
});

test("onChange is not called when invalid value", () => {
  render(<DateTimePicker value={new Date()} onChange={onChange} />);

  expect(screen.getByPlaceholderText("mm/dd/yyyy hh:mm (a|p)m")).toBeVisible();
  fireEvent.change(screen.getByPlaceholderText("mm/dd/yyyy hh:mm (a|p)m"), {
    target: { value: "invalid" },
  });
  expect(onChange).not.toHaveBeenCalled();
});

test("sets passed value correctly", () => {
  const date = new Date();
  const formattedDMY = new Intl.DateTimeFormat("en-US", {
    month: "2-digit",
    day: "2-digit",
    year: "numeric",
  }).format(date);
  const formattedTime = new Intl.DateTimeFormat("en-US", {
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
  const formattedDate = `${formattedDMY} ${formattedTime}`;
  render(<DateTimePicker value={date} onChange={onChange} />);

  expect(screen.getByPlaceholderText("mm/dd/yyyy hh:mm (a|p)m")).toHaveValue(formattedDate);
});

test("displays label correctly", () => {
  const label = "testing";
  render(<DateTimePicker value={new Date()} onChange={onChange} label={label} />);

  expect(screen.getByLabelText(label)).toBeVisible();
});

test("is disabled", () => {
  render(<DateTimePicker value={new Date()} onChange={onChange} disabled />);

  expect(screen.getByPlaceholderText("mm/dd/yyyy hh:mm (a|p)m")).toBeDisabled();
});
