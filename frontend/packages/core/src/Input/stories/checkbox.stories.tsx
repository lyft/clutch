import React from "react";
import type { Meta } from "@storybook/react";

import type { CheckboxPanelProps } from "../checkbox";
import CheckboxPanel from "../checkbox";

export default {
  title: "Core/Input/CheckboxPanel",
  component: CheckboxPanel,
  argTypes: {
    onChange: { action: "onChange event" },
  },
} as Meta;

const Template = (props: CheckboxPanelProps) => <CheckboxPanel {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  header: "Select all that apply:",
  options: {
    "Option 1": false,
    "Option 2": false,
    "Option 3": false,
  },
};

export const Disabled = Template.bind({});
Disabled.args = {
  ...Primary.args,
  disabled: true,
};
