import React from "react";
import styled from "@emotion/styled";
import {
  Avatar as MuiAvatar,
  ClickAwayListener,
  Divider as MuiDivider,
  Grow as MuiGrow,
  IconButton,
  ListItemIcon,
  ListItemText as MuiListItemText,
  MenuItem as MuiMenuItem,
  MenuList as MuiMenuList,
  Paper as MuiPaper,
  Popper as MuiPopper,
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
  // avatar on header
  ".MuiAvatar-root": {
    height: "32px",
    width: "32px",
    fontSize: "14px",
    lineHeight: "18px",
  },
});

// header and menu avatar
const Avatar = styled(MuiAvatar)({
  backgroundColor: "#727FE1",
  color: "#FFFFFF",
  fontWeight: 500,
});

const Paper = styled(MuiPaper)({
  width: "242px",
  border: "1px solid #E7E7EA",
  boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
});

const Popper = styled(MuiPopper)({
  padding: "0 12px",
  marginLeft: "12px",
  zIndex: 1201,
});

const MenuList = styled(MuiMenuList)({
  padding: "0px",
  borderRadius: "4px",
  ".MuiMenuItem-root": {
    "&:hover": {
      backgroundColor: "#E7E7EA",
    },
    "&:active": {
      backgroundColor: "#EBEDFB",
    },
  },
});

// user details menu item
const AvatarMenuItem = styled(MuiMenuItem)({
  height: "52px",
  margin: "16px 0 16px 0",
  padding: "0 16px 0 16px",
});

const AvatarListItemIcon = styled(ListItemIcon)({
  minWidth: "inherit",
  width: "48px",
  // avatar on menu
  ".MuiAvatar-root": {
    height: "48px",
    width: "48px",
    fontSize: "20px",
    lineHeight: "24px",
  },
});

const AvatarListItemText = styled(MuiListItemText)({
  paddingLeft: "16px",
  margin: "0px",
  ".MuiTypography-root": {
    color: "rgba(13, 16, 48, 0.6)",
    fontSize: "14px",
    lineHeight: "24px",
  },
});

// default menu items
const MenuItem = styled(MuiMenuItem)({
  height: "48px",
  padding: "12px",
});

const ListItemText = styled(MuiListItemText)({
  margin: "0px",
  ".MuiTypography-root": {
    color: "#0D1030",
    fontSize: "14px",
    lineHeight: "24px",
  },
});

const Divider = styled(MuiDivider)({
  backgroundColor: "#E7E7EA",
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

interface UserAvatarProps {
  initials: string;
}

const UserAvatar: React.FC<UserAvatarProps> = ({ initials }) => {
  return <Avatar>{initials}</Avatar>;
};

interface UserData {
  value: string;
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
    <>
      <UserPhoto
        ref={anchorRef}
        edge="end"
        aria-controls={open ? "account-options" : undefined}
        aria-haspopup="true"
        onClick={handleToggle}
      >
        <UserAvatar initials={userInitials} />
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
                  <AvatarMenuItem>
                    <AvatarListItemIcon>
                      <UserAvatar initials={userInitials} />
                    </AvatarListItemIcon>
                    <AvatarListItemText>{user}</AvatarListItemText>
                  </AvatarMenuItem>
                  {data?.length === 0 ? null : <Divider />}
                  {data?.map(d => (
                    <MenuItem>
                      <ListItemText>{d.value}</ListItemText>
                    </MenuItem>
                  ))}
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </>
  );
};

export { UserInformation, userId };
