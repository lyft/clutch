import React from "react";
import styled from "@emotion/styled";
import {
  Avatar as MuiAvatar,
  Box,
  ClickAwayListener,
  Grow,
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

const UserPhoto = styled(IconButton)`
  padding: 0.06rem 0rem 0rem 0.75rem;
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

const ItemText = styled(Typography)`
  color: #0d1030;
  font-size: 0.88rem;
  opacity: 0.6;
`;

const Paper = styled(MuiPaper)`
  width: 16.63rem;
  border: 0.06rem solid #e2e2e6;
  box-shadow: 0rem 0.31rem 0.94rem rgba(53, 72, 212, 0.2);
  border: 0.06rem solid #e2e2e6;
`;

const AvatarMenuItem = styled(MuiMenuItem)`
  &:focus {
    background: transparent;
  }
  &:hover {
    background: transparent;
  }
`;

const AvatarListItemIcon = styled(ListItemIcon)`
  margin-left: 0.5rem;
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

const UserAvatar: React.FC = () => {
  return (
    <Avatar>
      <Initials>{userId().slice(0, 2).toUpperCase()}</Initials>
    </Avatar>
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
        onMouseEnter={handleToggle}
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
              offset: "-115,0",
            },
          },
        }}
      >
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{ transformOrigin: placement === "bottom" ? "center top" : "center bottom" }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList autoFocusItem={open} id="account-options" onKeyDown={handleListKeyDown}>
                  <AvatarMenuItem>
                    <AvatarListItemIcon>
                      <UserAvatar />
                    </AvatarListItemIcon>
                    <ListItemText>
                      <ItemText>{userId()}</ItemText>
                    </ListItemText>
                  </AvatarMenuItem>
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
