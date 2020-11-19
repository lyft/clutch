import React from "react";
import styled from "@emotion/styled";
import {
  Avatar as MuiAvatar,
  Box,
  ClickAwayListener,
  Grow,
  IconButton,
  MenuList,
  Paper,
  Popper,
  Typography,
} from "@material-ui/core";
import Cookies from "js-cookie";
import jwtDecode from "jwt-decode";

const UserPhoto = styled(IconButton)`
  padding: 0.06rem 0rem 0.25rem 0.75rem;
  margin-right: 0.25rem;
`;

const Avatar = styled(MuiAvatar)`
  background-color: #d7dadb;
  height: 2rem;
  width: 2rem;
`;

const Initials = styled(Typography)`
  color: #02acbe;
  font-size: 1rem;
`;

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
    </Box>
  );
};

export { UserInformation, userId };
