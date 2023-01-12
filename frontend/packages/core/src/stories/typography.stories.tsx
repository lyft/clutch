import * as React from "react";
import type { Meta, Story } from "@storybook/react";

import type { TypographyProps } from "../typography";
import { Typography as TypographyComponent, VARIANTS } from "../typography";

export default {
  title: "Core/Typography",
  component: TypographyComponent,
  argTypes: {
    variant: {
      options: VARIANTS,
      control: { type: "select" },
    },
    color: {
      control: { type: "color" },
    },
  },
} as Meta;

const Template: Story<TypographyProps> = ({ variant, children, ...props }) => (
  <TypographyComponent variant={variant} {...props}>
    {children}
  </TypographyComponent>
);

export const Typography = Template.bind({});
Typography.args = {
  children: "Some text",
  variant: "h1",
};
