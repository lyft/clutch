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

export const PrimaryDemo = ({ ...props }) => {
  const [dateValue, setDateValue] = React.useState<Date | null>(props.value);

  return (
    <DateTimePicker
      label={props.label}
      onChange={(newValue: unknown) => {
        setDateValue(newValue as Date);
      }}
      value={dateValue ?? props.value}
    />
  );
};

PrimaryDemo.args = {
  label: "My Label",
  value: new Date(),
} as DateTimePickerProps;

export const Disabled = Template.bind({});
Disabled.args = {
  ...PrimaryDemo.args,
  disabled: true,
} as DateTimePickerProps;

export const WithError = Template.bind({});
WithError.args = {
  ...PrimaryDemo.args,
  error: true,
  helperText: "error in the field",
  onChange: (newValue: unknown) => null,
} as DateTimePickerProps;

export const WithMinDate = Template.bind({});
WithMinDate.args = {
  ...PrimaryDemo.args,
  minDate: new Date(),
  onChange: (newValue: unknown) => null,
} as DateTimePickerProps;

export const WithMaxDate = Template.bind({});
WithMaxDate.args = {
  ...PrimaryDemo.args,
  maxDate: new Date(),
  onChange: (newValue: unknown) => null,
} as DateTimePickerProps;
