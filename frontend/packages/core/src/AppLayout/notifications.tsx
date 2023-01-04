import React from "react";
import styled from "@emotion/styled";
import NotificationsIcon from "@mui/icons-material/Notifications";
import {
  ClickAwayListener,
  Grow as MuiGrow,
  IconButton,
  ListItemText as MuiListItemText,
  MenuItem as MuiMenuItem,
  MenuList,
  Paper as MuiPaper,
  Popper as MuiPopper,
} from "@mui/material";

const StyledNotificationsIcon = styled(IconButton)({
  color: "#ffffff",
  marginRight: "8px",
  padding: "12px",
  "&:hover": {
    background: "#2d3db4",
  },
  "&:active": {
    background: "#2938a5",
  },
});

const Popper = styled(MuiPopper)({
  padding: "0 12px",
  marginLeft: "12px",
});

const Paper = styled(MuiPaper)({
  width: "242px",
  border: "1px solid #E7E7EA",
  boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
});

const MenuItem = styled(MuiMenuItem)({
  height: "48px",
  padding: "12px",
  "&:hover": {
    backgroundColor: "#E7E7EA",
  },
  "&:active": {
    backgroundColor: "#EBEDFB",
  },
});

const ListItemText = styled(MuiListItemText)({
  margin: "0px",
  ".MuiTypography-root": {
    color: "#0D1030",
    fontSize: "14px",
    lineHeight: "24px",
  },
});

const Grow = styled(MuiGrow)((props: { placement: string }) => ({
  transformOrigin: props.placement,
}));

interface NotificationsData {
  value: string;
}

export interface NotificationsProp {
  data?: NotificationsData[];
}

const Notifications: React.FC<NotificationsProp> = ({ data }) => {
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

  const handleListKeyDown = event => {
    if (event.key === "Tab") {
      event.preventDefault();
      setOpen(false);
    }
  };

  return (
    <>
      <StyledNotificationsIcon
        ref={anchorRef}
        edge="end"
        aria-controls={open ? "notification-options" : undefined}
        aria-haspopup="true"
        onClick={handleToggle}
      >
        <NotificationsIcon />
      </StyledNotificationsIcon>
      <Popper open={open} anchorEl={anchorRef.current} transition placement="bottom-end">
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            placement={placement === "bottom" ? "center top" : "center bottom"}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList
                  autoFocusItem={open}
                  id="notification-options"
                  onKeyDown={handleListKeyDown}
                >
                  {data?.map(d => {
                    return (
                      <MenuItem>
                        <ListItemText>{d.value}</ListItemText>
                      </MenuItem>
                    );
                  })}
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  );
};

export default Notifications;
