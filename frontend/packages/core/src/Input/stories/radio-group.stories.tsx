import React from "react";
import type { Meta } from "@storybook/react";

import type { RadioGroupProps } from "../radio-group";
import { RadioGroup } from "../radio-group";

export default {
  title: "Core/Input/RadioGroup",
  component: RadioGroup,
  argTypes: {
    onChange: { action: "onChange event" },
  },
} as Meta;

const Template = (props: RadioGroupProps) => <RadioGroup {...props} />;

export const Primary = Template.bind({});
const options = [{ label: "red" }, { label: "green" }, { label: "blue" }];
Primary.argTypes = {
  defaultOption: {
    control: {
      type: "select",
      options: options.map((_: any, i: number) => i),
    },
  },
};
Primary.args = {
  defaultOption: 1,
  name: "colorOptions",
  label: "Favorite color",
  options,
};

export const UniqueValues = Template.bind({});
UniqueValues.argTypes = Primary.argTypes;
UniqueValues.args = {
  ...Primary.args,
  options: [
    { label: "red", value: "#FF0000" },
    { label: "green", value: "#00FF00" },
    { label: "blue", value: "#0000FF" },
  ],
};
