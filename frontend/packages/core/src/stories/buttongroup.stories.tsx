import React from "react";
import type { Meta } from "@storybook/react";

import type { ButtonGroupProps } from "../button";
import { ButtonGroup, useWizardContext } from "@clutch-sh/core";

export default {
  title: "Core/Button Group",
  component: ButtonGroup,
  argTypes: {
    onClick: { action: "onClick event" },
  },
} as Meta;

const Template = (props: ButtonGroupProps) => <ButtonGroup {...props} />;

export const Default = Template.bind({});
Default.args = {
  buttons: [
    {
      text: "Back",
    },
    {
      text: "Next",
    },
  ]
};

export const Destructive = Template.bind({});
Destructive.args = {
  buttons: [
    {
      text: "Back",
    },
    {
      text: "Delete",
      destructive: true,
    },
  ]
};