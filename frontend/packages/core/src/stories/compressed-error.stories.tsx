import React from "react";
import type { Meta } from "@storybook/react";

import type { CompressedErrorProps } from "../error";
import { CompressedError } from "../error";

export default {
  title: "Core/CompressedError",
  component: CompressedError,
} as Meta;

const longMessage =
  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.";
const longMessageNoPeriods = longMessage.replace(/\./g, "");
const longMessageNoStops = longMessageNoPeriods.replace(/[^\w]/g, "");

const options = {
  default: longMessage,
  "Without Periods": longMessageNoPeriods,
  "Without Stops": longMessageNoStops,
};
const Template = (props: CompressedErrorProps) => {
  const selectedMessage = options[props.message] || props.message;
  return <CompressedError {...props} message={selectedMessage} />;
};

export const Primary = Template.bind({});
Primary.args = {
  message: "The server returned a 500.",
};

export const CustomTitle = Template.bind({});
CustomTitle.args = {
  ...Primary.args,
  title: "An error has occurred",
};

export const LongMessage = Template.bind({});
LongMessage.argTypes = {
  message: {
    control: {
      type: "select",
      options: Object.keys(options),
    },
  },
};
LongMessage.args = {
  message: Object.keys(options)[0],
};
