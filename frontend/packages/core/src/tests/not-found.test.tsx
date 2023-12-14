import React from "react";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import { Theme } from "../AppProvider/themes";
import NotFound from "../not-found";

test("renders correctly", () => {
  const { asFragment } = render(
    <Theme>
      <NotFound />
    </Theme>
  );

  expect(asFragment()).toMatchSnapshot();
});
