import React from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

import "@testing-library/jest-dom";

import * as ApplicationContext from "../../Contexts/app-context";
import contextValues from "../../Contexts/tests/testContext";
import { ThemeProvider } from "../../Theme";
import { NPSHeader } from "..";

beforeEach(() => {
  jest.spyOn(ApplicationContext, "useAppContext").mockReturnValue(contextValues);
  jest.useFakeTimers();
});

test("renders correctly", () => {
  render(
    <ThemeProvider>
      <NPSHeader />
    </ThemeProvider>
  );

  expect(screen.getAllByRole("button")).toHaveLength(1);
  expect(screen.getByRole("button")).toHaveAttribute("id", "headerFeedbackIcon");
});

test("opens a popper on click of feedback icon", async () => {
  const user = userEvent.setup({ delay: null });
  render(
    <ThemeProvider>
      <NPSHeader />
    </ThemeProvider>
  );

  const feedbackButton = await screen.findByRole("button");
  await user.click(feedbackButton);
  expect(screen.getByRole("tooltip")).toBeInTheDocument();
});
