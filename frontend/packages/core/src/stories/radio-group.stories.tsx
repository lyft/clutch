import React from "react";
import type { Meta } from "@storybook/react";

import type { RadioGroupProps } from "../Input/radio-group";
import { RadioGroup } from "../Input/radio-group";

export default {
  title: "Core/Input/RadioGroup",
  component: RadioGroup,
  argTypes: {
    onChange: { action: "onChange event" },
  },
} as Meta;

const Template = (props: RadioGroupProps) => <RadioGroup {...props} />;

export const WithLabelsOnly = Template.bind({});
WithLabelsOnly.args = {
  defaultOption: 1,
  name: "colorOptions",
  label: "Favorite color",
  options: [{ label: "red" }, { label: "green" }, { label: "blue" }],
};

export const WithLabelsAndUniqueValues = Template.bind({});
WithLabelsAndUniqueValues.args = {
  name: "colorOptionsWithValues",
  label: "Favorite color",
  options: [
    { label: "red", value: "#FF0000" },
    { label: "green", value: "#00FF00" },
    { label: "blue", value: "#0000FF" },
  ],
};
