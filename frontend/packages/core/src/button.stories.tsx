import React from "react";
import type { Meta } from "@storybook/react";

import type { ButtonProps } from "./button";
import { Button } from "./button";

export default {
  title: "Core/Button",
  component: Button,
  argTypes: {
    onClick: { action: "onClick event" },
  },
} as Meta;

const Template = (props: ButtonProps) => <Button {...props} />;

export const Default = Template.bind({});
Default.args = {
  text: "click here",
};

export const Destructive = Template.bind({});
Destructive.args = {
  text: "delete",
  destructive: true,
};
