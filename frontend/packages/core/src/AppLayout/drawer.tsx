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
  minWidth: `${drawerWidth}`,
  "div[class*='MuiDrawer-paper']": {
    minWidth: `${drawerWidth}`,
    top: "64px",
    backgroundColor: "linear-gradient(0deg, #FFFFFF, #FFFFFF), #FFFFFF",
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
  padding: "16px 8px 0px 8px",
  "&:hover": {
    backgroundColor: "#E7E7EA",
  },
  "&:active": {
    backgroundColor:
      "linear-gradient(0deg, rgba(255, 255, 255, 0.85), rgba(255, 255, 255, 0.85)), #0D1030",
  },
  "&.Mui-selected": {
    backgroundColor: "#EBEDFB",
    ".MuiAvatar-root": {
      color: "#FFFFFF",
      background: "#3548D4",
    },
  },
});

const Avatar = styled(MuiAvatar)({
  color: "rgba(13, 16, 48, 0.38)",
  background: "rgba(13, 16, 48, 0.12)",
  height: "36px",
  width: "36px",
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

// sidebar submenu
const Popper = styled(MuiPopper)({
  zIndex: 1200,
  paddingTop: "16px",
});

const Paper = styled(MuiPaper)({
  minWidth: "230px",
  border: "1px solid #E7E7EA",
  boxShadow: "0px 10px 24px rgba(35, 48, 143, 0.3)",
});

// sudebar submenu groupings
const LinkListItem = styled(ListItem)({
  backgroundColor: "#FFFFFF",
  "&.Mui-selected": {
    backgroundColor: "#EBEDFB",
    ".MuiTypography-root": {
      color: "#3548D4",
    },
  },
  "&:hover": {
    backgroundColor: "#E7E7EA",
  },
  height: "48px",
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
        <Avatar>{heading.charAt(0)}</Avatar>
        <GroupHeading align="center">{heading}</GroupHeading>
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
