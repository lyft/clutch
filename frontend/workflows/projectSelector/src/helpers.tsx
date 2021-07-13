import * as React from "react";
import _ from "lodash";

import type { Action, ProjectState, State } from "./types";
import { Group } from "./types";

const PROJECT_TYPE_URL = "type.googleapis.com/clutch.core.project.v1.Project";

interface DependencyMappings {
  upstreams?: { [dependency: string]: { [project: string]: boolean } };
  downstreams?: { [dependency: string]: { [project: string]: boolean } };
}

const StateContext = React.createContext<State | undefined>(undefined);
const useReducerState = () => {
  return React.useContext(StateContext);
};

const DispatchContext = React.createContext<(action: Action) => void | undefined>(undefined);
const useDispatch = () => {
  return React.useContext(DispatchContext);
};

// TODO(perf): call with useMemo().
const deriveSwitchStatus = (state: State, group: Group): boolean => {
  return (
    Object.keys(state[group]).length > 0 &&
    Object.keys(state[group]).every(key => state[group][key].checked)
  );
};

const updateGroupstate = (
  state: State,
  group: Group,
  project: string,
  projectState: ProjectState
): State => {
  const newState = { ...state };
  if (project in newState[group]) {
    // preserve the checked value if the project is already in the group
    newState[group][project].checked = state[group][project].checked;
    newState[group][project].custom = projectState.custom;
  } else {
    // we set the projectState to the default state passed in
    newState[group][project] = projectState;
  }

  return newState;
};

const dependencyToProjects = (state: State, group: Group): DependencyMappings => {
  const upstreamMap = {};
  const downstreamMap = {};

  const projects = Object.keys(state[group]);
  projects.forEach(p => {
    const { upstreams, downstreams } = state.projectData[p]?.dependencies || {};
    upstreams?.[PROJECT_TYPE_URL]?.ids.forEach(u => {
      if (!upstreamMap[u]) {
        upstreamMap[u] = { [p]: true };
      } else {
        upstreamMap[u][p] = true;
      }
    });
    downstreams?.[PROJECT_TYPE_URL]?.ids.forEach(d => {
      if (!downstreamMap[d]) {
        downstreamMap[d] = { [p]: true };
      } else {
        downstreamMap[d][p] = true;
      }
    });
  });

  return { upstreams: upstreamMap, downstreams: downstreamMap };
};

const exclusiveProjectDependencies = (
  state: State,
  group: Group,
  project: string
): { upstreams: string[]; downstreams: string[] } => {
  const dependencyMap = dependencyToProjects(state, group);

  const upstreams = [];
  const downstreams = [];
  _.forIn(dependencyMap.upstreams, (v, k) => {
    if (v[project] && Object.keys(v).length === 1) {
      upstreams.push(k);
    }
  });

  _.forIn(dependencyMap.downstreams, (v, k) => {
    if (v[project] && Object.keys(v).length === 1) {
      downstreams.push(k);
    }
  });
  return { upstreams, downstreams };
};

// TODO: if upstreams/downstreams are shared by all other projects that have already hidden their upstreams/downstreams
// hide those upstreams/downstreams as well
const updateHiddenState = (state: State, group: Group, project: string): State => {
  const newState = { ...state };

  const dependencyMapping = dependencyToProjects(state, group);

  _.forIn(dependencyMapping.upstreams, (v, k) => {
    if (v[project]) {
      if (Object.keys(v).length === 1) {
        newState[Group.UPSTREAM][k].hidden = !state[Group.UPSTREAM][k].hidden;
      } else {
        // the state could have changed and now that project is shared by others so reset to false
        newState[Group.UPSTREAM][k].hidden = false;
      }
    }
  });

  _.forIn(dependencyMapping.downstreams, (v, k) => {
    if (v[project]) {
      if (Object.keys(v).length === 1) {
        newState[Group.DOWNSTREAM][k].hidden = !state[Group.DOWNSTREAM][k].hidden;
      } else {
        // the state could have changed and now that project is shared by others so reset to false
        newState[Group.UPSTREAM][k].hidden = false;
      }
    }
  });

  return newState;
};

export {
  deriveSwitchStatus,
  DispatchContext,
  exclusiveProjectDependencies,
  PROJECT_TYPE_URL,
  StateContext,
  updateGroupstate,
  updateHiddenState,
  useDispatch,
  useReducerState,
};
