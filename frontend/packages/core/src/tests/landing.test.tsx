import React from "react";
import { MemoryRouter } from "react-router-dom";
import { render } from "@testing-library/react";

import "@testing-library/jest-dom";

import * as appContext from "../Contexts/app-context";
import Landing from "../landing";
import { ThemeProvider } from "../Theme";

jest.spyOn(appContext, "useAppContext").mockReturnValue({ workflows: [] });

test("renders correctly", () => {
  const { asFragment } = render(
    <MemoryRouter>
      <ThemeProvider>
        <Landing />
      </ThemeProvider>
    </MemoryRouter>
  );

  expect(asFragment()).toMatchSnapshot();
});
