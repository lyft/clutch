import * as React from "react";
import type { Meta } from "@storybook/react";

import type { DateTimePickerProps } from "../date-time";
import DateTimePicker from "../date-time";

export default {
  title: "Core/Input/DateTimePicker",
  component: DateTimePicker,
  argTypes: {
    value: {
      control: "date",
    },
  },
} as Meta;

const Template = (props: DateTimePickerProps) => <DateTimePicker {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  label: "My Label",
  onChange: () => {},
  value: new Date(),
} as DateTimePickerProps;

export const Disabled = Template.bind({});
Disabled.args = {
  ...Primary.args,
  disabled: true,
} as DateTimePickerProps;
