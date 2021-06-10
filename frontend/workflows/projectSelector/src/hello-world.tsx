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
  CUSTOM,
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
  [Group.CUSTOM]: GroupState;
  [Group.UPSTREAM]: GroupState;
  [Group.DOWNSTREAM]: GroupState;
}

interface GroupState {
  [s: string]: ProjectState;
}

interface ProjectState {
  checked: boolean;
  hidden?: boolean; // upstreams and downstreams are hidden when their parent is unchecked unless other parents also use them.
}

const selectorReducer = (state: State, action: Action): State => {
  const { type, payload } = action;

  switch (type) {
    case ActionKind.ADD_PROJECTS:
      // TODO: don't allow adding a group to CUSTOM if it's already in PROJECTS.
      return {
        ...state,
        [payload.group]: {
          ...state[payload.group],
          ...Object.fromEntries(payload.projects.map(v => [v, { checked: true }])),
        },
      };
    case ActionKind.REMOVE_PROJECTS:
      // TODO: also remove any upstreams or downstreams related (only) to the project.
      return {
        ...state,
        [payload.group]: _.omit(state[payload.group], payload.projects),
      };
    case ActionKind.TOGGLE_PROJECTS:
      // TODO: hide upstreams and downstreams if group is PROJECTS OR CUSTOM
      const { group, projects } = payload;
      return {
        ...state,
        [group]: {
          ...state[group],
          ...Object.fromEntries(
            projects.map(key => [
              key,
              { ...state[group][key], checked: !state[group][key].checked },
            ])
          ),
        },
      };
    case ActionKind.ONLY_PROJECTS:
      const applyOnlyToGroups = [payload.group];
      if (payload.group === Group.PROJECTS) {
        // If the group is the PROJECTS group, toggle both PROJECTS and CUSTOM.
        applyOnlyToGroups.push(Group.CUSTOM);
      } else if (payload.group === Group.CUSTOM) {
        applyOnlyToGroups.push(Group.PROJECTS);
      }

      const newOnlyProjectState = { ...state };

      applyOnlyToGroups.forEach(group => {
        newOnlyProjectState[group] = Object.fromEntries(
          Object.keys(state[group]).map(key => [
            key,
            { ...state[group][key], checked: payload.projects.includes(key) },
          ])
        );
      });

      return newOnlyProjectState;

    case ActionKind.TOGGLE_ENTIRE_GROUP:
      const applicableGroups = [payload.group];
      if (payload.group === Group.PROJECTS) {
        // If the group is the PROJECTS group, toggle both PROJECTS and CUSTOM.
        applicableGroups.push(Group.CUSTOM);
      }

      const newCheckedValue = !determineSwitchStatus(state, applicableGroups);
      const newGroupToggledState = { ...state };
      applicableGroups.forEach(group => {
        newGroupToggledState[group] = Object.fromEntries(
          Object.keys(state[group]).map(key => [
            key,
            { ...state[group][key], checked: newCheckedValue },
          ])
        );
      });

      return newGroupToggledState;
    default:
      throw new Error(`unknown resolver action: ${type}`);
  }
};

// TODO(perf): call with useMemo().
const determineSwitchStatus = (state: State, groups: Group[]): boolean => {
  return (
    Object.keys(state[Group.CUSTOM]).length > 0 &&
    Object.keys(state[Group.CUSTOM]).every(key => state[Group.CUSTOM][key].checked)
  );
};

const initialState: State = {
  [Group.PROJECTS]: {},
  [Group.CUSTOM]: {},
  [Group.UPSTREAM]: {},
  [Group.DOWNSTREAM]: {},
};

const Selector = styled.div({
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

  const changeHandler = ({ target }) => {
    // setProjects({ ...projects, [target.name]: target.checked });
  };

  const handleAdd = () => {
    if (customProject === "") {
      return;
    }
    dispatch({
      type: ActionKind.ADD_PROJECTS,
      payload: { group: Group.CUSTOM, projects: [customProject] },
    });
    setCustomProject("");
  };

  return (
    <Selector>
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
      <div>
        Projects
        <Switch
          onChange={() =>
            dispatch({ type: ActionKind.TOGGLE_ENTIRE_GROUP, payload: { group: Group.PROJECTS } })
          }
          checked={determineSwitchStatus(state, [Group.PROJECTS, Group.CUSTOM])}
          disabled={
            Object.keys(state[Group.PROJECTS]).length + Object.keys(state[Group.CUSTOM]).length == 0
          }
        />
      </div>
      <div>
        {/* {Object.keys(projects).map(key => (
          <div key={key}>
            <Checkbox name={key} onChange={changeHandler} checked={projects[key]} /> {key}
          </div>
        ))} */}
        {Object.keys(state[Group.CUSTOM]).map(key => (
          <div key={key} className="project">
            <Checkbox
              name={key}
              onChange={() =>
                dispatch({
                  type: ActionKind.TOGGLE_PROJECTS,
                  payload: { group: Group.CUSTOM, projects: [key] },
                })
              }
              checked={state[Group.CUSTOM][key].checked ? true : false}
            />
            {key}
            <div
              className="only"
              onClick={() => 
                dispatch({type: ActionKind.ONLY_PROJECTS, payload: {group: Group.CUSTOM, projects: [key]}})
              }
            >
              Only
            </div>
            <ClearIcon
              onClick={() =>
                dispatch({
                  type: ActionKind.REMOVE_PROJECTS,
                  payload: { group: Group.CUSTOM, projects: [key] },
                })
              }
            />
          </div>
        ))}
      </div>
      <Divider />
      <div>
        <ExpandMoreIcon />
        Downstreams
        <Switch />
      </div>
      <div>
        {upstreams.map(v => (
          <div key={v}>
            <Checkbox /> {v}
          </div>
        ))}
      </div>
      <Divider />
      <div>
        <ExpandMoreIcon />
        Upstreams
        <Switch />
      </div>
      <div>
        {downstreams.map(v => (
          <div key={v}>
            <Checkbox /> {v}
          </div>
        ))}
      </div>
    </Selector>
  );
};

const HelloWorld = () => (
  <>
    <ProjectSelector />
  </>
);

export default HelloWorld;
