import React from "react";
import { render, unmountComponentAtNode } from "react-dom";

import "@testing-library/jest-dom";

import { Button } from "../Input";

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
  render(<Button text="test" />, container);

  expect(container.innerHTML).toMatchSnapshot();
});

test("Neutral Button Component", () => {
  render(<Button variant="neutral" text="test" />, container);

  expect(container.innerHTML).toMatchSnapshot();
});

test("Destructive Button Component", () => {
  render(<Button variant="destructive" text="test" />, container);

  expect(container.innerHTML).toMatchSnapshot();
});

test("Small Button Component", () => {
  render(<Button size="small" text="test" />, container);

  expect(container.innerHTML).toMatchSnapshot();
});

test("Large Button Component", () => {
  render(<Button size="large" text="test" />, container);

  expect(container.innerHTML).toMatchSnapshot();
});
