import React from "react";
import { render, screen, waitForElementToBeRemoved } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { capitalize } from "lodash";

import "@testing-library/jest-dom";

import { ThemeProvider } from "../../Theme";
import EmojiRatings from "../emojiRatings";

const stringExample = [
  {
    emoji: 1,
    label: "bad",
  },
  {
    emoji: 2,
    label: "ok",
  },
  {
    emoji: 3,
    label: "great",
  },
];

const emojiMap = {
  1: "SAD",
  2: "NEUTRAL",
  3: "HAPPY",
};

test("will display a given list of emojis and their capitalized labels", () => {
  render(
    <ThemeProvider>
      <EmojiRatings ratings={stringExample} setRating={() => {}} />
    </ThemeProvider>
  );

  const elements = screen.getAllByRole("button");

  expect(elements).toHaveLength(stringExample.length);

  elements.forEach((element, i) => {
    expect(element).toHaveAttribute("aria-label", capitalize(stringExample[i].label));
  });
});

test("all emojis have an initial opacity of 0.5 when not selected", () => {
  render(
    <ThemeProvider>
      <EmojiRatings ratings={stringExample} setRating={() => {}} />
    </ThemeProvider>
  );

  const elements = screen.getAllByRole("button");

  elements.forEach((element, i) => {
    expect(element).toHaveStyle("opacity: 0.5");
  });
});

test("emojis will have a tooltip show on hover", async () => {
  const user = userEvent.setup();
  render(
    <ThemeProvider>
      <EmojiRatings ratings={stringExample} setRating={() => {}} />
    </ThemeProvider>
  );

  await user.hover(screen.getByLabelText(/Great/i));
  expect(await screen.findByText("Great")).toHaveClass("MuiTooltip-tooltip");
  await user.unhover(screen.getByLabelText(/Great/i));
  await waitForElementToBeRemoved(() => screen.queryByText("Great"));
});

test("emojis will update opacity to 1 on selection", async () => {
  const user = userEvent.setup();
  render(
    <ThemeProvider>
      <EmojiRatings ratings={stringExample} setRating={() => {}} />
    </ThemeProvider>
  );

  await user.click(screen.getByLabelText(/Great/i));
  expect(screen.getByLabelText(/Great/i)).toHaveStyle("opacity: 1");
});

test("will return a given emoji on select", async () => {
  const user = userEvent.setup();
  let selected: any = null;

  render(
    <ThemeProvider>
      <EmojiRatings
        ratings={stringExample}
        setRating={rating => {
          selected = rating;
        }}
      />
    </ThemeProvider>
  );

  expect(selected).toBeNull();

  await user.click(screen.getByLabelText(/Ok/i));

  expect(selected.emoji).toEqual(emojiMap[stringExample[1].emoji]);
});
