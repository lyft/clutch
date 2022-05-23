import * as React from "react";
import { Checkbox, Switch } from "@clutch-sh/core";
import styled from "@emotion/styled";
import IconButton from "@material-ui/core/IconButton";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import ClearIcon from "@material-ui/icons/Clear";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";

import { deriveSwitchStatus, useDispatch, useReducerState } from "./helpers";
import ProjectLinks from "./project-links";
import type { Group } from "./types";

const StyledGroup = styled.div({
  fontWeight: 500,
  marginLeft: "4px",
  marginTop: "9px",
  display: "block",
});

const StyledGroupTitle = styled.span({
  marginRight: "4px",
  display: "inline-block",
});

const StyledCount = styled.span({
  color: "rgba(13, 16, 48, 0.6)",
  backgroundColor: "rgba(13, 16, 48, 0.03)",
  fontVariantNumeric: "tabular-nums",
  borderRadius: "4px",
  fontWeight: "bold",
  fontSize: "12px",
  padding: "1px 4px",
  marginRight: "4px",
  marginBottom: "10px",
  marginTop: "2px",
  display: "inline-block",
});

// This div used to have `padding: "0 25px 0 8px"` but that made it look weird when we implemented quicklinks
// because the "only" and "x" buttons are hidden when the popper is expanded and mouse is no longer hovering.
const StyledMenuItem = styled.div({
  height: "48px",
  display: "flex",
  alignItems: "center",
  "&:hover": {
    backgroundColor: "rgba(13, 16, 48, 0.03)",
  },
  "&:hover > div": {
    display: "inline-flex", // Unhide hidden only button and x if necessary.
  },
});

const StyledProjectHeader = styled.div({
  display: "flex",
  maxWidth: "100%",
  alignItems: "flex-start",
  justifyContent: "space-between",
  minHeight: "40px",
  padding: "0 12px",
});

const StyledHeaderColumn = styled.div((props: { grow?: boolean }) => ({
  display: "flex",
  minHeight: "38px",
  alignItems: "center",
  flexGrow: props.grow ? 1 : 0,
}));

const StyledNoProjectsText = styled.div({
  color: "rgba(13, 16, 48, 0.38)",
  textAlign: "center",
  fontSize: "12px",
  marginBottom: "16px",
});

const StyledAllText = styled.div({
  color: "rgba(13, 16, 48, 0.38)",
});

const StyledMenuItemName = styled.span({
  flexGrow: 1,
  whiteSpace: "nowrap",
  textOverflow: "ellipsis",
  overflow: "hidden",
  marginLeft: "4px",
  marginRight: "0px",
  display: "block",
  maxWidth: "160px",
});

const StyledClearIcon = styled.span({
  ".MuiIconButton-root": {
    padding: "6px",
    color: "rgba(13, 16, 48, 0.38)",
  },
  ".MuiIconButton-root:hover": {
    backgroundColor: "rgb(245, 246, 253)",
  },
  ".MuiIconButton-root:active": {
    backgroundColor: "rgba(0,0,0, 0.1)",
  },
});

const StyledOnlyButton = styled.button({
  border: "none",
  cursor: "pointer",
  borderRadius: "4px",
  fontSize: "14px",
  padding: "5px",
  marginRight: "4px",
  color: "rgba(53, 72, 212, 1)",
  backgroundColor: "unset",
  fontFamily: "inherit",
  "&:hover": {
    backgroundColor: "#f5f6fd",
  },
  "&:active": {
    backgroundColor: "#D7DAF6",
  },
});

const StyledHoverOptions = styled.div({
  alignItems: "center",
});

interface ProjectGroupProps {
  title: string;
  group: Group;
  displayToggleHelperText?: boolean;
}

const ProjectGroup: React.FC<ProjectGroupProps> = ({ title, group, displayToggleHelperText }) => {
  const dispatch = useDispatch();
  const state = useReducerState();

  const [collapsed, setCollapsed] = React.useState(false);

  const groupKeys = Object.keys(state?.[group] ?? {});
  const numProjects = groupKeys.length;
  const checkedProjects = groupKeys.filter(k => state?.[group][k].checked);

  // We need to keep track of which project has its quick links open so that we know
  // to hide the other projects' buttons
  const [quickLinksWindowKey, setQuickLinksWindowKey] = React.useState<string>("");

  const onCloseQuickLinks = () => {
    setQuickLinksWindowKey("");
  };

  const onOpenQuickLinks = (projectName: string) => {
    setQuickLinksWindowKey(projectName);
  };

  return (
    <>
      <StyledProjectHeader>
        <StyledHeaderColumn onClick={() => setCollapsed(!collapsed)}>
          {collapsed ? <ChevronRightIcon /> : <ExpandMoreIcon />}
        </StyledHeaderColumn>
        <StyledHeaderColumn grow>
          <StyledGroup>
            <StyledGroupTitle>{title}</StyledGroupTitle>
            <StyledCount>
              {checkedProjects.length}
              {numProjects > 0 && ` / ${numProjects}`}
            </StyledCount>
          </StyledGroup>
        </StyledHeaderColumn>
        <StyledHeaderColumn>
          {displayToggleHelperText && <StyledAllText>All</StyledAllText>}
          <Switch
            onChange={() =>
              dispatch &&
              dispatch({
                type: "TOGGLE_ENTIRE_GROUP",
                payload: { group, projects: groupKeys },
              })
            }
            checked={deriveSwitchStatus(state, group)}
            disabled={numProjects === 0 || state?.loading}
          />
        </StyledHeaderColumn>
      </StyledProjectHeader>
      {!collapsed && (
        <div>
          {numProjects === 0 && (
            <StyledNoProjectsText>No projects in this group yet.</StyledNoProjectsText>
          )}
          {groupKeys.sort().map(key => (
            <StyledMenuItem key={key}>
              <Checkbox
                name={key}
                size="small"
                disabled={state?.loading}
                onChange={() =>
                  dispatch &&
                  dispatch({
                    type: "TOGGLE_PROJECTS",
                    payload: { group, projects: [key] },
                  })
                }
                checked={!!state?.[group][key].checked}
              />
              <StyledMenuItemName>{key}</StyledMenuItemName>
              <StyledHoverOptions hidden={quickLinksWindowKey !== key}>
                <StyledOnlyButton
                  onClick={() =>
                    !state?.loading &&
                    dispatch &&
                    dispatch({
                      type: "ONLY_PROJECTS",
                      payload: { group, projects: [key] },
                    })
                  }
                >
                  Only
                </StyledOnlyButton>
                <StyledClearIcon>
                  {state?.[group][key].custom && (
                    <IconButton
                      onClick={() =>
                        !state?.loading &&
                        dispatch &&
                        dispatch({
                          type: "REMOVE_PROJECTS",
                          payload: { group, projects: [key] },
                        })
                      }
                    >
                      <ClearIcon />
                    </IconButton>
                  )}
                </StyledClearIcon>
              </StyledHoverOptions>
              {!state?.loading && state?.projectData?.[key]?.linkGroups && (
                <ProjectLinks
                  linkGroups={state?.projectData?.[key]?.linkGroups ?? []}
                  onOpen={() => onOpenQuickLinks(key)}
                  onClose={onCloseQuickLinks}
                  showOpenButton={quickLinksWindowKey !== key}
                />
              )}
            </StyledMenuItem>
          ))}
        </div>
      )}
    </>
  );
};

export default ProjectGroup;
