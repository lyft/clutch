import React from "react";
import type { Meta } from "@storybook/react";

import type { StatusProps } from "../icon";
import { StatusIcon } from "../icon";

export default {
  title: "Core/Icon/Status",
  component: StatusIcon,
} as Meta;

const Template = (props: StatusProps) => <StatusIcon {...props} />;

export const Primary = Template.bind({});

export const Success = Template.bind({});
Success.args = {
  variant: "success",
};

export const Failure = Template.bind({});
Failure.args = {
  variant: "failure",
};
