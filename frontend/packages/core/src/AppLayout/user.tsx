import React from "react";
import styled from "@emotion/styled";
import {
  Avatar as MuiAvatar,
  Box,
  ClickAwayListener,
  Divider as MuiDivider,
  Grow as MuiGrow,
  IconButton,
  ListItemIcon,
  ListItemText as MuiListItemText,
  MenuItem as MuiMenuItem,
  MenuList,
  Paper as MuiPaper,
  Popper as MuiPopper,
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
  ".avatar-header .MuiAvatar-root": {
    height: "32px",
    width: "32px",
  },
});

const Avatar = styled(MuiAvatar)({
  backgroundColor: "#727FE1",
});

const Initials = styled(Typography)({
  color: "#FFFFFF",
  fontSize: "14px",
  fontWeight: 500,
  lineHeight: "18px",
});

const Paper = styled(MuiPaper)({
  width: "242px",
  border: "1px solid #E7E7EA",
  boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
});

const Popper = styled(MuiPopper)({
  padding: "0 12px",
  marginLeft: "12px",
  zIndex: 1101,
});

const MenuItem = styled(MuiMenuItem)({
  "&:hover": {
    backgroundColor: "#E7E7EA",
  },
  "&.MuiListItem-root.Mui-focusVisible": {
    backgroundColor: "#DBDBE0",
  },
  "&:active": {
    backgroundColor: "#EBEDFB",
  },
});

const Divider = styled(MuiDivider)({
  backgroundColor: "#E7E7EA",
});

const AvatarListItemIcon = styled(ListItemIcon)({
  minWidth: "inherit",
  width: "48px",
  ".avatar-menu .MuiAvatar-root": {
    height: "48px",
    width: "48px",
  },
  ".avatar-menu .MuiTypography-root": {
    fontSize: "20px",
    lineHeight: "24px",
  },
});

const AvatarListItemText = styled(MuiListItemText)({
  paddingLeft: "16px",
  ".MuiTypography-root": {
    color: "rgba(13, 16, 48, 0.6)",
    fontSize: "14px",
    lineHeight: "24px",
  },
});

const ListItemText = styled(MuiListItemText)({
  ".MuiTypography-root": {
    color: "#0D1030",
    fontSize: "14px",
    lineHeight: "24px",
  },
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

export interface UserAvatarProps {
  initials: string;
}

const UserAvatar: React.FC<UserAvatarProps> = ({ initials }) => {
  return (
    <Avatar>
      <Initials>{initials}</Initials>
    </Avatar>
  );
};

interface UserData {
  value: string;
  user: string;
}

export interface UserInformationProps {
  data?: UserData[];
  user?: string;
}

// TODO (sperry): investigate using popover instead of popper
const UserInformation: React.FC<UserInformationProps> = ({ data, user = userId() }) => {
  const userInitials = user.slice(0, 2).toUpperCase();
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
        <div className="avatar-header">
          <UserAvatar initials={userInitials} />
        </div>
      </UserPhoto>
      <Popper open={open} anchorEl={anchorRef.current} transition placement="bottom-end">
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            placement={placement === "bottom" ? "center top" : "center bottom"}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList autoFocusItem={open} id="account-options" onKeyDown={handleListKeyDown}>
                  <MenuItem>
                    <AvatarListItemIcon>
                      <div className="avatar-menu">
                        <UserAvatar initials={userInitials} />
                      </div>
                    </AvatarListItemIcon>
                    <AvatarListItemText>{user}</AvatarListItemText>
                  </MenuItem>
                  <Divider />
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
    </Box>
  );
};

export { UserInformation, userId };
