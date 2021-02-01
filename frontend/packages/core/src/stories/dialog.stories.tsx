import React, { useState } from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import { Button } from "../button";
import type { DialogProps } from "../dialog";
import { Dialog, DialogActions, DialogContent } from "../dialog";
import { TextField } from "../Input/text-field";

export default {
  title: "Core/Dialog",
  component: Dialog,
} as Meta;

const Template = ({ content, open, ...props }: DialogProps & { content: React.ReactNode }) => {
  const [isOpen, setIsOpen] = useState(open);

  return (
    <>
      <Button disabled={!open} text="Show Dialog" onClick={() => setIsOpen(true)} />
      <Dialog open={isOpen} onClose={() => setIsOpen(false)} {...props}>
        <DialogContent>{content}</DialogContent>
        <DialogActions>
          <Button text="Close" variant="neutral" onClick={() => setIsOpen(false)} />
          <Button text="Continue" onClick={action("continue button click")} />
        </DialogActions>
      </Dialog>
    </>
  );
};

export const Primary = Template.bind({});
Primary.args = {
  title: "Title",
  content:
    "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
  open: true,
};

export const WithTextField = Template.bind({});
WithTextField.args = {
  ...Primary.args,
  content: (
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
