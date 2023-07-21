import * as React from "react";
import type { Meta } from "@storybook/react";

import type { TimePickerProps } from "../time-picker";
import TimePicker from "../time-picker";

export default {
  title: "Core/Input/TimePicker",
  component: TimePicker,
  argTypes: {
    label: {
      control: "text",
    },
    value: {
      control: "date",
    },
  },
} as Meta;

const Template = (props: TimePickerProps) => <TimePicker {...props} />;

export const PrimaryDemo = ({ ...props }) => {
  const [timeValue, setTimeValue] = React.useState<Date | null>(props.value);

  return (
    <TimePicker
      label={props.label}
      onChange={(newValue: unknown) => {
        setTimeValue(newValue as Date);
      }}
      value={timeValue ?? props.value}
    />
  );
};

PrimaryDemo.args = {
  label: "My Label",
  value: new Date(),
} as TimePickerProps;

export const Disabled = Template.bind({});
Disabled.args = {
  ...PrimaryDemo.args,
  disabled: true,
} as TimePickerProps;
