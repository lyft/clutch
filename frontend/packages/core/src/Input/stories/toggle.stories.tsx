import * as React from "react";
import type { Meta } from "@storybook/react";

import type { SwitchProps } from "../toggle";
import { Switch } from "../toggle";

export default {
  title: "Core/Input/Switch",
  component: Switch,
  name: "Switch",
  argTypes: {
    onChange: { action: "onChange event" },
  },
} as Meta;

const Template = (props: SwitchProps) => <Switch {...props} />;

export const Primary = Template.bind({});

export const Checked = Template.bind({});
Checked.props = {
  checked: true,
};

export const Disabled = Template.bind({});
Disabled.props = {
  disabled: true,
};
