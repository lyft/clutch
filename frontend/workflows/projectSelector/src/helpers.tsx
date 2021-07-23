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

const DispatchContext = React.createContext<(action: Action) => void | undefined>(() => undefined);
const useDispatch = () => {
  return React.useContext(DispatchContext);
};

// TODO(perf): call with useMemo().
const deriveSwitchStatus = (state: State | undefined, group: Group): boolean => {
  const groupKeys = Object.keys(state?.[group] || []);
  return groupKeys.length > 0 && groupKeys.every(key => state?.[group][key].checked);
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
    upstreams?.[PROJECT_TYPE_URL]?.ids?.forEach(u => {
      if (!upstreamMap[u]) {
        upstreamMap[u] = { [p]: true };
      } else {
        upstreamMap[u][p] = true;
      }
    });
    downstreams?.[PROJECT_TYPE_URL]?.ids?.forEach(d => {
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

  const upstreams = [] as string[];
  const downstreams = [] as string[];
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
// TODO: (perf/efficiency) compute the projects that should be displayed rather than computing the hidden projects
// returns the upstreams/downstreams that should be hidden based on the checked status of/exclusivity to project(s)
const deriveHiddenDependencies = (state: State): { upstreams: string[]; downstreams: string[] } => {
  const upstreams = [] as string[];
  const downstreams = [] as string[];

  const uncheckedProjects = Object.keys(state[Group.PROJECTS]).filter(
    p => !state[Group.PROJECTS][p].checked
  );
  // no unchecked projects, so don't go through rest of the flow
  if (uncheckedProjects.length === 0) {
    return { upstreams, downstreams };
  }

  // get the relationship b/w upstreams/downstreams to Group.Projects
  const depMapping = dependencyToProjects(state, Group.PROJECTS);

  uncheckedProjects.forEach(project => {
    _.forIn(depMapping.upstreams, (v, k) => {
      // if dependency is exclusive to the unchecked project, hide the dependency
      if (v[project] && Object.keys(v).length === 1) {
        upstreams.push(k);
      } else if (Object.keys(v).every(p => !state[Group.PROJECTS][p].checked)) {
        // if all Group.Projects that share the dependency are unchecked, hide the dependency
        upstreams.push(k);
      }
    });

    _.forIn(depMapping.downstreams, (v, k) => {
      // if dependency is exclusive to the unchecked project, hide the dependency
      if (v[project] && Object.keys(v).length === 1) {
        downstreams.push(k);
      } else if (Object.keys(v).every(p => !state[Group.PROJECTS][p].checked)) {
        // if all Group.Projects that share the dependency are unchecked, hide the dependency
        downstreams.push(k);
      }
    });
  });

  return { upstreams, downstreams };
};

/* filter out the state data on the following criteria:
if projects in Group.Projects are unchecked, remove the related upstream(s)/downstream(s)
*/
const deriveStateData = (state: State): State => {
  // get the upstreams/downstreams that should be omitted from the final state data
  const { upstreams, downstreams } = deriveHiddenDependencies(state);

  const newState = { ...state };
  if (upstreams.length > 0) {
    newState[Group.UPSTREAM] = _.omit(state[Group.UPSTREAM], upstreams);
  }
  if (downstreams.length > 0) {
    newState[Group.DOWNSTREAM] = _.omit(state[Group.DOWNSTREAM], downstreams);
  }

  return newState;
};

export {
  deriveStateData,
  deriveSwitchStatus,
  DispatchContext,
  exclusiveProjectDependencies,
  PROJECT_TYPE_URL,
  StateContext,
  updateGroupstate,
  useDispatch,
  useReducerState,
};
