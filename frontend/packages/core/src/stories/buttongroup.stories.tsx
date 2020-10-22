import React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import type { ButtonGroupProps } from "../button";
import { ButtonGroup } from "../button";

export default {
  title: "Core/Button Group",
  component: ButtonGroup,
} as Meta;

const Template = (props: ButtonGroupProps) => <ButtonGroup {...props} />;

export const Default = Template.bind({});
Default.args = {
  buttons: [
    {
      text: "Back",
      onClick: action("onClick event"),
    },
    {
      text: "Next",
      onClick: action("onClick event"),
    },
  ],
};

export const Destructive = Template.bind({});
Destructive.args = {
  buttons: [
    {
      text: "Back",
      onClick: action("onClick event"),
    },
    {
      text: "Delete",
      destructive: true,
      onClick: action("onClick event"),
    },
  ],
};
