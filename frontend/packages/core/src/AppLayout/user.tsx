import React from "react";
import styled from "@emotion/styled";
import {
  Avatar as MuiAvatar,
  Box,
  ClickAwayListener,
  Grow as MuiGrow,
  IconButton,
  ListItemIcon,
  ListItemText,
  MenuItem as MuiMenuItem,
  MenuList,
  Paper as MuiPaper,
  Popper,
  Typography,
} from "@material-ui/core";
import Cookies from "js-cookie";
import jwtDecode from "jwt-decode";

const UserPhoto = styled(IconButton)({
  padding: "12px",
  "&:hover": {
    background: "#2d3db4",
  },
  "&:active": {
    background: "#2938a5",
  },
});

const Avatar = styled(MuiAvatar)({
  backgroundColor: "#dce7f4",
  height: "28px",
  width: "28px",
});

const AvatarBackdrop = styled(MuiAvatar)({
  backgroundColor: "#f6faff",
  height: "32px",
  width: "32px",
});

const Initials = styled(Typography)({
  color: "#0d1030",
  opacity: 0.6,
  fontSize: "14px",
  fontWeight: 500,
});

const Paper = styled(MuiPaper)({
  width: "266px",
  border: "1px solid #e2e2e6",
  boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
});

const UserProfileMenuItem = styled(MuiMenuItem)({
  "&:focus": {
    background: "transparent",
  },
  "&:hover": {
    background: "transparent",
  },
});

const AvatarListItemIcon = styled(ListItemIcon)({
  marginLeft: "8px",
});

const AvatarListItemText = styled(Typography)({
  color: "#0d1030",
  fontSize: "14px",
  opacity: "0.6",
});

const Grow = styled(MuiGrow)((props: { placement: string }) => ({
  transformOrigin: props.placement,
}));

interface JwtToken {
  sub: string;
}

const userId = (): string => {
  // Check JWT token for subject and display if available.
  const token = Cookies.get("token");
  if (!token) {
    return "Anonymous";
  }
  let subject = "Unknown user";
  try {
    const decoded = jwtDecode(token) as JwtToken;
    if (decoded?.sub) {
      subject = decoded.sub;
    }
  } catch {}
  return subject;
};

const UserAvatar: React.FC = () => {
  return (
    <AvatarBackdrop>
      <Avatar>
        <Initials>{userId().slice(0, 2).toUpperCase()}</Initials>
      </Avatar>
    </AvatarBackdrop>
  );
};

const UserInformation: React.FC = () => {
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
      <UserPhoto
        ref={anchorRef}
        edge="end"
        aria-controls={open ? "account-options" : undefined}
        aria-haspopup="true"
        onClick={handleToggle}
      >
        <UserAvatar />
      </UserPhoto>
      <Popper
        open={open}
        anchorEl={anchorRef.current}
        role={undefined}
        transition
        popperOptions={{
          modifiers: {
            offset: {
              offset: "-116, -112",
            },
          },
        }}
      >
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            placement={placement === "bottom" ? "center top" : "center bottom"}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList autoFocusItem={open} id="account-options" onKeyDown={handleListKeyDown}>
                  <UserProfileMenuItem>
                    <AvatarListItemIcon>
                      <UserAvatar />
                    </AvatarListItemIcon>
                    <ListItemText>
                      <AvatarListItemText>{userId()}</AvatarListItemText>
                    </ListItemText>
                  </UserProfileMenuItem>
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </Box>
  );
};

export { UserInformation, userId };
