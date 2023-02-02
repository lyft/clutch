import React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import {
  Button,
  ButtonGroup as ButtonGroupComponent,
  ButtonGroupProps,
  ButtonProps,
} from "../button";

export default {
  title: "Core/Buttons/Button Group",
  component: ButtonGroupComponent,
} as Meta;

const Template = ({ children, ...props }: ButtonGroupProps) => (
  <ButtonGroupComponent {...props}>
    <Button text="Back" variant="neutral" onClick={action("onClick event")} />
    {children as React.ReactElement<ButtonProps>}
  </ButtonGroupComponent>
);

export const Primary = Template.bind({});
Primary.args = {
  children: <Button text="Next" onClick={action("onClick event")} />,
  border: "top",
};

export const Destructive = Template.bind({});
Destructive.args = {
  children: <Button text="Terminate" variant="destructive" onClick={action("onClick event")} />,
  border: "top",
};
