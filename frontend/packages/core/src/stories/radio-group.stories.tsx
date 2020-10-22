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

export const Primary = Template.bind({});
Primary.args = {
  name: "continue",
  label: "Favorite color",
  options: [{ label: "red" }, { label: "green" }, { label: "blue" }],
};
