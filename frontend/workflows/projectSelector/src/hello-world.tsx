import * as React from "react";

import _ from "lodash";

import { Button, Checkbox, Switch, TextField } from "@clutch-sh/core";
import { Divider, IconButton } from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import HelpOutlineIcon from "@material-ui/icons/HelpOutline";
import ClearIcon from "@material-ui/icons/Clear";
import LayersIcon from "@material-ui/icons/Layers";

import styled from "@emotion/styled";

enum ActionKind {
  ADD_PROJECTS,
  REMOVE_PROJECTS,
  TOGGLE_PROJECTS,
  TOGGLE_ENTIRE_GROUP,
  ONLY_PROJECTS,
}

enum Group {
  PROJECTS,
  UPSTREAM,
  DOWNSTREAM,
}

interface Action {
  type: ActionKind;
  payload: Payload;
}

interface Payload {
  group: Group;
  projects?: string[];
}

interface State {
  [Group.PROJECTS]: GroupState;
  [Group.UPSTREAM]: GroupState;
  [Group.DOWNSTREAM]: GroupState;
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
  const { type, payload } = action;

  switch (type) {
    case ActionKind.ADD_PROJECTS:
      // TODO: don't allow adding if it already exists.
      return {
        ...state,
        [payload.group]: {
          ...state[payload.group],
          ...Object.fromEntries(payload.projects.map(v => [v, { checked: true, custom: true }])),
        },
      };
    case ActionKind.REMOVE_PROJECTS:
      // TODO: also remove any upstreams or downstreams related (only) to the project.
      return {
        ...state,
        [payload.group]: _.omit(state[payload.group], payload.projects),
      };
    case ActionKind.TOGGLE_PROJECTS:
      // TODO: hide upstreams and downstreams if group is PROJECTS
      return {
        ...state,
        [payload.group]: {
          ...state[payload.group],
          ...Object.fromEntries(
            payload.projects.map(key => [
              key,
              { ...state[payload.group][key], checked: !state[payload.group][key].checked },
            ])
          ),
        },
      };
    case ActionKind.ONLY_PROJECTS:
      const newOnlyProjectState = { ...state };

      newOnlyProjectState[payload.group] = Object.fromEntries(
        Object.keys(state[payload.group]).map(key => [
          key,
          { ...state[payload.group][key], checked: payload.projects.includes(key) },
        ])
      );

      return newOnlyProjectState;

    case ActionKind.TOGGLE_ENTIRE_GROUP:
      const newCheckedValue = !deriveSwitchStatus(state, payload.group);
      const newGroupToggledState = { ...state };
      newGroupToggledState[payload.group] = Object.fromEntries(
        Object.keys(state[payload.group]).map(key => [
          key,
          { ...state[payload.group][key], checked: newCheckedValue },
        ])
      );

      return newGroupToggledState;
    default:
      throw new Error(`unknown resolver action: ${type}`);
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
      {collapsible && <ExpandMoreIcon />} {title}
      <Switch
        onChange={() =>
          dispatch({
            type: ActionKind.TOGGLE_ENTIRE_GROUP,
            payload: { group: group },
          })
        }
        checked={deriveSwitchStatus(state, group)}
        disabled={Object.keys(state[group]).length == 0}
      />
      <div>
        {Object.keys(state[group]).length == 0 && <div>No projects in this group.</div>}
        {Object.keys(state[group]).map(key => (
          <div key={key} className="project">
            <Checkbox
              name={key}
              onChange={() =>
                dispatch({
                  type: ActionKind.TOGGLE_PROJECTS,
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
                  type: ActionKind.ONLY_PROJECTS,
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
                    type: ActionKind.REMOVE_PROJECTS,
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
    // Here we would call the API and load the initial set of projects.
  }, []);

  const upstreams = ["authors", "thumbnails"];
  const downstreams = ["coffee", "shelves"];

  const handleAdd = () => {
    if (customProject === "") {
      return;
    }
    dispatch({
      type: ActionKind.ADD_PROJECTS,
      payload: { group: Group.PROJECTS, projects: [customProject] },
    });
    setCustomProject("");
  };

  return (
    <DispatchContext.Provider value={dispatch}>
      <StateContext.Provider value={state}>
        <SelectorContainer>
          <div>
            <LayersIcon />
            Dash
          </div>
          <div>
            <TextField
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
