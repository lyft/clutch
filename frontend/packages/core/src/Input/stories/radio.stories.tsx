import React from "react";
import type { Meta } from "@storybook/react";

import type { RadioProps } from "../radio";
import Radio from "../radio";

export default {
  title: "Core/Input/Radio",
  component: Radio,
} as Meta;

const Template = (props: RadioProps) => <Radio {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  selected: false,
};

export const Selected = Template.bind({});
Selected.args = {
  selected: true,
};

export const Disabled = Template.bind({});
Disabled.args = {
  disabled: true,
};

export const NameValueAndOnChange = Template.bind({});
NameValueAndOnChange.args = {
  name: "foo",
  value: "bar",
  onChange: e => e.target.value,
};
