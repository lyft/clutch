import React from "react";
import styled from "@emotion/styled";

import { IconButton } from "../button";

const StyledIcon = styled(IconButton)<{ open?: boolean }>(
  {
    color: "#ffffff",
    marginRight: "8px",
    padding: "12px",
    "&:hover": {
      background: "#2d3db4",
    },
    "&:active": {
      background: "#2938a5",
    },
  },
  props => ({
    background: props.open ? "#2d3db4" : "unset",
  })
);

const HeaderPopper = ({ children }) => {
  return (
    <>
      <StyledIcon>{children}</StyledIcon>
    </>
  );
};

export default HeaderPopper;
