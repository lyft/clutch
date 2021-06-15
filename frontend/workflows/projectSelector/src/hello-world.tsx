import * as React from "react";

import _ from "lodash";

import { Button, Checkbox, Switch, TextField } from "@clutch-sh/core";
import { Divider, IconButton, LinearProgress } from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import ExpandLessIcon from "@material-ui/icons/ExpandLess";
import HelpOutlineIcon from "@material-ui/icons/HelpOutline";
import ClearIcon from "@material-ui/icons/Clear";
import LayersIcon from "@material-ui/icons/Layers";

import styled from "@emotion/styled";

enum Group {
  PROJECTS,
  UPSTREAM,
  DOWNSTREAM,
}

type UserActionKind =
  | "ADD_PROJECTS"
  | "REMOVE_PROJECTS"
  | "TOGGLE_PROJECTS"
  | "TOGGLE_ENTIRE_GROUP"
  | "ONLY_PROJECTS";

interface UserAction {
  type: UserActionKind;
  payload: UserPayload;
}

interface UserPayload {
  group: Group;
  projects?: string[];
}

type BackgroundActionKind = "HYDRATE_START" | "HYDRATE_END";

interface BackgroundAction {
  type: BackgroundActionKind;
  payload?: BackgroundPayload;
}

interface BackgroundPayload {
  result: any;
}

type Action = BackgroundAction | UserAction;

interface State {
  [Group.PROJECTS]: GroupState;
  [Group.UPSTREAM]: GroupState;
  [Group.DOWNSTREAM]: GroupState;

  projectData: { [projectName: string]: Project };
  loading: boolean;
}

// TODO: subout with full manifest structure (from proto def)
interface Project {
  upstreams: string[];
  downstreams: string[];
}

interface GroupState {
  [projectName: string]: ProjectState;
}

interface ProjectState {
  checked: boolean;
  // TODO: hidden should be derived?
  hidden?: boolean; // upstreams and downstreams are hidden when their parent is unchecked unless other parents also use them.
  custom?: boolean;
}

const StateContext = React.createContext<State | undefined>(undefined);
const useReducerState = () => {
  return React.useContext(StateContext);
};

const DispatchContext = React.createContext<(action: Action) => void | undefined>(undefined);
const useDispatch = () => {
  return React.useContext(DispatchContext);
};

const selectorReducer = (state: State, action: Action): State => {
  switch (action.type) {
    case "ADD_PROJECTS": {
      // a given custom project may already exist in the group so don't trigger a state update for those duplicates
      const uniqueCustomProjects = action.payload.projects.filter(
        (project: string) => !(project in state[action.payload.group])
      );
      if (uniqueCustomProjects.length === 0) {
        return state;
      }
      return {
        ...state,
        [action.payload.group]: {
          ...state[action.payload.group],
          ...Object.fromEntries(
            uniqueCustomProjects.map(v => [v, { checked: true, custom: true }])
          ),
        },
      };
    }
    case "REMOVE_PROJECTS": {
      // TODO: also remove any upstreams or downstreams related (only) to the project.
      // if group == Groups.PROJECT, hide exclusive downstream upstreams
      //
      return {
        ...state,
        [action.payload.group]: _.omit(state[action.payload.group], action.payload.projects),
      };
    }
    case "TOGGLE_PROJECTS": {
      // TODO: hide exclusive upstreams and downstreams if group is PROJECTS
      return {
        ...state,
        [action.payload.group]: {
          ...state[action.payload.group],
          ...Object.fromEntries(
            action.payload.projects.map(key => [
              key,
              {
                ...state[action.payload.group][key],
                checked: !state[action.payload.group][key].checked,
              },
            ])
          ),
        },
      };
    }
    case "ONLY_PROJECTS": {
      const newOnlyProjectState = { ...state };

      newOnlyProjectState[action.payload.group] = Object.fromEntries(
        Object.keys(state[action.payload.group]).map(key => [
          key,
          { ...state[action.payload.group][key], checked: action.payload.projects.includes(key) },
        ])
      );

      return newOnlyProjectState;
    }
    case "TOGGLE_ENTIRE_GROUP": {
      const newCheckedValue = !deriveSwitchStatus(state, action.payload.group);
      const newGroupToggledState = { ...state };
      newGroupToggledState[action.payload.group] = Object.fromEntries(
        Object.keys(state[action.payload.group]).map(key => [
          key,
          { ...state[action.payload.group][key], checked: newCheckedValue },
        ])
      );

      return newGroupToggledState;
    }
    // Background actions.
    case "HYDRATE_START": {
      return { ...state, loading: true };
    }
    case "HYDRATE_END": {
      const newPostAPICallState = { ...state, loading: false };
      // TODO: handle payload.
      _.forIn(action.payload.result, (v, k) => {
        // Add each project to the projects list.
        state[Group.PROJECTS][k] = { checked: true };
        state.projectData[k] = {};

        // Add each upstream.
        v.upstreams.forEach(v => {
          state[Group.UPSTREAM][v] = { checked: false };
          state.projectData[v] = {};
        });

        // Add each downstream.
        v.downstreams.forEach(v => {
          state[Group.DOWNSTREAM][v] = { checked: false };
          state.projectData[v] = {};
        });

        // Update project data for each.
      });
      return newPostAPICallState;
    }
    default:
      throw new Error(`unknown resolver action`);
  }
};

const fakeAPI = (state: State) => {
  return {
    clutch: {
      upstreams: ["rides", "locations"],
      downstreams: ["queueworker"],
    },
  };
};

// TODO(perf): call with useMemo().
const deriveSwitchStatus = (state: State, group: Group): boolean => {
  return (
    Object.keys(state[group]).length > 0 &&
    Object.keys(state[group]).every(key => state[group][key].checked)
  );
};

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

const StyledTitle = styled.span({
  fontWeight: 500,
});

const MenuItem = styled.div({
  display: "flex",
  alignItems: "center",
  "&:hover": {
    backgroundColor: "rgba(13, 16, 48, 0.03)",
  },
  "&:hover > span": {
    display: "inline-flex", // Unhide "only" hidden span.
  },
});

const StyledHeader = styled.div({
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

const initialState: State = {
  [Group.PROJECTS]: {},
  [Group.UPSTREAM]: {},
  [Group.DOWNSTREAM]: {},
  projectData: {},
  loading: false,
};

const StyledClearIcon = styled.span({
  color: "grey",
  display: "inline-flex",
  width: "24px",
});

const StyledOnlyButton = styled.span({
  marginRight: "10px",
  color: "rgba(53, 72, 212, 1)",
});

const ProjectGroup = ({
  title,
  group,
  collapsible,
}: {
  title: string;
  group: Group;
  collapsible?: boolean;
}) => {
  const dispatch = useDispatch();
  const state = useReducerState();

  const [collapsed, setCollapsed] = React.useState(false);

  const numProjects = Object.keys(state[group]).length;
  const checkedProjects = Object.keys(state[group]).filter(k => state[group][k].checked);

  return (
    <>
      <StyledHeader>
        <StyledHeaderColumn>
          {collapsible && (
            <StyledHeaderColumn onClick={() => setCollapsed(!collapsed)}>
              {collapsed ? <ChevronRightIcon /> : <ExpandMoreIcon />}
            </StyledHeaderColumn>
          )}
          <StyledTitle>
            {title}
            <StyledCount>
              {checkedProjects.length}
              {numProjects > 0 && "/" + numProjects}
            </StyledCount>
          </StyledTitle>
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
      </StyledHeader>
      {!collapsed && (
        <div>
          {numProjects == 0 && <div>No projects in this group.</div>}
          {Object.keys(state[group]).map(key => (
            <MenuItem key={key}>
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
            </MenuItem>
          ))}
        </div>
      )}
    </>
  );
};

const SelectorContainer = styled.div({
  backgroundColor: "#F9FAFE",
});

const ProjectSelector = () => {
  // On load, we'll request a list of owned projects and their upstreams and downstreams from the API.
  // The API will contain information about the relationships between projects and upstreams and downstreams.
  // By default, the owned projects will be checked and others will be unchecked.
  // If a project is unchecked, the upstream and downstreams related to it disappear from the list.
  // If a project is rechecked, the checks were preserved.

  const [customProject, setCustomProject] = React.useState("");

  const [state, dispatch] = React.useReducer(selectorReducer, initialState);

  React.useEffect(() => {
    console.log("effect");
    // Determine if any hydration is required.
    // - Are any services missing from state.projectdata?
    // - Are projects empty (first load)?
    // - Is loading not already in progress?

    let allPresent = true;
    _.forEach(Object.keys(state[Group.PROJECTS]), p => {
      if (!(p in state.projectData)) {
        allPresent = false;
        return false; // Stop iteration.
      }
      return true; // Continue.
    });

    if (!state.loading && (Object.keys(state[Group.PROJECTS]).length == 0 || !allPresent)) {
      console.log("calling API!", state.loading);
      dispatch({ type: "HYDRATE_START" });
      // TODO: call API and use payload.
      setTimeout(
        () => dispatch({ type: "HYDRATE_END", payload: { result: fakeAPI(state) } }),
        1000
      );
    }
  }, [state[Group.PROJECTS]]);

  const handleAdd = () => {
    if (customProject === "") {
      return;
    }
    dispatch({
      type: "ADD_PROJECTS",
      payload: { group: Group.PROJECTS, projects: [customProject] },
    });
    setCustomProject("");
  };

  return (
    <DispatchContext.Provider value={dispatch}>
      <StateContext.Provider value={state}>
        <SelectorContainer>
          {state.loading && <LinearProgress color="secondary" />}
          <div>
            <LayersIcon />
            Dash
          </div>
          <div>
            <TextField
              disabled={state.loading}
              placeholder="Add a project"
              value={customProject}
              onChange={e => setCustomProject(e.target.value)}
              onKeyDown={e => e.key === "Enter" && handleAdd()}
            />
          </div>
          <ProjectGroup title="Projects" group={Group.PROJECTS} />
          <Divider />
          <ProjectGroup title="Upstreams" group={Group.UPSTREAM} collapsible />
          <Divider />
          <ProjectGroup title="Downstreams" group={Group.DOWNSTREAM} collapsible />
        </SelectorContainer>
      </StateContext.Provider>
    </DispatchContext.Provider>
  );
};

const HelloWorld = () => (
  <>
    <ProjectSelector />
  </>
);

export default HelloWorld;
