import React from "react";
import { action } from "@storybook/addon-actions";
import type { Meta } from "@storybook/react";

import type { ButtonGroupProps, ButtonProps } from "../button";
import { Button, ButtonGroup } from "../button";

export default {
  title: "Core/Buttons/Button Group",
  component: ButtonGroup,
} as Meta;

const Template = ({ children, ...props }: ButtonGroupProps) => (
  <ButtonGroup {...props}>
    <Button text="Back" variant="neutral" onClick={action("onClick event")} />
    {children as React.ReactElement<ButtonProps>}
  </ButtonGroup>
);

export const Primary = Template.bind({});
Primary.args = {
  children: <Button text="Next" onClick={action("onClick event")} />,
};

export const Destructive = Template.bind({});
Destructive.args = {
  children: <Button text="Terminate" variant="destructive" onClick={action("onClick event")} />,
};
