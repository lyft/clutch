import * as React from "react";
import type { Meta } from "@storybook/react";

import type { CheckboxProps } from "../checkbox";
import { Checkbox } from "../checkbox";

export default {
  title: "Core/Input/Checkbox",
  component: Checkbox,
  argTypes: {
    onChange: { action: "onChange event" },
  },
} as Meta;

const Template = (props: CheckboxProps) => <Checkbox {...props} />;

export const Unselected = Template.bind({});
Unselected.args = {
  disabled: false,
  checked: false,
};

export const Selected = Template.bind({});
Selected.args = {
  disabled: false,
  checked: true,
};

export const Disabled = Template.bind({});
Disabled.args = {
  disabled: true,
  checked: false,
};
