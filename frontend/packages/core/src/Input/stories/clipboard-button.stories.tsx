import React from "react";
import type { Meta } from "@storybook/react";

import type { ClipboardButtonProps } from "../button";
import { ClipboardButton as ClipboardButtonComponent } from "..";

export default {
  title: "Core/Buttons/Clipboard Button",
  component: ClipboardButtonComponent,
} as Meta;

const Template = (props: ClipboardButtonProps) => <ClipboardButtonComponent {...props} />;

export const ClipboardButton = Template.bind({});
ClipboardButton.args = {
  text: "https://clutch.sh",
};
