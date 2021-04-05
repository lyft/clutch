import React from "react";
import { IconButton } from "@material-ui/core";
import MoreVertIcon from "@material-ui/icons/MoreVert";
import type { Meta } from "@storybook/react";

import type { PopperProps } from "../popper";
import { Popper, PopperItem } from "../popper";

export default {
  title: "Core/Popper",
  component: Popper,
} as Meta;

const Template = ({children, ...props}: PopperProps) => {
  const anchorRef = React.useRef(null);
  const [open, setOpen] = React.useState(false);
  
  return (
    <>
      <IconButton
        disableRipple
        disabled={React.Children.count(children) <= 0}
        ref={anchorRef}
        onClick={() => setOpen(true)}
      >
        <MoreVertIcon />
      </IconButton>
      <Popper open={open} anchorRef={anchorRef} onClickAway={() => setOpen(false)} {...props}>
        {children}
      </Popper>
    </>
  );
};

export const Primary = Template.bind({});
Primary.args = {
  children: (
    <>
      <PopperItem>
        Item 1
      </PopperItem>
      <PopperItem>
        Item 2
      </PopperItem>
    </>
  )
};

export const NoChildren = Template.bind({});
