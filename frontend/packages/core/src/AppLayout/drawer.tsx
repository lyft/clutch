import * as React from "react";
import { Link as RouterLink, useLocation } from "react-router-dom";
import styled from "@emotion/styled";
import {
  alpha,
  Avatar as MuiAvatar,
  Drawer as MuiDrawer,
  List,
  ListItemButton,
  Theme,
  Typography,
} from "@mui/material";
import _ from "lodash";

import type { WorkflowIcon } from "../AppProvider";
import type { Workflow } from "../AppProvider/workflow";
import { useAppContext } from "../Contexts";
import { useNavigate } from "../navigation";
import type { PopperItemProps } from "../popper";
import { Popper, PopperItem } from "../popper";
import { THEME_VARIANTS } from "../Theme/colors";

import { filterHiddenRoutes, routesByGrouping, sortedGroupings, workflowByRoute } from "./utils";

// sidebar
const DrawerPanel = styled(MuiDrawer)(({ theme }: { theme: Theme }) => ({
  width: "100px",
  overflowY: "auto",
  ".MuiDrawer-paper": {
    top: "unset",
    width: "inherit",
    backgroundColor:
      theme.palette.mode === THEME_VARIANTS.light
        ? theme.palette.contrastColor
        : theme.palette.background.paper,
    boxShadow: `0px 5px 15px ${alpha(theme.palette.primary[400], 0.2)}`,
    position: "relative",
    display: "flex",
  },
}));

// sidebar groupings
const GroupList = styled(List)({
  padding: "0px",
});

const GroupListItem = styled(ListItemButton)<{ icon: number }>(
  ({ theme }: { theme: Theme }) => ({
    flexDirection: "column",
    minHeight: "82px",
    padding: "16px 8px 16px 8px",
    height: "fit-content",
    "&:hover": {
      backgroundColor: theme.palette.primary[100],
    },
    "&:active": {
      backgroundColor: theme.palette.primary[300],
    },
    "&.Mui-selected": {
      backgroundColor: theme.palette.primary[200],
      "&:hover": {
        backgroundColor: theme.palette.primary[100],
      },
      "&:active": {
        backgroundColor: theme.palette.primary[300],
      },
    },
  }),
  props => ({
    // avatar and label
    "&:hover, &:active, &.Mui-selected": {
      ".MuiAvatar-root": {
        backgroundColor: props.icon ? "unset" : props.theme.palette.primary[600],
      },
      ".MuiTypography-root": {
        color: props.theme.palette.primary[600],
      },
    },
  })
);

const GroupHeading = styled(Typography)(({ theme }: { theme: Theme }) => ({
  color: alpha(theme.palette.secondary[900], 0.6),
  fontWeight: 500,
  fontSize: "14px",
  lineHeight: "18px",
  flexGrow: 1,
  paddingTop: "11px",
  width: "100%",
  textOverflow: "ellipsis",
  overflow: "hidden",
}));

const IconAvatar = styled(MuiAvatar)({
  height: "24px",
  width: "24px",
});

const Avatar = styled(IconAvatar)(({ theme }: { theme: Theme }) => ({
  background: alpha(theme.palette.secondary[900], 0.6),
  color: theme.palette.contrastColor,
  fontSize: "14px",
  borderRadius: "4px",
}));

interface GroupProps {
  heading: string;
  open: boolean;
  selected: boolean;
  icon: WorkflowIcon;
  updateOpenGroup: (heading: string) => void;
  closeGroup: () => void;
  children: React.ReactElement<PopperItemProps> | React.ReactElement<PopperItemProps>[];
}

const Group = ({
  heading,
  open = false,
  selected = false,
  icon,
  updateOpenGroup,
  closeGroup,
  children,
}: GroupProps) => {
  const navigate = useNavigate();
  const anchorRef = React.useRef(null);
  const singleChild = React.Children.count(children) === 1;
  const validIcon = icon.path && icon.path.length > 0;

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
        selected={open || selected}
        ref={anchorRef}
        aria-controls={open ? "workflow-options" : undefined}
        aria-haspopup="true"
        icon={+validIcon}
        onClick={() => {
          if (singleChild) {
            navigate(children[0].props.to);
          }
          updateOpenGroup(heading);
        }}
      >
        {validIcon ? (
          <IconAvatar src={icon.path}>{heading.charAt(0)}</IconAvatar>
        ) : (
          <Avatar>{heading.charAt(0)}</Avatar>
        )}
        <GroupHeading align="center">{heading}</GroupHeading>
        <Popper
          open={open && !singleChild}
          onClickAway={closeGroup}
          anchorRef={anchorRef}
          id="workflow-options"
        >
          {children}
        </Popper>
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
    <PopperItem
      selected={isSelected}
      component={RouterLink}
      componentProps={{ to }}
      data-qa="workflowGroupItem"
    >
      {text}
    </PopperItem>
  );
};

const Drawer: React.FC = () => {
  const { workflows } = useAppContext();
  const filteredWorkflows = filterHiddenRoutes(workflows);
  const [activeWorkflow, setActiveWorkflow] = React.useState<Workflow>(null);
  const [openGroup, setOpenGroup] = React.useState("");
  const location = useLocation();

  const updateOpenGroup = (group: string) => {
    setOpenGroup(openGroup === group ? "" : group);
  };

  React.useEffect(() => {
    setActiveWorkflow(workflowByRoute(filteredWorkflows, location.pathname));
  }, [location]);

  // Will hide the drawer if there are no visible workflows
  if (!filteredWorkflows.length) {
    return null;
  }

  return (
    <DrawerPanel data-qa="drawer" variant="permanent">
      {sortedGroupings(filteredWorkflows).map(grouping => {
        const value = routesByGrouping(filteredWorkflows)[grouping];
        const sortedWorkflows = _.sortBy(value.workflows, w => w.displayName);
        return (
          <Group
            key={grouping}
            heading={grouping}
            open={openGroup === grouping}
            selected={openGroup === grouping || activeWorkflow?.group === grouping}
            icon={value.icon}
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
