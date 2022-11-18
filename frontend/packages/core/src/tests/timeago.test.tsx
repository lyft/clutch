import React from "react";
import { render, screen } from "@testing-library/react";

import TimeAgo from "../timeago";

const getTimestamp = (days, date = new Date()) => {
  date.setDate(date.getDate() - days);

  return date.getTime();
};

test("loads and displays a timeago in short format given a timestamp", async () => {
  render(<TimeAgo live={false} date={getTimestamp(2)} />);

  expect(screen.getByText("2d")).toBeInTheDocument();
});

test("loads and displays a timeago in long format given a timestamp", async () => {
  render(<TimeAgo live={false} short={false} date={getTimestamp(6)} />);

  expect(screen.getByText("6 days")).toBeInTheDocument();
});

test("loads and displays a short month format given a timestamp", async () => {
  render(<TimeAgo live={false} date={getTimestamp(35)} />);

  expect(screen.getByText("1mo")).toBeInTheDocument();
});
