import * as React from "react";
import { Checkbox, Switch } from "@clutch-sh/core";
import styled from "@emotion/styled";
import IconButton from "@material-ui/core/IconButton";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import ClearIcon from "@material-ui/icons/Clear";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";

import { deriveSwitchStatus, useDispatch, useReducerState } from "./helpers";
import type { Group } from "./types";

const StyledCount = styled.span({
  color: "rgba(13, 16, 48, 0.6)",
  backgroundColor: "rgba(13, 16, 48, 0.03)",
  fontVariantNumeric: "tabular-nums",
  letterSpacing: "2px",
  borderRadius: "4px",
  fontWeight: "bold",
  fontSize: "12px",
  padding: "1px 2px 1px 4px",
  margin: "0 4px",
});

const StyledProjectTitle = styled.span({
  fontWeight: 500,
  marginLeft: "4px",
});

const StyledMenuItem = styled.div({
  padding: "0 25px 0 8px",
  height: "48px",
  display: "flex",
  alignItems: "center",
  "&:hover": {
    backgroundColor: "rgba(13, 16, 48, 0.03)",
  },
  "&:hover > span": {
    display: "inline-flex", // Unhide hidden button spans.
  },
});

const StyledProjectHeader = styled.div({
  display: "flex",
  alignItems: "center",
  justifyContent: "space-between",
  height: "48px",
  padding: "0 12px",
});

const StyledHeaderColumn = styled.div({
  display: "flex",
  alignItems: "center",
});

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
  marginLeft: "8px",
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

const StyledHoverOptions = styled.span({
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

  const numProjects = Object.keys(state[group]).length;
  const checkedProjects = Object.keys(state[group]).filter(k => state[group][k].checked);

  return (
    <>
      <StyledProjectHeader>
        <StyledHeaderColumn>
          <StyledHeaderColumn onClick={() => setCollapsed(!collapsed)}>
            {collapsed ? <ChevronRightIcon /> : <ExpandMoreIcon />}
          </StyledHeaderColumn>
          <StyledProjectTitle>
            {title}
            <StyledCount>
              {checkedProjects.length}
              {numProjects > 0 && `/${numProjects}`}
            </StyledCount>
          </StyledProjectTitle>
        </StyledHeaderColumn>
        <StyledHeaderColumn>
          {displayToggleHelperText && <StyledAllText>All</StyledAllText>}
          <Switch
            onChange={() =>
              dispatch({
                type: "TOGGLE_ENTIRE_GROUP",
                payload: { group },
              })
            }
            checked={deriveSwitchStatus(state, group)}
            disabled={numProjects === 0 || state.loading}
          />
        </StyledHeaderColumn>
      </StyledProjectHeader>
      {!collapsed && (
        <div>
          {numProjects === 0 && (
            <StyledNoProjectsText>No projects in this group yet.</StyledNoProjectsText>
          )}
          {Object.keys(state[group])
            .sort()
            .map(key => (
              <StyledMenuItem key={key}>
                <Checkbox
                  name={key}
                  size="small"
                  disabled={state.loading}
                  onChange={() =>
                    dispatch({
                      type: "TOGGLE_PROJECTS",
                      payload: { group, projects: [key] },
                    })
                  }
                  checked={!!state[group][key].checked}
                />
                <StyledMenuItemName>{key}</StyledMenuItemName>
                <StyledHoverOptions hidden>
                  <StyledOnlyButton
                    onClick={() =>
                      !state.loading &&
                      dispatch({
                        type: "ONLY_PROJECTS",
                        payload: { group, projects: [key] },
                      })
                    }
                  >
                    Only
                  </StyledOnlyButton>
                  <StyledClearIcon>
                    {state[group][key].custom && (
                      <IconButton
                        onClick={() =>
                          !state.loading &&
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
              </StyledMenuItem>
            ))}
        </div>
      )}
    </>
  );
};

export default ProjectGroup;
