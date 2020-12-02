import React from "react";
import type { Meta } from "@storybook/react";

import { Button } from "../button";
import type { DialogProps } from "../dialog";
import { Dialog } from "../dialog";
import { TextField } from "../Input/text-field";

export default {
  title: "Core/Dialog",
  component: Dialog,
  argTypes: {
    onClose: { action: "onClose event" },
  },
} as Meta;

const Template = ({ children, ...props }: DialogProps) => (
  <Dialog
    actions={
      <>
        <Button text="Back" variant="neutral" />
        <Button text="Continue" />
      </>
    }
    {...props}
  >
    {children}
  </Dialog>
);

export const Primary = Template.bind({});
Primary.args = {
  title: "Title",
  children:
    "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
  open: true,
};

export const WithTextField = Template.bind({});
WithTextField.args = {
  ...Primary.args,
  children: (
    <>
      Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
      labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
      laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
      voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat
      non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
      <TextField label="Title" placeholder="Input Text" />
    </>
  ),
};
