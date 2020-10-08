import React from "react";
import type { Meta } from "@storybook/react";

import type { TextFieldProps } from "../text-field";
import TextField from "../text-field";

export default {
  title: "Core/Input/TextField",
  component: TextField,
  argTypes: {
    onReturn: { action: "onReturn event" },
    onChange: { action: "onChange event" },
  },
} as Meta;

const Template = (props: TextFieldProps) => <TextField {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  maxWidth: "",
};

export const WithLabel = Template.bind({});
WithLabel.args = {
  ...Primary.args,
  placeholder: "",
  label: "TextField",
};

export const WithType = Template.bind({});
WithType.args = {
  ...Primary.args,
  label: "Date",
  type: "datetime-local",
};
