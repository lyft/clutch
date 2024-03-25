import React from "react";
import styled from "@emotion/styled";
import {
  alpha,
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
  Theme,
} from "@mui/material";
import Cookies from "js-cookie";
import jwtDecode from "jwt-decode";
import * as _ from "lodash";

const UserPhoto = styled(IconButton)(({ theme }: { theme: Theme }) => ({
  padding: "12px",
  "&:hover": {
    background: theme.palette.primary[600],
  },
  "&:active": {
    background: theme.palette.primary[700],
  },
  // avatar on header
  ".MuiAvatar-root": {
    height: "32px",
    width: "32px",
    fontSize: "14px",
    lineHeight: "18px",
  },
}));

// header and menu avatar
const Avatar = styled(MuiAvatar)(({ theme }: { theme: Theme }) => ({
  backgroundColor: theme.palette.primary[500],
  color: theme.palette.contrastColor,
  fontWeight: 500,
}));

const Paper = styled(MuiPaper)(({ theme }: { theme: Theme }) => ({
  width: "242px",
  border: `1px solid ${theme.palette.secondary[100]}`,
  boxShadow: `0px 5px 15px ${alpha(theme.palette.primary[600], 0.2)}`,
}));

const Popper = styled(MuiPopper)({
  padding: "0 12px",
  offset: "12px",
  zIndex: 1201,
});

const MenuList = styled(MuiMenuList)(({ theme }: { theme: Theme }) => ({
  padding: "0px",
  borderRadius: "4px",
  ".MuiMenuItem-root": {
    "&:hover": {
      backgroundColor: theme.palette.secondary[200],
    },
    "&:active": {
      backgroundColor: theme.palette.primary[200],
    },
  },
}));

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

const AvatarListItemText = styled(MuiListItemText)(({ theme }: { theme: Theme }) => ({
  paddingLeft: "16px",
  margin: "0px",
  wordBreak: "break-all",
  textWrap: "wrap",
  ".MuiTypography-root": {
    color: alpha(theme.palette.secondary[900], 0.9),
    fontSize: "14px",
    lineHeight: "24px",
  },
}));

// default menu items
const MenuItem = styled(MuiMenuItem)({
  height: "48px",
  padding: "12px",
});

const ListItemText = styled(MuiListItemText)(({ theme }: { theme: Theme }) => ({
  margin: "0px",
  ".MuiTypography-root": {
    color: theme.palette.secondary[900],
    fontSize: "14px",
    lineHeight: "24px",
  },
}));

const Divider = styled(MuiDivider)(({ theme }: { theme: Theme }) => ({
  backgroundColor: theme.palette.secondary[100],
}));

const Grow = styled(MuiGrow)((props: { placement: string }) => ({
  transformOrigin: props.placement,
}));

interface JwtToken {
  sub: string;
}

const userId = (): string => {
  if (process.env.NODE_ENV === "development") {
    if (process.env.REACT_APP_USER_ID) {
      return process.env.REACT_APP_USER_ID;
    }
  }
  // Check JWT token for subject and display if available.
  const token = Cookies.get("token");
  if (!token) {
    // eslint-disable-next-line
    console.info("No user token set in development - returning Anonymous");
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
const UserInformation: React.FC<UserInformationProps> = ({
  data,
  user = userId(),
  children = [],
}) => {
  const userInitials = user.slice(0, 2).toUpperCase();
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef(null);

  const handleToggle = () => {
    setOpen(!open);
  };

  const handleClose = event => {
    if (event.target.localName === "body") {
      return;
    }
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
                  {data?.map((d, i) => (
                    // eslint-disable-next-line react/no-array-index-key
                    <React.Fragment key={i}>
                      <MenuItem>
                        <ListItemText>{d.value}</ListItemText>
                      </MenuItem>
                      {i > 0 && i < data.length && <Divider />}
                    </React.Fragment>
                  ))}
                  {_.castArray(children).length > 0 && <Divider />}
                  <div style={{ marginBottom: "8px" }}>
                    {_.castArray(children)?.map((c, i) => (
                      <>
                        <MenuItem>{c}</MenuItem>
                        {i < _.castArray(children).length - 1 && <Divider />}
                      </>
                    ))}
                  </div>
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
