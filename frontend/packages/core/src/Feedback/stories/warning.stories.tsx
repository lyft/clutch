import React from "react";
// TODO: remove when https://github.com/lyft/clutch/pull/607/files#r512255592 lands
// eslint-disable-next-line import/no-extraneous-dependencies
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import type { WarningProps } from "../warning";
import Warning from "../warning";

export default {
  title: "Core/Warning",
  component: Warning,
} as Meta;

const Template = (props: WarningProps) => <Warning {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  message: "Informational warning",
};

export const CloseAction = Template.bind({});
CloseAction.args = {
  ...Primary.args,
  onClose: action("on-close"),
};
