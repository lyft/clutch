import React from "react";
import { Popover, Typography } from "@material-ui/core";
import HelpIcon from "@material-ui/icons/Help";
import styled from "styled-components";

const HelpIconContainer = styled.div`
  ${({ theme }) => `
  display: flex;
  color: ${theme.palette.text.secondary};
  padding: 5px;
  `}
`;

const Hint: React.FC = ({ children }) => {
  const [anchorEl, setAnchorEl] = React.useState<HTMLElement | null>(null);

  const handleOpen = (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
    setAnchorEl(e.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);

  return (
    <>
      <HelpIconContainer
        aria-owns={open ? "help-popover" : undefined}
        aria-haspopup="true"
        onMouseEnter={handleOpen}
        onMouseLeave={handleClose}
      >
        <HelpIcon fontSize="small" />
      </HelpIconContainer>
      <Popover
        id="help-popover"
        style={{ pointerEvents: "none" }}
        open={open}
        anchorEl={anchorEl}
        anchorOrigin={{
          vertical: "bottom",
          horizontal: "left",
        }}
        transformOrigin={{
          vertical: "top",
          horizontal: "left",
        }}
        onClose={handleClose}
        disableRestoreFocus
      >
        <Typography>{children}</Typography>
      </Popover>
    </>
  );
};

export default Hint;
