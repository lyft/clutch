import React from 'react';
import type { Meta } from "@storybook/react";

import type { ButtonProps } from './button'

import { Button } from './button';

export default {
  title: 'Core/Button/Button',
  component: Button,

} as Meta;

const Template = (props: ButtonProps) => <Button {...props} />;

export const Default = Template.bind({});
Default.args = {
  text: "click here",
};

export const Destructive  = Template.bind({});
Destructive.args = {
  ...Default.args,
  text: "delete",
  destructive: true,
};