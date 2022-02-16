import React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import type { WarningProps } from "../toast";
import { Warning } from "../toast";

export default {
  title: "Core/Feedback/Warning",
  component: Warning,
} as Meta;

const Template = (props: WarningProps) => <Warning duration={null} {...props} />;

export const Primary = Template.bind({});
Primary.args = {
  message: "Informational warning",
};

export const CloseAction = Template.bind({});
CloseAction.args = {
  ...Primary.args,
  onClose: action("on-close"),
};
