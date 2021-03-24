import * as React from "react";
import SearchIcon from '@material-ui/icons/Search';
import type { Meta } from "@storybook/react";

import type { IconButtonProps } from "../button";
import { IconButton } from "../button";

export default {
  title: "Core/Buttons/Icon Button",
  component: IconButton,
  argTypes: {
    onClick: { action: "onClick event" },
  },
} as Meta;

const Template = (props: IconButtonProps) => <IconButton {...props}><SearchIcon /></IconButton>;

export const Primary = Template.bind({});

export const Disabled = Template.bind({});
Disabled.args = {
  disabled: true,
};
