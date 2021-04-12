import * as React from "react";
import styled from "@emotion/styled";
import type {
  ClickAwayListenerProps,
  ListItemProps,
  PopperProps as MuiPopperProps,
} from "@material-ui/core";
import {
  ClickAwayListener,
  Collapse,
  List,
  ListItem as MuiListItem,
  ListItemText as MuiListItemText,
  Paper as MuiPaper,
  Popper as MuiPopper,
} from "@material-ui/core";

const StyledPopper = styled(MuiPopper)({
  zIndex: 1201,
  paddingTop: "16px",
});

const Paper = styled(MuiPaper)({
  minWidth: "fit-content",
  border: "1px solid #E7E7EA",
  boxShadow: "0px 10px 24px rgba(35, 48, 143, 0.3)",
  ".MuiListItem-root[id='popperItem']": {
    backgroundColor: "#FFFFFF",
    height: "48px",
    "&:hover": {
      backgroundColor: "#F5F6FD",
    },
    "&:active": {
      backgroundColor: "#D7DAF6",
    },
    "&.Mui-selected": {
      backgroundColor: "#FFFFFF",
      "&:hover": {
        backgroundColor: "#F5F6FD",
      },
      "&:active": {
        backgroundColor: "#D7DAF6",
      },
    },
    "&:hover, &:active, &.Mui-selected": {
      ".MuiTypography-root": {
        color: "#3548D4",
      },
    },
  },
});

const ListItem = styled(MuiListItem)({
  padding: "0",
});

const PopperItemIcon = styled.div({
  margin: "12px 5px 12px 12px",
  height: "24px",
  width: "24px",
});

const ListItemText = styled(MuiListItemText)({
  ".MuiTypography-root": {
    color: "rgba(13, 16, 48, 0.6)",
    fontWeight: 500,
    fontSize: "14px",
    lineHeight: "18px",
    padding: "15px 15px",
  },
});

export interface PopperItemProps extends Pick<ListItemProps, "selected"> {
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
  anchorRef: React.MutableRefObject<HTMLElement>;
  children: React.ReactElement<PopperItemProps> | React.ReactElement<PopperItemProps>[];
}
const Popper = ({
  open,
  anchorRef,
  onClickAway,
  placement = "right-start",
  children,
}: PopperProps) => (
  <Collapse in={open} timeout="auto" unmountOnExit>
    <StyledPopper open={open} anchorEl={anchorRef.current} transition placement={placement}>
      <Paper>
        <ClickAwayListener onClickAway={onClickAway}>
          <List component="div" disablePadding id="workflow-options">
            {children}
          </List>
        </ClickAwayListener>
      </Paper>
    </StyledPopper>
  </Collapse>
);

export { Popper, PopperItem };
