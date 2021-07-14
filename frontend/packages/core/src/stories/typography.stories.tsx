import * as React from "react";
import type { Meta, Story } from "@storybook/react";

import type { TypographyProps } from "../typography";
import { Typography, VARIANTS } from "../typography";

export default {
  title: "Core/Typography",
  component: Typography,
  argTypes: {
    variant: {
      options: VARIANTS,
      control: { type: "select" },
    },
  },
} as Meta;

const Template: Story<TypographyProps> = ({ variant, children }) => (
  <Typography variant={variant}>{children}</Typography>
);

export const Primary = Template.bind({});
Primary.args = {
  children: "Some text",
  variant: "h1",
};
