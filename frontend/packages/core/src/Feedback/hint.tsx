import React from "react";
import HelpIcon from "@mui/icons-material/Help";
import { Popover, Theme, Typography } from "@mui/material";

import styled from "../styled";

const HelpIconContainer = styled("div")(({ theme }: { theme: Theme }) => ({
  display: "flex",
  color: theme.palette.secondary[300],
  padding: "5px",
}));

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
