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

export const Basic = Template.bind({});
Basic.args = {
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
  ...Basic.args,
  disabled: true,
};

export const Error = Template.bind({});
Error.args = {
  ...Basic.args,
  error: true,
  helperText: "There was a problem!",
};

export const CustomValues = Template.bind({});
CustomValues.args = {
  ...Basic.args,
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
  ...Basic.args,
  label: null,
};
