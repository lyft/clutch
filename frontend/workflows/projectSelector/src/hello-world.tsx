import * as React from "react";

import _ from "lodash";

import { Button, Checkbox, Switch, TextField } from "@clutch-sh/core";
import { Divider, IconButton } from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import HelpOutlineIcon from "@material-ui/icons/HelpOutline";
import ClearIcon from "@material-ui/icons/Clear";
import LayersIcon from "@material-ui/icons/Layers";

enum ActionKind {
  ADD_PROJECTS,
  REMOVE_PROJECTS,
  TOGGLE_PROJECTS,
  ONLY_PROJECTS,
}

interface Action {
  type: ActionKind;
  payload: Payload;
}

interface Payload {
  projects?: string[];
}

interface State {
  projects: any;
  customProjects: any;
  upstreams: any;
  downstreams: any;
}

const selectorReducer = (state: State, action: Action): State => {
  const { type, payload } = action;

  switch (type) {
    case ActionKind.ADD_PROJECTS:
      return {
        ...state,
        customProjects: {
          ...state.customProjects,
          ...Object.fromEntries(payload.projects.map(v => [v, null])),
        },
      };
    case ActionKind.REMOVE_PROJECTS:
      return {
        ...state,
        customProjects: _.omit(state.customProjects, payload.projects),
      };
    default:
      throw new Error(`unknown resolver action: ${type}`);
  }
};

enum GroupSelectedStatus {
  NONE,
  SOME,
  ALL,
}

const determineSelectedStatus = (group: string): GroupSelectedStatus => {};

const initialState: State = {
  projects: {},
  customProjects: {},
  upstreams: {},
  downstreams: {},
};

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
    dispatch({ type: ActionKind.ADD_PROJECTS, payload: { projects: [customProject] } });
    setCustomProject("");
  };

  return (
    <div>
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
        <IconButton onClick={handleAdd}>
          <AddIcon />
        </IconButton>
      </div>
      <div>
        Projects
        <Switch />
      </div>
      <div>
        {/* {Object.keys(projects).map(key => (
          <div key={key}>
            <Checkbox name={key} onChange={changeHandler} checked={projects[key]} /> {key}
          </div>
        ))} */}
        {Object.keys(state.customProjects).map(key => (
          <div key={key}>
            <Checkbox name={key} onChange={changeHandler} checked />
            {key}
            <ClearIcon
              onClick={() =>
                dispatch({ type: ActionKind.REMOVE_PROJECTS, payload: { projects: [key] } })
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
    </div>
  );
};

const HelloWorld = () => (
  <>
    <ProjectSelector />
  </>
);

export default HelloWorld;
