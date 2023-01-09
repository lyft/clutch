import * as React from "react";
import type { Meta } from "@storybook/react";

import { Button as ButtonComponent, ButtonProps } from "../button";

const VARIANTS = ["neutral", "primary", "danger", "destructive", "secondary"];

export default {
  title: "Core/Buttons/Button",
  component: ButtonComponent,
  argTypes: {
    onClick: { action: "onClick event" },
    variant: {
      options: VARIANTS,
      control: {
        type: "select",
      },
    },
  },
} as Meta;

const Template = (props: ButtonProps) => <ButtonComponent {...props} />;

export const Button = Template.bind({});
Button.args = {
  text: "Text",
  variant: "primary",
  disabled: false,
};
