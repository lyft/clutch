import React from "react";
import type { Meta } from "@storybook/react";

import Paper from "../../paper";
import { NPSWizard } from "..";

export default {
  title: "Core/NPS/Wizard",
  component: NPSWizard,
} as Meta;

const Template = () => (
  <Paper>
    <NPSWizard />
  </Paper>
);

export const Primary = Template.bind({});
