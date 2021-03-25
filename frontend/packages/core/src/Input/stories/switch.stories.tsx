import * as React from "react";
import type { Meta } from "@storybook/react";

import type { SwitchProps } from "../_switch";
import Switch from "../_switch";

export default {
  title: "Core/Input/Switch",
  component: Switch,
  argTypes: {
    onChange: { action: "onChange event" },
  },
} as Meta;

const Template = (props: SwitchProps) => <Switch {...props} />;

export const Primary = Template.bind({});

export const Checked = Template.bind({});
Checked.args = {
  checked: true,
};

export const Disabled = Template.bind({});
Disabled.args = {
  disabled: true,
};
