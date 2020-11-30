import * as React from "react";
import type { Meta } from "@storybook/react";

import type { TextFieldProps } from "../text-field";
import TextField from "../text-field";

export default {
  title: "Core/Input/TextField",
  component: TextField,
} as Meta;

const Template = (props: TextFieldProps) => <TextField {...props} />;

export const Basic = Template.bind({});
Basic.args = {
  label: "My Label",
  placeholder: "This is a placeholder, start typing",
};

export const Disabled = Template.bind({});
Disabled.args = {
  ...Basic.args,
  disabled: true,
};

export const Error = Template.bind({});
Error.args = {
  ...Basic.args,
  error: true,
  helperText: "There was a problem!",
};

export const WithoutLabel = Template.bind({});
WithoutLabel.args = {
  ...Basic.args,
  label: null,
};

export const MultipleLines = Template.bind({});
MultipleLines.args = {
  ...Basic.args,
  multiline: true,
  defaultValue: "This is\nan example\nof multiline content",
};
