import React from "react";
// TODO: remove when https://github.com/lyft/clutch/pull/607/files#r512255592 lands
// eslint-disable-next-line import/no-extraneous-dependencies
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import type { ErrorProps } from "../error";
import { Error } from "../error";

export default {
  title: "Core/Error",
  component: Error,
} as Meta;

const Template = (props: ErrorProps) => <Error {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  message: "An error occurred",
};

export const Retry = Template.bind({});
Retry.args = {
  ...Primary.args,
  onRetry: action("retry-click"),
};
