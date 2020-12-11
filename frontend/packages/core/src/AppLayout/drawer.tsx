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

import { useAppContext } from "../Contexts";

import { routesByGrouping, sortedGroupings } from "./utils";

// sidebar
const drawerWidth = "100px";

const DrawerPanel = styled(MuiDrawer)({
  flexShrink: 0,
  width: `${drawerWidth}`,
  "div[class*='MuiDrawer-paper']": {
    width: `${drawerWidth}`,
    top: "64px",
    backgroundColor: "#FFFFFF",
    boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
  },
});

// sidebar groupings
const GroupList = styled(List)({
  padding: "0px",
});

const GroupListItem = styled(ListItem)({
  flexDirection: "column",
  height: "82px",
  padding: "16px 8px 16px 8px",
  "&:hover": {
    backgroundColor: "#F5F6FD",
    "&.Mui-selected": {
      backgroundColor: "#F5F6FD",
    },
  },
  "&:active": {
    backgroundColor: "#D7DAF6",
    "&.Mui-selected": {
      backgroundColor: "#D7DAF6",
    },
  },
  "&:hover, &:active, &.Mui-selected": {
    ".MuiAvatar-root": {
      backgroundColor: "#3548D4",
    },
    ".heading": {
      color: "#3548D4",
    },
    ".initials": {
      color: "#FFFFFF",
    },
  },
  "&.Mui-selected": {
    backgroundColor: "#FFFFFF",
  },
});

const GroupHeading = styled(Typography)({
  color: "rgba(13, 16, 48, 0.6)",
  fontWeight: 500,
  fontSize: "14px",
  lineHeight: "18px",
  flexGrow: 1,
  paddingTop: "8px",
  paddingBottom: "8px",
});

const Avatar = styled(MuiAvatar)({
  background: "rgba(13, 16, 48, 0.6)",
  height: "24px",
  width: "24px",
});

const Initials = styled(Typography)({
  color: "#FFFFFF",
  fontSize: "14px",
});

// sidebar submenu
const Popper = styled(MuiPopper)({
  zIndex: 1200,
  paddingTop: "16px",
});

const Paper = styled(MuiPaper)({
  width: "230px",
  border: "1px solid #E7E7EA",
  boxShadow: "0px 10px 24px rgba(35, 48, 143, 0.3)",
});

// sudebar submenu groupings
const LinkListItem = styled(ListItem)({
  backgroundColor: "#FFFFFF",
  height: "48px",
  "&:hover": {
    backgroundColor: "#F5F6FD",
    "&.Mui-selected": {
      backgroundColor: "#F5F6FD",
    },
  },
  "&:active": {
    backgroundColor: "#D7DAF6",
    "&.Mui-selected": {
      backgroundColor: "#D7DAF6",
    },
  },
  "&:hover, &:active, &.Mui-selected": {
    ".MuiTypography-root": {
      color: "#3548D4",
    },
  },
  "&.Mui-selected": {
    backgroundColor: "#FFFFFF",
  },
});

const LinkHeading = styled(Typography)({
  color: "rgba(13, 16, 48, 0.6)",
  fontWeight: 500,
  fontSize: "14px",
  lineHeight: "18px",
});

interface GroupProps {
  heading: string;
  open: boolean;
  updateOpenGroup: (heading: string) => void;
}

const Group: React.FC<GroupProps> = ({ heading, open = false, updateOpenGroup, children }) => {
  const childrenList = React.Children.toArray(children);

  const [openList, setListOpen] = React.useState(false);

  const anchorRef = React.useRef(null);

  const handleToggle = () => {
    setListOpen(!openList);
  };

  const handleClose = event => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }
    setListOpen(false);
  };

  let headingPath = window.location.pathname.replace("/", "").split("/")[0];

  if (headingPath === "ec2") {
    headingPath = "aws";
  }

  const isSelected = headingPath == heading.toLowerCase();

  // n.b. if a Workflow Grouping has no workflows in it don't display it even if
  // it's not explicitly marked as hidden.
  if (childrenList.length === 0) {
    return null;
  }

  return (
    <GroupList data-qa="workflowGroup">
      <GroupListItem
        button
        selected={isSelected}
        ref={anchorRef}
        aria-controls={openList ? "workflow-options" : undefined}
        aria-haspopup="true"
        onClick={() => {
          handleToggle();
          updateOpenGroup(heading);
        }}
      >
        <Avatar>
          <Initials className="initials">
            {heading.charAt(0)}
          </Initials>
        </Avatar>
        <GroupHeading className="heading" align="center">{heading}</GroupHeading>
        <Collapse in={open} timeout="auto" unmountOnExit>
          <Popper
            open={openList}
            anchorEl={anchorRef.current}
            role={undefined}
            transition
            placement="right-start"
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <List component="div" disablePadding id="workflow-options">
                  {childrenList.map((c: React.ReactElement) => {
                    return React.cloneElement(c);
                  })}
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
    <LinkListItem
      selected={isSelected}
      component={RouterLink}
      to={to}
      dense
      data-qa="workflowGroupItem"
    >
      <ListItemText>
        <LinkHeading>{text}</LinkHeading>
      </ListItemText>
    </LinkListItem>
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
        return (
          <Group
            key={grouping}
            heading={grouping}
            open={openGroup === grouping}
            updateOpenGroup={updateOpenGroup}
          >
            {value.workflows.map(workflow => (
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
