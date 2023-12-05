import React from "react";
import { BrowserRouter } from "react-router-dom";
import { render, waitFor } from "@testing-library/react";

import "@testing-library/jest-dom";

import * as appContext from "../../Contexts/app-context";
import { client } from "../../Network";
import { ThemeProvider } from "../../Theme";
import AppLayout from "..";

jest.spyOn(appContext, "useAppContext").mockReturnValue({ workflows: [] });
jest.spyOn(client, "post").mockReturnValue(
  new Promise((resolve, reject) => {
    resolve({
      data: {},
    });
  })
);

test("renders correctly", async () => {
  const { asFragment } = render(
    <BrowserRouter>
      <ThemeProvider>
        <AppLayout />
      </ThemeProvider>
    </BrowserRouter>
  );

  await waitFor(() => {
    expect(asFragment()).toMatchSnapshot();
  });
});
