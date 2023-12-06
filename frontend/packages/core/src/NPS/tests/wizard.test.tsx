import React from "react";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

import "@testing-library/jest-dom";

import { client } from "../../Network";
import { ThemeProvider } from "../../Theme";
import { NPSWizard } from "..";

const defaultResult = {
  prompt: "Test Prompt",
  freeformPrompt: "Test Freeform",
  ratingLabels: [
    {
      emoji: "SAD",
      label: "bad",
    },
    {
      emoji: "NEUTRAL",
      label: "ok",
    },
    {
      emoji: "HAPPY",
      label: "great",
    },
  ],
};

beforeEach(() => {
  jest.spyOn(client, "post").mockReturnValue(
    new Promise((resolve, reject) => {
      resolve({
        data: {
          originSurvey: {
            WIZARD: defaultResult,
            HEADER: defaultResult,
          },
        },
      });
    })
  );
});

afterEach(() => {
  jest.restoreAllMocks();
});

test("renders correctly", async () => {
  render(
    <ThemeProvider>
      <NPSWizard />
    </ThemeProvider>
  );

  expect(await screen.findByTestId("nps-wizard")).toBeVisible();
});

test("renders the container with a bluish background", async () => {
  render(
    <ThemeProvider>
      <NPSWizard />
    </ThemeProvider>
  );

  expect(await screen.findByTestId("nps-wizard")).toHaveStyle({
    background: "#F9F9FE",
  });
});

test("removes the bluish background", async () => {
  const user = userEvent.setup();
  render(
    <ThemeProvider>
      <NPSWizard />
    </ThemeProvider>
  );

  const emojiButton = await screen.findByLabelText(/Great/i);
  await user.click(emojiButton);

  const submitButton = await screen.findByText("Submit");
  await user.click(submitButton);

  expect(screen.getByTestId("nps-wizard")).toHaveStyle({
    background: "unset",
  });
});
