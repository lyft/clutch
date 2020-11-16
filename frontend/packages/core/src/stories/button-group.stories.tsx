import React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import type { ButtonGroupProps } from "../button";
import { ButtonGroup } from "../button";

export default {
  title: "Core/Buttons/Button Group",
  component: ButtonGroup,
} as Meta;

const Template = (props: ButtonGroupProps) => <ButtonGroup {...props} />;

const sharedArgs = [
  { text: "Back", variant: "neutral", onClick: action("onClick event") },
  action("onClick event"),
];

export const Primary = Template.bind({});
Primary.args = {
  buttons: [
    sharedArgs[0],
    {
      text: "Next",
      onClick: sharedArgs[1],
    },
  ],
};

export const Destructive = Template.bind({});
Destructive.args = {
  buttons: [
    sharedArgs[0],
    {
      text: "Terminate",
      variant: "destructive",
      onClick: sharedArgs[1],
    },
  ],
};
