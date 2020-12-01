import React from "react";
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

const Template = (props: SelectProps) => <Select name="demo" {...props} />;

export const Basic = Template.bind({});
Basic.args = {
  options: [
    {
      label: "option 1",
    },
    {
      label: "option 2",
    },
  ],
};

export const CustomValues = Template.bind({});
CustomValues.args = {
  options: [
    {
      label: "option 1",
      value: "value 1",
    },
    {
      label: "option 2",
      value: "value 2",
    },
  ],
};

export const WithLabel = Template.bind({});
WithLabel.argTypes = {
  defaultOption: {
    control: {
      type: "select",
      options: Basic.args.options.map((_: any, i: number) => i),
    },
  },
};
WithLabel.args = {
  ...Basic.args,
  defaultOption: 0,
  label: "Please select one",
  maxWidth: "100px",
};
