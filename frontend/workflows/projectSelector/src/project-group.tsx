import * as React from "react";
import { Checkbox, Switch } from "@clutch-sh/core";
import styled from "@emotion/styled";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import ClearIcon from "@material-ui/icons/Clear";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";

import type { Group } from "./hello-world";
import { deriveSwitchStatus, useDispatch, useReducerState } from "./hello-world";

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
});

const StyledMenuItem = styled.div({
  display: "flex",
  alignItems: "center",
  "&:hover": {
    backgroundColor: "rgba(13, 16, 48, 0.03)",
  },
  "&:hover > span": {
    display: "inline-flex", // Unhide "only" hidden span.
  },
});

const StyledProjectHeader = styled.div({
  display: "flex",
  alignItems: "center",
  justifyContent: "space-between",
});

const StyledHeaderColumn = styled.div({
  display: "flex",
  alignItems: "center",
});

const StyledAllText = styled.div({
  color: "rgba(13, 16, 48, 0.38)",
});

const StyledMenuItemName = styled.span({
  flexGrow: 1,
});

const StyledClearIcon = styled.span({
  color: "rgba(13, 16, 48, 0.38)",
  display: "inline-flex",
  width: "24px",
});

const StyledOnlyButton = styled.span({
  marginRight: "10px",
  color: "rgba(53, 72, 212, 1)",
});

interface ProjectGroupProps {
  title: string;
  group: Group;
  collapsible?: boolean;
}

const ProjectGroup: React.FC<ProjectGroupProps> = ({ title, group, collapsible }) => {
  const dispatch = useDispatch();
  const state = useReducerState();

  const [collapsed, setCollapsed] = React.useState(false);

  const numProjects = Object.keys(state[group]).length;
  const checkedProjects = Object.keys(state[group]).filter(k => state[group][k].checked);

  return (
    <>
      <StyledProjectHeader>
        <StyledHeaderColumn>
          {collapsible && (
            <StyledHeaderColumn onClick={() => setCollapsed(!collapsed)}>
              {collapsed ? <ChevronRightIcon /> : <ExpandMoreIcon />}
            </StyledHeaderColumn>
          )}
          <StyledProjectTitle>
            {title}
            <StyledCount>
              {checkedProjects.length}
              {numProjects > 0 && "/" + numProjects}
            </StyledCount>
          </StyledProjectTitle>
        </StyledHeaderColumn>
        <StyledHeaderColumn>
          {!collapsible && <StyledAllText>All</StyledAllText>}
          <Switch
            onChange={() =>
              dispatch({
                type: "TOGGLE_ENTIRE_GROUP",
                payload: { group: group },
              })
            }
            checked={deriveSwitchStatus(state, group)}
            disabled={numProjects == 0 || state.loading}
          />
        </StyledHeaderColumn>
      </StyledProjectHeader>
      {!collapsed && (
        <div>
          {numProjects == 0 && <div>No projects in this group.</div>}
          {Object.keys(state[group]).map(key => (
            <StyledMenuItem key={key}>
              <Checkbox
                name={key}
                disabled={state.loading}
                onChange={() =>
                  dispatch({
                    type: "TOGGLE_PROJECTS",
                    payload: { group, projects: [key] },
                  })
                }
                checked={state[group][key].checked ? true : false}
              />
              <StyledMenuItemName>{key}</StyledMenuItemName>
              <StyledOnlyButton
                hidden={true}
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
                  <ClearIcon
                    onClick={() =>
                      !state.loading &&
                      dispatch({
                        type: "REMOVE_PROJECTS",
                        payload: { group, projects: [key] },
                      })
                    }
                  />
                )}
              </StyledClearIcon>
            </StyledMenuItem>
          ))}
        </div>
      )}
    </>
  );
};

export default ProjectGroup;
