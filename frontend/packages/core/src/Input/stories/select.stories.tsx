import * as React from "react";
import type { Meta } from "@storybook/react";

import type { SelectProps } from "../select";
import Select from "../select";

export default {
  title: "Core/Input/Select",
  component: Select,
  argTypes: {
    onChange: { action: "onChange event" },
    options: { control: { type: "object" } },
  },
} as Meta;

const Template = (props: SelectProps) => <Select name="storybookDemo" {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  label: "My Label",
  options: [
    {
      label: "Option 1",
    },
    {
      label: "Option 2",
    },
  ],
};

export const Disabled = Template.bind({});
Disabled.args = {
  ...Primary.args,
  disabled: true,
};

export const Error = Template.bind({});
Error.args = {
  ...Primary.args,
  error: true,
  helperText: "There was a problem!",
};

export const CustomValues = Template.bind({});
CustomValues.args = {
  ...Primary.args,
  options: [
    {
      label: "Option 1",
      value: "VALUE_ONE",
    },
    {
      label: "Option 2",
      value: "VALUE_TWO",
    },
  ],
};

export const WithoutLabel = Template.bind({});
WithoutLabel.args = {
  ...Primary.args,
  label: null,
};

export const WithStartAdornment = Template.bind({});
WithStartAdornment.args = {
  ...Primary.args,
  options: [
    {
      label: "Option 1",
      startAdornment: <img src="https://clutch.sh/img/microsite/logo.svg" alt="logo" />,
    },
  ],
};
