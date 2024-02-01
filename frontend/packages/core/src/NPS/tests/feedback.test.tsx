import React from "react";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { capitalize } from "lodash";

import "@testing-library/jest-dom";

import contextValues from "../../Contexts/tests/testContext";
import { client } from "../../Network";
import { ThemeProvider } from "../../Theme";
import NPSFeedback, { defaults, FEEDBACK_MAX_LENGTH } from "../feedback";
import { generateFeedbackTypes } from "../header";

const feedbackTypes = generateFeedbackTypes(contextValues.workflows);

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

describe("api success", () => {
  let spy;
  beforeEach(() => {
    spy = jest.spyOn(client, "post").mockReturnValue(
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

  test("renders survey text prompt", async () => {
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    expect(await screen.findByText(defaultResult.prompt)).toBeVisible();
  });

  test("renders emojis to <EmojiRatings />", async () => {
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    const emojiButtons = await screen.findAllByRole("button");

    expect(defaultResult.ratingLabels.map(({ label }) => capitalize(label))).toEqual(
      emojiButtons.map(element => element.getAttribute("aria-label"))
    );
  });

  test("renders text placeholder", async () => {
    const user = userEvent.setup();
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    await user.click(await screen.findByLabelText(/Great/i));

    expect(await screen.findByPlaceholderText(defaultResult.freeformPrompt)).toBeVisible();
  });

  test("displays a successful submission alert after submit", async () => {
    const user = userEvent.setup();
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    await user.click(await screen.findByLabelText(/Great/i));
    await user.click(await screen.findByText("Submit"));

    expect(await screen.findByText("Thank you for your feedback!")).toBeVisible();
  });

  test("sends feedback upon emoji selection change", async () => {
    const user = userEvent.setup();
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );
    spy.mockClear();
    expect(spy).not.toHaveBeenCalled();

    await user.click(await screen.findByLabelText(/Great/i));
    await user.click(await screen.findByText("Submit"));

    await waitFor(() => {
      expect(spy).toHaveBeenCalled();
      expect(spy).toHaveBeenCalledTimes(1);
    });
  });
});

describe("api failure", () => {
  beforeEach(() => {
    jest.spyOn(client, "post").mockReturnValue(
      new Promise((resolve, reject) => {
        reject(new Error("Test Error"));
      })
    );
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  test("renders default text prompt", async () => {
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    expect(await screen.findByText(defaults.prompt as string)).toBeVisible();
  });

  test("renders default emojis to <EmojiRatings />", async () => {
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    const emojiButtons = await screen.findAllByRole("button");

    if (!defaults.ratingLabels) {
      return;
    }

    expect(defaults.ratingLabels.map(({ label }) => capitalize(label as string))).toEqual(
      emojiButtons.map(element => element.getAttribute("aria-label"))
    );
  });

  test("renders default text placeholder", async () => {
    const user = userEvent.setup();
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    await user.click(await screen.findByLabelText(/Great/i));

    expect(await screen.findByPlaceholderText(defaults.freeformPrompt as string)).toBeVisible();
  });
});

describe("basic rendering", () => {
  const maxLength = FEEDBACK_MAX_LENGTH;
  const generateRandomString = (length, rs = "") => {
    let randomString = rs;
    randomString += Math.random().toString(20).substr(2, length);
    if (randomString.length > length) return randomString.slice(0, length);
    return generateRandomString(length, randomString);
  };

  beforeEach(() => {
    jest.spyOn(client, "post").mockReturnValue(
      new Promise((resolve, reject) => {
        resolve({
          data: {
            originSurvey: {
              WIZARD: defaultResult,
            },
          },
        });
      })
    );
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  test("will not display feedback form or submit unless emoji is selected", () => {
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    expect(screen.getByTestId("feedback-items-container").childElementCount).toBe(2);
  });

  test("will display text prompt at top", async () => {
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    const feedbackItems = await screen.findByTestId("feedback-items-container");
    expect(feedbackItems.childNodes[0].firstChild).toHaveTextContent(defaultResult.prompt);
  });

  test("will display <EmojiRatings /> below prompt", async () => {
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    const feedbackItems = await screen.findByTestId("feedback-items-container");
    feedbackItems.childNodes[1].childNodes.forEach(node => {
      expect(node).toHaveAttribute("aria-label");
    });
  });

  test("displays a text form and submit buttons after selection of emoji", async () => {
    const user = userEvent.setup();
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );
    expect(screen.getByTestId("feedback-items-container").childElementCount).toBe(2);

    await user.click(await screen.findByLabelText(/Great/i));

    expect(await (await screen.findByTestId("feedback-items-container")).childElementCount).toBe(4);
    expect(await screen.findByRole("textbox")).toBeVisible();
    expect(
      await (await screen.findAllByRole("button")).filter(element =>
        element.hasAttribute("aria-label")
      )
    ).toHaveLength(defaultResult.ratingLabels.length);
  });

  test("will update the length on feedback if input is given", async () => {
    const testValue = "Some Feedback Text";
    const user = userEvent.setup();
    const { container } = render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    await user.click(await screen.findByLabelText(/Great/i));

    expect(container.querySelector(".MuiFormHelperText-root")).toHaveTextContent(
      `0 / ${maxLength}`
    );

    const textbox = await screen.findByPlaceholderText(defaultResult.freeformPrompt);
    await user.click(textbox);
    await user.paste(testValue);

    expect(container.querySelector(".MuiFormHelperText-root")).toHaveTextContent(
      `${testValue.trim().length} / ${maxLength}`
    );
    expect(await screen.findByPlaceholderText(defaultResult.freeformPrompt)).toHaveValue(testValue);
  });

  test("will display an error on feedback if more input is given than maxLength", async () => {
    const testValue = generateRandomString(FEEDBACK_MAX_LENGTH + 1);
    const user = userEvent.setup();
    const { container } = render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    user.click(await screen.findByLabelText(/Great/i));

    await waitFor(() => {
      expect(container.querySelector(".MuiFormHelperText-root")).toHaveTextContent(
        `0 / ${maxLength}`
      );
    });

    const textbox = await screen.findByPlaceholderText(defaultResult.freeformPrompt);
    await user.click(textbox);
    await user.paste(testValue);

    expect(await screen.findByRole("textbox")).toHaveValue(testValue);
    expect(container.querySelector(".MuiFormHelperText-root")).toHaveTextContent(
      `${testValue.trim().length} / ${maxLength}`
    );
    expect(container.querySelector(".MuiFormHelperText-root")).toHaveClass("Mui-error");
  });

  test("will disable the submit button upon error", async () => {
    const testValue = generateRandomString(FEEDBACK_MAX_LENGTH + 1);
    const user = userEvent.setup();
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    user.click(await screen.findByLabelText(/Great/i));

    expect(await screen.findByText("Submit")).toBeEnabled();

    const textbox = await screen.findByPlaceholderText(defaultResult.freeformPrompt);
    await user.click(textbox);
    await user.paste(testValue);

    expect(await screen.findByText("Submit")).toBeDisabled();
  });
});

describe("Wizard Origin Rendering", () => {
  beforeEach(() => {
    jest.spyOn(client, "post").mockReturnValue(
      new Promise((resolve, reject) => {
        resolve({
          data: {
            originSurvey: {
              WIZARD: defaultResult,
            },
          },
        });
      })
    );
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  test("renders correctly", () => {
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    expect(screen.getByTestId("feedback-component")).toBeVisible();
  });

  test("styles the submit button correctly", async () => {
    const user = userEvent.setup();
    render(
      <ThemeProvider>
        <NPSFeedback origin="WIZARD" />
      </ThemeProvider>
    );

    user.click(await screen.findByLabelText(/Great/i));

    expect(await screen.findByText("Submit")).toHaveStyle("background-color: transparent");
  });
});

describe("Header Origin Rendering", () => {
  beforeEach(() => {
    jest.spyOn(client, "post").mockReturnValue(
      new Promise((resolve, reject) => {
        resolve({
          data: {
            originSurvey: {
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

  test("renders correctly", () => {
    render(
      <ThemeProvider>
        <NPSFeedback origin="HEADER" feedbackTypes={feedbackTypes} />
      </ThemeProvider>
    );

    expect(screen.getByTestId("feedback-component")).toBeVisible();
  });

  test("styles the submit button correctly", async () => {
    const user = userEvent.setup();
    render(
      <ThemeProvider>
        <NPSFeedback origin="HEADER" feedbackTypes={feedbackTypes} />
      </ThemeProvider>
    );

    user.click(await screen.findByLabelText(/Great/i));

    expect(await screen.findByText("Submit")).toHaveStyle("background-color: #3548D4");
  });
});
