import * as React from "react";

import _ from "lodash";

import { Button, Checkbox, Switch, TextField } from "@clutch-sh/core";
import { Divider, IconButton, LinearProgress } from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
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

  projectData: { [key: string]: Project };
  loading: boolean;
}

interface Project {
  upstreams: string[];
  downstreams: string[];
}

interface GroupState {
  [s: string]: ProjectState;
}

interface ProjectState {
  checked: boolean;
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
    case "ADD_PROJECTS":
      // TODO: don't add if it already exists.
      // TODO: refresh API if project or its upstreams and downstreams are not present in state.
      return {
        ...state,
        [action.payload.group]: {
          ...state[action.payload.group],
          ...Object.fromEntries(action.payload.projects.map(v => [v, { checked: true, custom: true }])),
        },
      };
    case "REMOVE_PROJECTS":
      // TODO: also remove any upstreams or downstreams related (only) to the project.
      return {
        ...state,
        [action.payload.group]: _.omit(state[action.payload.group], action.payload.projects),
      };
    case "TOGGLE_PROJECTS":
      // TODO: hide upstreams and downstreams if group is PROJECTS
      return {
        ...state,
        [action.payload.group]: {
          ...state[action.payload.group],
          ...Object.fromEntries(
            action.payload.projects.map(key => [
              key,
              { ...state[action.payload.group][key], checked: !state[action.payload.group][key].checked },
            ])
          ),
        },
      };
    case "ONLY_PROJECTS":
      const newOnlyProjectState = { ...state };

      newOnlyProjectState[action.payload.group] = Object.fromEntries(
        Object.keys(state[action.payload.group]).map(key => [
          key,
          { ...state[action.payload.group][key], checked: action.payload.projects.includes(key) },
        ])
      );

      return newOnlyProjectState;

    case "TOGGLE_ENTIRE_GROUP":
      const newCheckedValue = !deriveSwitchStatus(state, action.payload.group);
      const newGroupToggledState = { ...state };
      newGroupToggledState[action.payload.group] = Object.fromEntries(
        Object.keys(state[action.payload.group]).map(key => [
          key,
          { ...state[action.payload.group][key], checked: newCheckedValue },
        ])
      );

      return newGroupToggledState;

    // Background actions.
    case "HYDRATE_START":
      return { ...state, loading: true };

    case "HYDRATE_END":
      return { ...state, loading: false };
    default:
      throw new Error(`unknown resolver action`);
  }
};

// TODO(perf): call with useMemo().
const deriveSwitchStatus = (state: State, group: Group): boolean => {
  return (
    Object.keys(state[group]).length > 0 &&
    Object.keys(state[group]).every(key => state[group][key].checked)
  );
};

const initialState: State = {
  [Group.PROJECTS]: {},
  [Group.UPSTREAM]: {},
  [Group.DOWNSTREAM]: {},
  loading: true,
  projectData: {}
};

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

  return (
    <div>
      <div>
        {collapsible && <ExpandMoreIcon />}
        {title}
        {!collapsible && "All"}
        <Switch
          onChange={() =>
            dispatch({
              type: "TOGGLE_ENTIRE_GROUP",
              payload: { group: group },
            })
          }
          checked={deriveSwitchStatus(state, group)}
          disabled={Object.keys(state[group]).length == 0}
        />
      </div>
      <div>
        {Object.keys(state[group]).length == 0 && <div>No projects in this group.</div>}
        {Object.keys(state[group]).map(key => (
          <div key={key} className="project">
            <Checkbox
              name={key}
              onChange={() =>
                dispatch({
                  type: "TOGGLE_PROJECTS",
                  payload: { group, projects: [key] },
                })
              }
              checked={state[group][key].checked ? true : false}
            />
            {key}
            <div
              className="only"
              onClick={() =>
                dispatch({
                  type: "ONLY_PROJECTS",
                  payload: { group, projects: [key] },
                })
              }
            >
              Only
            </div>
            {state[group][key].custom && (
              <ClearIcon
                onClick={() =>
                  dispatch({
                    type: "REMOVE_PROJECTS",
                    payload: { group, projects: [key] },
                  })
                }
              />
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

const SelectorContainer = styled.div({
  backgroundColor: "#F5F6FD",
  ".project": {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  ".only": {
    color: "#3548D4",
  },
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
    // - Are projects empty?
    let allPresent = true;
    Object.keys(state[Group.PROJECTS]).forEach(p => {
      if (allPresent && !(p in state.projectData)) {
        allPresent = false;
      }
    });

    if (Object.keys(state[Group.PROJECTS]).length == 0 || !allPresent) {
      console.log("calling API!")
      dispatch({ type: "HYDRATE_START"});
      // Here we would call the API and load projects if needed.
      dispatch({ type: "HYDRATE_END", payload: {result: {}} });
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
