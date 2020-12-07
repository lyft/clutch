import React from "react";
import type { Meta } from "@storybook/react";

import { Button } from "../button";
import type { DialogProps } from "../dialog";
import Dialog from "../dialog";

export default {
  title: "Core/Dialog",
  component: Dialog,
  argTypes: {
    onClose: { action: "onClose event" },
  },
} as Meta;

const Template = (props: DialogProps) => (
  <Dialog {...props}>
    <Button text="Yes" />
    <Button text="No" />
  </Dialog>
);

export const Primary = Template.bind({});
Primary.args = {
  title: "Dialog's Title",
  content:
    "This is the content of the dialog. This is the content of the dialog. This is the content of the dialog.",
  open: true,
};
