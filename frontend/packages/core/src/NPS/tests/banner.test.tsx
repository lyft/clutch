import React from "react";
import { render, screen } from "@testing-library/react";

import "@testing-library/jest-dom";

import { ApplicationContext } from "../../Contexts";
import { HeaderItem } from "../../Contexts/app-context";
import { ThemeProvider } from "../../Theme";
import { Banner, FeedbackBannerProps } from "../banner";

const customRender = ({ ...props }: FeedbackBannerProps) => {
  let triggeredHeaderData = { NPS: {} };

  return render(
    <ApplicationContext.Provider
      // eslint-disable-next-line react/jsx-no-constructed-context-values
      value={{
        workflows: [],
        triggerHeaderItem: (item: HeaderItem, data: unknown) => {
          triggeredHeaderData = {
            ...triggeredHeaderData,
            [item]: {
              ...(data as any),
            },
          };
        },
        triggeredHeaderData,
      }}
    >
      <ThemeProvider>
        <Banner {...props} />
      </ThemeProvider>
    </ApplicationContext.Provider>
  );
};

test("An NPS banner component with default feedback text", () => {
  customRender({});

  const renderedText = screen.getByTestId("nps-banner-text");
  expect(renderedText.textContent).toBe("Enjoying this feature? We would like your feedback!");
  expect(renderedText).toBeVisible();
});

test("An NPS banner component with default button text", () => {
  customRender({});

  const renderedText = screen.getByTestId("nps-banner-button");
  expect(renderedText.textContent).toBe("Give Feedback");
  expect(renderedText).toBeVisible();
});

test("An NPS banner component with custom feedback text", () => {
  const customText = "Testing Feedback Text";
  customRender({ feedbackText: customText });

  const renderedText = screen.getByTestId("nps-banner-text");
  expect(renderedText.textContent).toBe(customText);
  expect(renderedText).toBeVisible();
});

test("An NPS banner component with custom button text", () => {
  const customText = "Testing Feedback Text";
  customRender({ feedbackButtonText: customText });

  const renderedText = screen.getByTestId("nps-banner-button");
  expect(renderedText.textContent).toBe(customText);
  expect(renderedText).toBeVisible();
});

test("An NPS banner component with a default container", () => {
  customRender({});

  expect(screen.getByTestId("nps-banner-container")).toBeVisible();
});

test("An NPS banner component not elevated on the page", () => {
  customRender({ elevated: false });

  expect(screen.queryByTestId("nps-banner-container")).not.toBeInTheDocument();
});
