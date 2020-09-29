import React from "react";
import { Meta } from "@storybook/react/types-6-0";

import TextField from "../text-field";
import type { TextFieldProps } from "../text-field";

export default {
  title: 'Core/Input/TextField',
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
  label: "TextField"
}