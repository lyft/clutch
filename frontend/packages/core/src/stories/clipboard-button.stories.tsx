import React from "react";
import type { Meta } from "@storybook/react";

import { ClipboardButton } from "../Input";
import type { ClipboardButtonProps } from "../Input/button";

export default {
  title: "Core/Buttons/Clipboard Button",
  component: ClipboardButton,
} as Meta;

const Template = (props: ClipboardButtonProps) => <ClipboardButton {...props} />;

export const Default = Template.bind({});
Default.args = {
  text: "https://clutch.sh",
};
