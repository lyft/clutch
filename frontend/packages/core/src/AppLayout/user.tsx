import React from "react";
import {
  Avatar as MuiAvatar,
  ClickAwayListener,
  Grid,
  Grow,
  IconButton,
  MenuList,
  Paper,
  Popper,
  Typography,
} from "@material-ui/core";
import Cookies from "js-cookie";
import jwtDecode from "jwt-decode";
import styled from "styled-components";

const UserPhoto = styled(IconButton)`
  ${({ theme }) => `
  padding-top: 1px;
  padding-bottom: ${theme.spacing(0.5)}px;
  `}
`;

const Avatar = styled(MuiAvatar)`
  ${({ theme }) => `
  background-color: ${theme.palette.text.secondary}
  `}
`;

const Initials = styled(Typography)`
  ${({ theme }) => `
  color: ${theme.palette.accent.main}
  `}
`;

const userId = (): string => {
  // Check JWT token for subject and display if available.
  const token = Cookies.get("token");
  if (!token) {
    return "Anonymous";
  }
  let subject = "Unknown user";
  try {
    const decoded = jwtDecode(token);
    if (decoded?.sub) {
      subject = decoded.sub;
    }
  } catch {}
  return subject;
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
    <Grid container alignItems="center" justify="flex-end">
      <div>{userId()}</div>
      <UserPhoto
        ref={anchorRef}
        edge="end"
        aria-controls={open ? "account-options" : undefined}
        aria-haspopup="true"
        onClick={handleToggle}
      >
        <Avatar>
          <Initials variant="h6">{userId().slice(0, 2).toUpperCase()}</Initials>
        </Avatar>
      </UserPhoto>
      <Popper open={open} anchorEl={anchorRef.current} role={undefined} transition>
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{ transformOrigin: placement === "bottom" ? "center top" : "center bottom" }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList autoFocusItem={open} id="account-options" onKeyDown={handleListKeyDown} />
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </Grid>
  );
};

export { UserInformation, userId };
