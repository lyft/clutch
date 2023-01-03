import * as React from "react";
import type { Meta } from "@storybook/react";

import PaperComponent from "../paper";

export default {
  title: "Core/Paper",
  component: PaperComponent,
} as Meta;

export const Paper = () => (
  <PaperComponent>
    <div>Some text in paper</div>
  </PaperComponent>
);
