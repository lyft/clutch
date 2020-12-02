import React from "react";
import styled from "@emotion/styled";
import {
  Box,
  ClickAwayListener,
  Grow,
  IconButton,
  MenuList,
  Paper,
  Popper,
} from "@material-ui/core";
import NotificationsIcon from "@material-ui/icons/Notifications";

const StyledNotificationsIcon = styled(IconButton)({
  color: "#ffffff",
  marginRight: "1rem",
  borderRadius: "8px",
  padding: "8px",
  "&:hover": {
    background: "#2d3db4",
  },
  "&:active": {
    background: "#2938a5",
  },
});

const Notifications: React.FC = () => {
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef(null);

  const handleToggle = () => {
    setOpen(!open);
  };

  const handleClose = event => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }
    setOpen(false);
  };

  function handleListKeyDown(event) {
    if (event.key === "Tab") {
      event.preventDefault();
      setOpen(false);
    }
  }

  return (
    <Box>
      <StyledNotificationsIcon
        ref={anchorRef}
        edge="end"
        aria-controls={open ? "notifcation-options" : undefined}
        aria-haspopup="true"
        onClick={handleToggle}
      >
        <NotificationsIcon />
      </StyledNotificationsIcon>
      <Popper open={open} anchorEl={anchorRef.current} role={undefined} transition>
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{ transformOrigin: placement === "bottom" ? "center top" : "center bottom" }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList
                  autoFocusItem={open}
                  id="notifcation-options"
                  onKeyDown={handleListKeyDown}
                />
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </Box>
  );
};

export default Notifications;
