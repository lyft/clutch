import React from "react";
import type { Meta } from "@storybook/react";

import type { StatusProps } from "../icon";
import { Status } from "../icon";

export default {
  title: "Core/Icon/Status",
  component: Status,
} as Meta;

const Template = (props: StatusProps) => <Status {...props} />;

export const Primary = Template.bind({});

export const Success = Template.bind({});
Success.args = {
  variant: "success",
};

export const Failure = Template.bind({});
Failure.args = {
  variant: "failure",
};
