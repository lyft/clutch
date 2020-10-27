import React from "react";
import type { Meta } from "@storybook/react";

import type { ClipboardButtonProps } from "../button";
import { ClipboardButton } from "../button";

export default {
  title: "Core/Clipboard Button",
  component: ClipboardButton,
} as Meta;

const Template = (props: ClipboardButtonProps) => <ClipboardButton {...props} />;

export const Default = Template.bind({});
Default.args = {
  text: "https://clutch.sh",
};
