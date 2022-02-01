import * as React from "react";
import { Checkbox, Switch } from "@clutch-sh/core";
import styled from "@emotion/styled";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import ClearIcon from "@mui/icons-material/Clear";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import IconButton from "@mui/material/IconButton";

import { deriveSwitchStatus, useDispatch, useReducerState } from "./helpers";
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

const StyledMenuItem = styled.div({
  padding: "0 25px 0 8px",
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
  width: "24px",
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
              <StyledHoverOptions hidden>
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
                      size="large"
                    >
                      <ClearIcon />
                    </IconButton>
                  )}
                </StyledClearIcon>
              </StyledHoverOptions>
            </StyledMenuItem>
          ))}
        </div>
      )}
    </>
  );
};

export default ProjectGroup;
