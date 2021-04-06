import * as React from "react";
import { Link as RouterLink } from "react-router-dom";
import styled from "@emotion/styled";
import {
  Avatar as MuiAvatar,
  ClickAwayListener,
  Collapse,
  Drawer as MuiDrawer,
  List,
  ListItem,
  ListItemText,
  Paper as MuiPaper,
  Popper as MuiPopper,
  Typography,
} from "@material-ui/core";
import _ from "lodash";

import { useAppContext } from "../Contexts";

import { routesByGrouping, sortedGroupings } from "./utils";

// sidebar
const DrawerPanel = styled(MuiDrawer)({
  width: "100px",
  ".MuiDrawer-paper": {
    top: "unset",
    width: "inherit",
    backgroundColor: "#FFFFFF",
    boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
    position: "relative",
    display: "flex",
  },
});

// sidebar groupings
const GroupList = styled(List)({
  padding: "0px",
});

const GroupListItem = styled(ListItem)({
  flexDirection: "column",
  minHeight: "82px",
  padding: "16px 8px 16px 8px",
  height: "fit-content",
  "&:hover": {
    backgroundColor: "#F5F6FD",
  },
  "&:active": {
    backgroundColor: "#D7DAF6",
  },
  // avatar and label
  "&:hover, &:active, &.Mui-selected": {
    ".MuiAvatar-root": {
      backgroundColor: "#3548D4",
    },
    ".MuiTypography-root": {
      color: "#3548D4",
    },
  },
  "&.Mui-selected": {
    backgroundColor: "#EBEDFB",
    "&:hover": {
      backgroundColor: "#F5F6FD",
    },
    "&:active": {
      backgroundColor: "#D7DAF6",
    },
  },
});

const GroupHeading = styled(Typography)({
  color: "rgba(13, 16, 48, 0.6)",
  fontWeight: 500,
  fontSize: "14px",
  lineHeight: "18px",
  flexGrow: 1,
  paddingTop: "11px",
  width: "100%",
  textOverflow: "ellipsis",
  overflow: "hidden",
});

const Avatar = styled(MuiAvatar)({
  background: "rgba(13, 16, 48, 0.6)",
  height: "24px",
  width: "24px",
  color: "#FFFFFF",
  fontSize: "14px",
  borderRadius: "4px",
});

// sidebar submenu
const Popper = styled(MuiPopper)({
  zIndex: 1201,
  paddingTop: "16px",
});

const Paper = styled(MuiPaper)({
  minWidth: "230px",
  border: "1px solid #E7E7EA",
  boxShadow: "0px 10px 24px rgba(35, 48, 143, 0.3)",
  // sidebar submenu groupings
  ".MuiListItem-root[data-qa='workflowGroupItem']": {
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

const LinkListItemText = styled(ListItemText)({
  ".MuiTypography-root": {
    color: "rgba(13, 16, 48, 0.6)",
    fontWeight: 500,
    fontSize: "14px",
    lineHeight: "18px",
  },
});

interface GroupProps {
  heading: string;
  open: boolean;
  updateOpenGroup: (heading: string) => void;
  closeGroup: () => void;
}

const Group: React.FC<GroupProps> = ({
  heading,
  open = false,
  updateOpenGroup,
  closeGroup,
  children,
}) => {
  const anchorRef = React.useRef(null);

  // n.b. if a Workflow Grouping has no workflows in it don't display it even if
  // it's not explicitly marked as hidden.
  if (React.Children.count(children) === 0) {
    return null;
  }
  // TODO (dschaller): revisit how we handle long groups once we have designs.
  // n.b. this is a stop-gap solution to prevent long groups from looking unreadable.
  return (
    <GroupList data-qa="workflowGroup">
      <GroupListItem
        button
        selected={open}
        ref={anchorRef}
        aria-controls={open ? "workflow-options" : undefined}
        aria-haspopup="true"
        onClick={() => {
          updateOpenGroup(heading);
        }}
      >
        <Avatar>{heading.charAt(0)}</Avatar>
        <GroupHeading align="center">{heading}</GroupHeading>
        <Collapse in={open} timeout="auto" unmountOnExit>
          <Popper open={open} anchorEl={anchorRef.current} transition placement="right-start">
            <Paper>
              <ClickAwayListener onClickAway={closeGroup}>
                <List component="div" disablePadding id="workflow-options">
                  {children}
                </List>
              </ClickAwayListener>
            </Paper>
          </Popper>
        </Collapse>
      </GroupListItem>
    </GroupList>
  );
};

interface LinkProps {
  to: string;
  text: string;
}

const Link: React.FC<LinkProps> = ({ to, text }) => {
  const isSelected = window.location.pathname.replace("/", "") === to;
  return (
    <ListItem
      selected={isSelected}
      component={RouterLink}
      to={to}
      dense
      data-qa="workflowGroupItem"
    >
      <LinkListItemText>{text}</LinkListItemText>
    </ListItem>
  );
};

const Drawer: React.FC = () => {
  const { workflows } = useAppContext();
  const [openGroup, setOpenGroup] = React.useState("");

  const updateOpenGroup = (group: string) => {
    if (openGroup === group) {
      setOpenGroup("");
    } else {
      setOpenGroup(group);
    }
  };

  return (
    <DrawerPanel data-qa="drawer" variant="permanent">
      {sortedGroupings(workflows).map(grouping => {
        const value = routesByGrouping(workflows)[grouping];
        const sortedWorkflows = _.sortBy(value.workflows, w => w.displayName);
        return (
          <Group
            key={grouping}
            heading={grouping}
            open={openGroup === grouping}
            updateOpenGroup={updateOpenGroup}
            closeGroup={() => setOpenGroup("")}
          >
            {sortedWorkflows.map(workflow => (
              <Link
                key={workflow.path.replace("/", "")}
                to={workflow.path}
                text={workflow.displayName}
              />
            ))}
          </Group>
        );
      })}
    </DrawerPanel>
  );
};

export default Drawer;
