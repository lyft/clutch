import React from "react";
import { Theme } from "./../packages/core/src/AppProvider/themes";

export const decorators = [
  (Story) => (
    <Theme variant="light">
      <Story />
    </Theme>
  ),
];