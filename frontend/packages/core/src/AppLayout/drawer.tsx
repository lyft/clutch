import * as React from "react";
import styled from "@emotion/styled";
import { Link as RouterLink } from "react-router-dom";
import {
  Collapse,
  Drawer as MuiDrawer,
  Grow as MuiGrow,
  List,
  ListItem,
  ListItemText,
  Popper,
  Paper as MuiPaper,
  Typography,
} from "@material-ui/core";

import { useTheme } from "../AppProvider/themes";
import { useAppContext } from "../Contexts";

import { routesByGrouping, sortedGroupings } from "./utils";

const mobileDrawerWidth = "90%";
const drawerWidth = "100px";

const DrawerPanel = styled(MuiDrawer)({
  flexShrink: 0,
  minWidth: `${drawerWidth}`,
  "@media screen and (max-width: 800px)": {
    minWidth: `${mobileDrawerWidth}`,
  },
  "div[class*='MuiDrawer-paper']": {
    minWidth: `${drawerWidth}`,
    backgroundColor: "linear-gradient(0deg, #FFFFFF, #FFFFFF), #FFFFFF",
    boxShadow: "0px 5px 15px rgba(53, 72, 212, 0.2)",
    "@media screen and (max-width: 800px)": {
      minWidth: `${mobileDrawerWidth}`,
    }
  }
});

const GroupHeading = styled(Typography)({
  color: "rgba(13, 16, 48, 0.6)",
  paddingTop: "0.25rem",
  fontWeight: 500,
  fontSize: "14px",
  lineHeight: "18px",
  textAlign: "center",
});

const GroupListItem = styled(ListItem)({
  height: "82px",
  "&:hover": {
    backgroundColor: "#EBEDFB",
  },
});

// const NavigationLink = styled(RouterLink)`
//   ${({ theme }) => `
//   color: ${theme.palette.text.secondary};
//   `}
// `;

const Paper = styled(MuiPaper)({
  minWidth: "230px",
  border: "1px solid #E7E7EA",
  boxShadow: "0px 10px 24px rgba(35, 48, 143, 0.3)",
})

const Grow = styled(MuiGrow)((props: { placement: string }) => ({
  transformOrigin: props.placement,
}));

interface GroupProps {
  heading: string;
  open: boolean;
  updateOpenGroup: (heading: string) => void;
  onNavigate: () => void;
}

const Group: React.FC<GroupProps> = ({
  heading,
  open = false,
  updateOpenGroup,
  onNavigate,
  children,
}) => {
  const childrenList = React.Children.toArray(children);

  const [openList, setOpen] = React.useState(false);
  const anchorRef = React.useRef(null);

  const handleToggle = () => {
    setOpen(!openList);
  };

  // n.b. if a Workflow Grouping has no workflows in it don't display it even if
  // it's not explicitly marked as hidden.
  if (childrenList.length === 0) {
    return null;
  }
  return (
    <List data-qa="workflowGroup">
      <GroupListItem
        button
        onClick={() => {
          updateOpenGroup(heading);
          handleToggle();
        }}
        ref={anchorRef}
        aria-controls={open ? "workflow-options" : undefined}
        aria-haspopup="true"
      >
        <GroupHeading>
          {heading}
        </GroupHeading>
      </GroupListItem>
      <Popper
        open={openList}
        anchorEl={anchorRef.current}
        role={undefined}
        transition
        placement="right-start"
      >
        {({ TransitionProps, placement }) => (
        <Grow
          {...TransitionProps}
          placement={placement === "bottom" ? "center top" : "center bottom"}
        >
          <Paper>
            <Collapse in={open} timeout="auto" unmountOnExit>
              <List
                component="div"
                disablePadding
                id="workflow-options"
              >
                {childrenList.map((c: React.ReactElement) => {
                  return React.cloneElement(c, { onClick: onNavigate });
                })}
              </List>
            </Collapse>
          </Paper>
        </Grow>
        )}
      </Popper>
    </List>
  );
};

interface LinkProps {
  to: string;
  text: string;
  onClick: () => void;
}

const Link: React.FC<LinkProps> = ({ to, text, onClick }) => {
  const theme = useTheme();
  const isSelected = window.location.pathname.replace("/", "") === to;
  const selectedStyle = isSelected ? { color: theme.palette.accent.main } : {};
  return (
    <ListItem
      component={RouterLink}
      onClick={onClick}
      to={to}
      dense
      data-qa="workflowGroupItem"
    >
      <ListItemText primary={text} style={selectedStyle} />
    </ListItem>
  );
};

interface DrawerProps {
  onClose: () => void;
}

const Drawer: React.FC<DrawerProps> = ({ onClose }) => {
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
    <DrawerPanel onClose={onClose} data-qa="drawer" variant="permanent">
      {sortedGroupings(workflows).map(grouping => {
        const value = routesByGrouping(workflows)[grouping];
        return (
          <Group
            key={grouping}
            heading={grouping}
            open={openGroup === grouping}
            updateOpenGroup={updateOpenGroup}
            onNavigate={onClose}
          >
            {value.workflows.map(workflow => (
              <Link
                key={workflow.path.replace("/", "")}
                to={workflow.path}
                text={workflow.displayName}
                onClick={onClose}
              />
            ))}
          </Group>
        );
      })}
    </DrawerPanel>
  );
};

export default Drawer;
