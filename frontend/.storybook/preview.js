import React from "react";
import { Theme } from "./../packages/core/src/AppProvider/themes";

export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
}

export const decorators = [
  (Story) => (
    <Theme variant="light">
      <Story />
    </Theme>
  ),
];