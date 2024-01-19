import * as React from "react";
import styled from "@emotion/styled";
import type {
  ClickAwayListenerProps,
  ListItemProps,
  PopperProps as MuiPopperProps,
  Theme,
} from "@mui/material";
import {
  alpha,
  ClickAwayListener,
  Collapse,
  List,
  ListItem as MuiListItem,
  ListItemText as MuiListItemText,
  Paper as MuiPaper,
  Popper as MuiPopper,
} from "@mui/material";

const StyledPopper = styled(MuiPopper)({
  zIndex: 1201,
  paddingTop: "16px",
});

const Paper = styled(MuiPaper)(({ theme }: { theme: Theme }) => ({
  minWidth: "fit-content",
  border: `1px solid ${theme.palette.secondary[200]}`,
  boxShadow: `0px 10px 24px ${alpha(theme.palette.primary[700], 0.3)}`,
  ".MuiListItem-root[id='popperItem']": {
    backgroundColor: theme.palette.contrastColor,
    height: "48px",
    "&:hover": {
      backgroundColor: theme.palette.primary[100],
    },
    "&:active": {
      backgroundColor: theme.palette.primary[300],
    },
    "&.Mui-selected": {
      backgroundColor: theme.palette.contrastColor,
      "&:hover": {
        backgroundColor: theme.palette.primary[100],
      },
      "&:active": {
        backgroundColor: theme.palette.primary[300],
      },
    },
    "&:hover, &:active, &.Mui-selected": {
      ".MuiTypography-root": {
        color: theme.palette.primary[600],
      },
    },
  },
}));

const ListItem = styled(MuiListItem)({
  padding: "0",
});

const PopperItemIcon = styled.div({
  margin: "12px 5px 12px 12px",
  height: "24px",
  width: "24px",
});

const ListItemText = styled(MuiListItemText)(({ theme }: { theme: Theme }) => ({
  ".MuiTypography-root": {
    color: alpha(theme.palette.secondary[900], 0.6),
    fontWeight: 500,
    fontSize: "14px",
    lineHeight: "18px",
    padding: "15px 15px",
  },
}));

export interface PopperItemProps extends Pick<ListItemProps, "selected" | "disabled"> {
  children: React.ReactNode;
  component?: React.ElementType;
  componentProps?: any;
  onClick?: () => void;
  icon?: React.ReactElement;
}

const PopperItem = ({ children, componentProps, onClick, icon, ...props }: PopperItemProps) => (
  <ListItem button onClick={onClick} id="popperItem" dense {...props} {...componentProps}>
    {icon && <PopperItemIcon>{icon}</PopperItemIcon>}
    <ListItemText>{children}</ListItemText>
  </ListItem>
);

export interface PopperProps
  extends Pick<ClickAwayListenerProps, "onClickAway">,
    Pick<MuiPopperProps, "placement"> {
  open: boolean;
  anchorRef?: React.MutableRefObject<HTMLElement | null> | null;
  children?: React.ReactElement<PopperItemProps> | React.ReactElement<PopperItemProps>[];
  id?: string;
}
const Popper = ({
  open,
  anchorRef,
  onClickAway,
  placement = "right-start",
  children,
  id,
}: PopperProps) => (
  <Collapse in={open} timeout="auto" unmountOnExit>
    <StyledPopper open={open} anchorEl={anchorRef?.current} placement={placement}>
      <Paper>
        <ClickAwayListener onClickAway={onClickAway}>
          <List component="div" disablePadding id={id}>
            {children}
          </List>
        </ClickAwayListener>
      </Paper>
    </StyledPopper>
  </Collapse>
);

export { Popper, PopperItem };
