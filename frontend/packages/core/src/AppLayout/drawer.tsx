import React from "react";
import { Link as RouterLink } from "react-router-dom";
import {
  Collapse,
  Divider as MuiDivider,
  Drawer as MuiDrawer,
  Grid,
  List,
  ListItem,
  ListItemText,
  SvgIcon,
  Typography,
} from "@material-ui/core";
import ExpandLess from "@material-ui/icons/ExpandLess";
import ExpandMore from "@material-ui/icons/ExpandMore";
import styled from "styled-components";

import { useTheme } from "../AppProvider/themes";
import { useAppContext } from "../Contexts";
import { TrendingUpIcon } from "../icon";

import Logo from "./logo";
import { userId } from "./user";
import { routesByGrouping, sortedGroupings } from "./utils";

const mobileDrawerWidth = "90%";
const drawerWidth = "300px";

const DrawerPanel = styled(MuiDrawer)`
  ${({ theme }) => `
  flex-shrink: 0;
  min-width: ${drawerWidth};
  @media screen and (max-width: 800px) {
    min-width: ${mobileDrawerWidth};
  }
  div[class*="MuiDrawer-paper"] {
    min-width: ${drawerWidth};
    background-color: ${theme.palette.secondary.main};
    padding: ${theme.spacing(2)}px;
    @media screen and (max-width: 800px) {
      min-width: ${mobileDrawerWidth};
    }
  }
  `}
`;

const DrawerHeader = styled(Grid)`
  ${({ theme }) => `
    justify: flex-start;
    direction: row;
    ${theme.mixins.toolbar};
  `}
`;

const GroupIcon = styled(SvgIcon)`
  ${({ theme }) => `
  color: ${theme.palette.primary.main};
  `}
`;

const GroupHeading = styled(Typography)`
  ${({ theme }) => `
  color: ${theme.palette.primary.main};
  padding-top: 0.25rem;
  font-weight: bolder;
  `}
`;

const TrendingIcon = styled(TrendingUpIcon)`
  fontsize: 20px;
  marginleft: 10px;
`;

const NavigationLink = styled(RouterLink)`
  ${({ theme }) => `
  color: ${theme.palette.text.secondary};
  `}
`;

interface HeaderProps {
  onNavigate: () => void;
}

const Header: React.FC<HeaderProps> = ({ onNavigate }) => {
  return (
    <DrawerHeader container flex-wrap="wrap">
      <Grid item>
        <NavigationLink to="/" onClick={onNavigate} data-qa="logo">
          <Logo />
        </NavigationLink>
      </Grid>
      <Grid item>
        <Grid container justify="center" direction="column">
          <Typography
            style={{ lineHeight: "1.3", margin: "-5px 0px 10px 10px" }}
            component="span"
            color="primary"
            data-qa="title"
          >
            <NavigationLink to="/" onClick={onNavigate}>
              <GroupHeading>clutch</GroupHeading>
            </NavigationLink>
            <Typography variant="caption">Welcome {userId()}!</Typography>
          </Typography>
        </Grid>
      </Grid>
    </DrawerHeader>
  );
};

const Divider = styled(MuiDivider)`
  ${({ theme }) => `
  background-color: ${theme.palette.accent.main};
  `}
`;

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

  // n.b. if a Workflow Grouping has no workflows in it don't display it even if
  // it's not explicitly marked as hidden.
  if (childrenList.length === 0) {
    return null;
  }
  return (
    <List data-qa="workflowGroup">
      <ListItem button onClick={() => updateOpenGroup(heading)}>
        <GroupHeading color="primary" variant="h6">
          {heading}
        </GroupHeading>
        <Grid container justify="flex-end" data-qa="toggle">
          {open ? <GroupIcon component={ExpandLess} /> : <GroupIcon component={ExpandMore} />}
        </Grid>
      </ListItem>
      {!open ? <Divider variant="middle" /> : null}
      <Collapse in={open} timeout="auto" unmountOnExit>
        <List component="div" disablePadding>
          {childrenList.map((c: React.ReactElement) => {
            return React.cloneElement(c, { onClick: onNavigate });
          })}
        </List>
      </Collapse>
    </List>
  );
};

interface LinkProps {
  to: string;
  text: string;
  onClick: () => void;
  trending?: boolean;
}

const Link: React.FC<LinkProps> = ({ to, text, onClick, trending = false }) => {
  const theme = useTheme();
  const isSelected = window.location.pathname.replace("/", "") === to;
  const selectedStyle = isSelected ? { color: theme.palette.accent.main } : {};
  return (
    <ListItem
      component={NavigationLink}
      onClick={onClick}
      to={to}
      dense
      data-qa="workflowGroupItem"
    >
      <ListItemText primary={text} style={selectedStyle} />
      {trending && <TrendingIcon />}
    </ListItem>
  );
};

interface DrawerProps {
  open: boolean;
  onClose: () => void;
}

const Drawer: React.FC<DrawerProps> = ({ open, onClose }) => {
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
    <DrawerPanel open={open} onClose={onClose} data-qa="drawer">
      <Header onNavigate={onClose} />
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
                trending={workflow.trending}
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
