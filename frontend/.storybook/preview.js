import React from "react";
import { Theme } from "./../packages/core/src/AppProvider/themes";

export const decorators = [
  (Story) => (
    <Theme variant="light">
      <Story />
    </Theme>
  ),
];

export const parameters = {  
  backgrounds: {
    default: "clutch",
    values: [
      {
        name: "clutch",
        value: "#f9fafe",
      },
      {
        name: "light",
        value: "#ffffff",
      },
    ]
  }
};