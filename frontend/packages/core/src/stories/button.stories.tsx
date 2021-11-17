import * as React from "react";
import type { Meta } from "@storybook/react";

import type { ButtonProps } from "../button";
import { Button } from "../button";

export default {
  title: "Core/Buttons/Button",
  component: Button,
  argTypes: {
    onClick: { action: "onClick event" },
  },
} as Meta;

const Template = (props: ButtonProps) => <Button {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  text: "Continue",
};

export const Destructive = Template.bind({});
Destructive.args = {
  text: "Delete",
  variant: "destructive",
};

export const Neutral = Template.bind({});
Neutral.args = {
  text: "Back",
  variant: "neutral",
};

export const Secondary = Template.bind({});
Secondary.args = {
  text: "Submit",
  variant: "secondary",
};

export const Disabled = Template.bind({});
Disabled.args = {
  ...Primary.args,
  disabled: true,
};
