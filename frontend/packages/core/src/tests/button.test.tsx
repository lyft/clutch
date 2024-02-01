import React from "react";
import { render, unmountComponentAtNode } from "react-dom";

import "@testing-library/jest-dom";

import { Button } from "../button";
import { ThemeProvider } from "../Theme";

let container: HTMLElement;
beforeEach(() => {
  // setup a DOM element as a render target
  container = document.createElement("div");
  document.body.appendChild(container);
});

afterEach(() => {
  // cleanup on exiting
  unmountComponentAtNode(container);
  container.remove();
});

test("Primary Button Component", () => {
  render(
    <ThemeProvider>
      <Button text="test" />
    </ThemeProvider>,
    container
  );

  expect(container.innerHTML).toMatchSnapshot();
});

test("Neutral Button Component", () => {
  render(
    <ThemeProvider>
      <Button variant="neutral" text="test" />
    </ThemeProvider>,
    container
  );

  expect(container.innerHTML).toMatchSnapshot();
});

test("Destructive Button Component", () => {
  render(
    <ThemeProvider>
      <Button variant="destructive" text="test" />
    </ThemeProvider>,
    container
  );

  expect(container.innerHTML).toMatchSnapshot();
});

test("Small Button Component", () => {
  render(
    <ThemeProvider>
      <Button size="small" text="test" />
    </ThemeProvider>,
    container
  );

  expect(container.innerHTML).toMatchSnapshot();
});

test("Large Button Component", () => {
  render(
    <ThemeProvider>
      <Button size="large" text="test" />
    </ThemeProvider>,
    container
  );

  expect(container.innerHTML).toMatchSnapshot();
});
