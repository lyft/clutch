import * as React from "react";
import SearchIcon from "@mui/icons-material/Search";
import type { Meta } from "@storybook/react";

import type { IconButtonProps } from "../button";
import { ICON_BUTTON_VARIANTS, IconButton as IconButtonComponent } from "../button";

export default {
  title: "Core/Buttons/Icon Button",
  component: IconButtonComponent,
  argTypes: {
    onClick: { action: "onClick event" },
    size: {
      options: ICON_BUTTON_VARIANTS,
      control: { type: "select" },
    },
  },
} as Meta;

const Template = (props: IconButtonProps) => (
  <IconButtonComponent {...props}>
    <SearchIcon />
  </IconButtonComponent>
);

export const IconButton = Template.bind({});
IconButton.args = {
  disabled: false,
};
